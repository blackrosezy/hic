#!/usr/bin/python

"""Hic

Usage:
  hic add <url> <container_name> [<port>]
  hic rm <url>
  hic sync
"""
from docopt import docopt
from tabulate import tabulate

from rediscli import RedisCli
from dockercli import DockerCli
from configuration import Configuration
from hicparser import Parser


class Hic:
    def __init__(self):
        self.dockercli = DockerCli()
        self.rediscli = RedisCli()
        self.parser = Parser()
        self.config = Configuration()

    def add_command(self, url, container_name, port=80):
        err, protocol, domain = self.parser.split_url(url)
        if err:
            print " => Invalid url '%s'." % url
            return
        self.add(protocol, domain, port, container_name)

    def remove_command(self, url):
        err, protocol, domain = self.parser.split_url(url)
        if err:
            print " => Invalid url '%s'." % url
            return
        found_url = False
        for item in self.config.list_config_flatten():
            if item['domain'] == domain:
                found_url = True
                full_ip = '%s://%s:%s' % (protocol, item['container_ip'], item['host_port'])
                url_key = 'frontend:%s' % domain
                self.rediscli.redis_conn.lrem(url_key, 0, full_ip)
                print " => Removing url '%s://%s' successful!" % (protocol, domain)
        if found_url:
            self.config.remove_config(domain)
        else:
            print " => Skipping remove '%s://%s'." % (protocol, domain)


    def add(self, protocol, domain, port, container_name):
        container_ip = ''
        host_port = 0
        for container in self.dockercli.list_containers_filtered():
            if container['name'] == container_name:
                container_ip, host_port = self.dockercli.get_ip_and_host_port(container, port)
                break
        if container_ip == '' or host_port == 0:
            print " => Cannot find container/port for '%s/%d'." % (container_name, port)
            return

        full_ip = '%s://%s:%s' % (protocol, container_ip, host_port)
        url_key = 'frontend:%s' % domain
        if self.rediscli.redis_conn.lindex(url_key, 0) == None:
            self.rediscli.redis_conn.rpush(url_key, '__%s__' % domain)

        if full_ip not in self.rediscli.redis_conn.lrange(url_key, 0, -1):
            self.rediscli.redis_conn.rpush(url_key, full_ip)
            print " => Adding '%s://%s' successful!" % (protocol, domain)
        self.config.update_config(protocol, domain, port, container_name)

    def sync(self):
        working_list = []
        # create working list
        for item in self.config.list_config_flatten():
            if item['host_port'] != 0:
                self.add(item['protocol'], item['domain'], item['port'], item['container_name'])
                working_list.append(item)
            else:
                print " => Skipping '%s://%s:%d' - '%s'." % (
                    item['protocol'], item['domain'], item['port'], item['container_name'])

        # remove unmatched redis items
        for item in self.rediscli.list_redis_items_flatten():
            found_item = False
            for working_item in working_list:
                if item['url'] == working_item['domain'] and item['ip'] == working_item['container_ip'] and item[
                    'protocol'] == working_item['protocol'] and item['host_port'] == working_item['host_port']:
                    found_item = True
                    break

            if found_item:
                continue

            url_key = 'frontend:%s' % item['url']
            self.rediscli.redis_conn.lrem(url_key, 0, item['original'])

        # remove urls without any ips
        for item in self.rediscli.list_redis_items():
            url_key = 'frontend:%s' % item['url']
            if not len(item['data']):
                self.rediscli.redis_conn.delete(url_key)

    def show(self):
        rows = []
        for item in self.config.list_config_flatten():
            row = []
            row.append('%s://%s' % (item['protocol'], item['domain'] ))
            row.append(item['container_ip'])
            port = 'Nil'
            if item['host_port']:
                port = item['host_port']
            row.append(port)
            row.append(item['port'])
            row.append(item['container_name'])

            rows.append(row)
        print tabulate(rows, headers=['URL', 'IP', 'HOST PORT', 'CONT. PORT', 'CONTAINER NAME'])


if __name__ == '__main__':
    cli = docopt(__doc__, version='Hic 2.0')
    hic = Hic()
    if cli['add']:
        port = 0
        try:
            if cli['<port>']:
                port = int(cli['<port>'])
        except ValueError:
            pass
        if port:
            hic.add_command(cli['<url>'], cli['<container_name>'], port)
        else:
            hic.add_command(cli['<url>'], cli['<container_name>'])
        hic.sync()
        hic.show()
    elif cli['rm']:
        hic.remove_command(cli['<url>'])
        hic.sync()
        hic.show()
    elif cli['sync']:
        hic.sync()
        hic.show()