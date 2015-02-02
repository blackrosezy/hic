import json

from hipache_cli.dockercli import DockerCli


class Configuration:
    CONFIG_FILE = '/root/.hic.json'

    def __init__(self):
        self.dockercli = DockerCli()

    def read_config(self):
        try:
            with open(self.CONFIG_FILE, 'r') as infile:
                data = json.load(infile)
        except (IOError, ValueError) as e:
            data = []
        return data

    def update_config(self, protocol, domain, port, container_name):
        data = self.read_config()

        match_item_url = None
        match_item_data = None
        for item_url in data:
            if item_url['domain'] == domain:
                match_item_url = item_url
                for item_data in item_url['data']:
                    if item_data['protocol'] == protocol and item_data['container_name'] == container_name and \
                                    item_data['port'] == port:
                        match_item_data = item_data
                        break
                break

        item_url = {}
        item_url['domain'] = domain
        item_url['data'] = []

        item_data = {}
        item_data['protocol'] = protocol
        item_data['container_name'] = container_name
        item_data['port'] = port

        if not match_item_url:
            item_url['data'].append(item_data)
            data.append(item_url)
        elif not match_item_data:
            match_item_url['data'].append(item_data)

        with open(self.CONFIG_FILE, 'w') as outfile:
            json.dump(data, outfile, indent=2)
        return data

    def remove_config(self, domain):
        new_config = []
        for item in self.list_config_flatten():
            if item['domain'] != domain:
                new_config.append(item)

        with open(self.CONFIG_FILE, 'w') as outfile:
            outfile.write('')

        for item in new_config:
            self.update_config(item['protocol'], item['domain'], item['port'], item['container_name'])


    def list_config_flatten(self):
        new_arr = []
        containers = self.dockercli.list_containers_filtered()
        for item in self.read_config():
            domain = item['domain']

            for item_data in item['data']:
                container_ip = ''
                host_port = 0
                for container in containers:
                    if container['name'] == item_data['container_name']:
                        container_ip, host_port = self.dockercli.get_ip_and_host_port(container, item_data['port'])
                        break

                new_map = {}
                new_map['domain'] = domain
                new_map['protocol'] = item_data['protocol']
                new_map['port'] = item_data['port']
                new_map['container_name'] = item_data['container_name']
                new_map['container_ip'] = container_ip
                new_map['host_port'] = host_port
                new_arr.append(new_map)
        return new_arr
