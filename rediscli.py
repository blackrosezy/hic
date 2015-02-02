import redis

from hicparser import Parser
from dockercli import DockerCli


class RedisCli:
    def __init__(self):
        self.dockercli = DockerCli()
        self.__connect_to_redis()

    def __connect_to_redis(self):
        hipache_container = self.__find_hipache_container()
        hipache_ip = hipache_container['NetworkSettings']['IPAddress']
        self.redis_conn = redis.StrictRedis(host=hipache_ip, port=6379)
        try:
            self.redis_conn.get('__test__')
        except redis.exceptions.ConnectionError:
            print " => Cannot connect to redis."
            exit(1)

    def __find_hipache_container(self):
        containers = self.dockercli.get_running_containers()
        for container in containers:
            container_info = container['more_info']
            environments = container_info['Config']['Env']
            for environment in environments:
                if environment == 'HIPACHE=__INSTANCE__':
                    return container_info
        print " => Cannot find hipache container."
        exit(1)


    def list_redis_items(self):
        redis_items = []
        frontends = self.redis_conn.keys()
        for item in frontends:
            if 'frontend:' not in item:
                continue

            item_url = {}
            item_url['url'] = item[9:]
            item_url['data'] = []

            ips = self.redis_conn.lrange(item, 0, -1)

            if len(ips):
                item_url['metadata'] = ips[0]

            for ip in ips[1:]:
                item_url['data'].append(ip)

            redis_items.append(item_url)
        return redis_items

    def list_redis_items_flatten(self):
        new_arr = []
        for item in self.list_redis_items():
            url = item['url']
            for item_data in item['data']:
                new_map = {}
                parser = Parser()
                err, protocol, ip, port = parser.split_ipdata(item_data)
                new_map['url'] = url
                new_map['protocol'] = protocol
                new_map['ip'] = ip
                new_map['host_port'] = port
                new_map['original'] = item_data
                new_arr.append(new_map)
        return new_arr
