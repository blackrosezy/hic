package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/olekukonko/tablewriter"
	"github.com/samalba/dockerclient"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	CONFIG_FILENAME = ".hic"
)

var (
	ContainerNotFound = errors.New("Container not found.")
	PortNotFound      = errors.New("Cannot find port.")
	IpNotFound        = errors.New("Cannot find Ip.")
)

type MappingType struct {
	url  string
	name string
	port int
}

type RedisType struct {
	url   string
	ip    string
	port  int
	key   string
	value string
}

func getContainerIpByName(docker *dockerclient.DockerClient, container_name string) (string, error) {
	containers, err := docker.ListContainers(false)
	if err != nil {
		log.Fatal(err)
	}

	container_id := "0"
	for _, c := range containers {
		names := c.Names
		for _, name := range names {
			if name == "/"+container_name {
				container_id = c.Id
				break
			}
		}
		if container_id != "0" {
			break
		}
	}

	if container_id == "0" {
		return "", ContainerNotFound
	}

	info, err := docker.InspectContainer(container_id)
	if err != nil {
		return "", err
	}

	container_ip := info.NetworkSettings.IpAddress

	return container_ip, nil
}

func getContainerIpByPort(docker *dockerclient.DockerClient, container_port int) (string, error) {
	containers, err := docker.ListContainers(false)
	if err != nil {
		log.Fatal(err)
	}

	container_id := "0"
	for _, c := range containers {
		ports := c.Ports
		for _, p := range ports {
			if p.PrivatePort == container_port {
				container_id = c.Id
				break
			}
		}
		if container_id != "0" {
			break
		}
	}

	if container_id == "0" {
		return "", ContainerNotFound
	}

	info, err := docker.InspectContainer(container_id)
	if err != nil {
		return "", err
	}

	container_ip := info.NetworkSettings.IpAddress

	return container_ip, nil
}

func getPortFromMixIp(mix_ip string) (int, error) {
	ip_components := strings.Split(mix_ip, ":")
	port := 0
	if len(ip_components) == 3 {
		port_s := ip_components[2]
		port_tmp, err := strconv.Atoi(port_s)
		if err != nil {
			log.Fatal(err)
		}
		port = port_tmp
	} else {
		return 0, PortNotFound
	}
	return port, nil
}

func getUrlFromMixKey(mix_ip string) (string, error) {
	url_components := strings.Split(mix_ip, "frontend:")
	url := ""
	if len(url_components) == 2 {
		url = url_components[1]
	} else {
		return "", IpNotFound
	}
	return url, nil
}

func getIpFromMixIp(mix_ip string) (string, error) {
	ip_components := strings.Split(mix_ip, ":")
	ip := ""
	if len(ip_components) == 3 {
		ip = ip_components[1]
		ip = strings.Replace(ip, "//", "", -1)
	} else {
		return "", IpNotFound
	}
	return ip, nil
}

func getRedisDataAsList(c redis.Conn) ([]RedisType, error) {
	var list []RedisType

	data, err := redis.Strings(c.Do("KEYS", "*"))
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range data {
		sub_data, _ := redis.Strings(c.Do("LRANGE", item, 0, -1))
		url, _ := getUrlFromMixKey(item)
		for _, sub_item := range sub_data {
			var i RedisType

			i.url = url

			ip, _ := getIpFromMixIp(sub_item)
			i.ip = ip

			port, _ := getPortFromMixIp(sub_item)
			i.port = port

			i.key = item
			i.value = sub_item

			list = append(list, i)
		}
	}

	return list, nil
}

func _Remove(c redis.Conn, query RedisType, remove_type int) {
	list, err := getRedisDataAsList(c)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range list {
		if remove_type == 1 {
			if item.ip == query.ip && item.url == query.url && item.port == query.port {
				key := fmt.Sprintf("frontend:%s", item.url)

				protocol := "http://"
				if item.port == 443 {
					protocol = "https://"
				}

				value := protocol + item.ip + ":" + strconv.Itoa(item.port)

				c.Do("LREM", key, 0, value)
			}
		} else if remove_type == 2 {
			if item.ip == query.ip && item.url == query.url {
				key := fmt.Sprintf("frontend:%s", item.url)

				protocol := "http://"
				if item.port == 443 {
					protocol = "https://"
				}

				value := protocol + item.ip + ":" + strconv.Itoa(item.port)

				c.Do("LREM", key, 0, value)
			}
		} else if remove_type == 3 {
			if item.url == query.url {
				key := fmt.Sprintf("frontend:%s", item.url)

				protocol := "http://"
				if item.port == 443 {
					protocol = "https://"
				}

				value := protocol + item.ip + ":" + strconv.Itoa(item.port)
				c.Do("LREM", key, 0, value)
			}
		}
	}

	list, err = getRedisDataAsList(c)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, item := range list {
		if item.url == query.url {
			count = count + 1
		}
	}

	if count == 1 {
		key := fmt.Sprintf("frontend:%s", query.url)
		items, _ := redis.Strings(c.Do("LRANGE", key, 0, -1))
		for _, item := range items {
			c.Do("LREM", key, 0, item)
		}
	}
}

func _Add(c redis.Conn, query RedisType) {
	_Remove(c, query, 1)

	key := fmt.Sprintf("frontend:%s", query.url)
	title, _ := redis.Strings(c.Do("LRANGE", key, 0, 0))
	if len(title) != 0 {
		if title[0] != "_"+query.url {
			c.Do("RPUSH", key, "_"+query.url)
		}
	} else {
		c.Do("RPUSH", key, "_"+query.url)
	}

	protocol := "http://"
	if query.port == 443 {
		protocol = "https://"
	}

	new_ip := protocol + query.ip + ":" + strconv.Itoa(query.port)
	c.Do("RPUSH", key, new_ip)
}

func Clear(c redis.Conn) {
	data, err := getRedisDataAsList(c)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range data {
		c.Do("LREM", item.key, 0, item.value)
	}
}

// c, <url>, <ip>, <private port>
func Remove(docker *dockerclient.DockerClient, c redis.Conn, url string, ip string, port int) {
	new_url := url
	if strings.HasPrefix(url, "http://") {
		new_url = strings.Replace(url, "http://", "", -1)
	}

	if strings.HasPrefix(url, "https://") {
		new_url = strings.Replace(url, "https://", "", -1)
	}

	var query RedisType
	query.port = port
	query.ip = ip
	query.url = new_url

	if ip != "" && port != 0 {
		_Remove(c, query, 1)
	} else if ip != "" && port == 0 {
		_Remove(c, query, 2)
	} else if ip == "" && port == 0 {
		_Remove(c, query, 3)
	}

	fmt.Println(" => Removed!")
}

// c, <container name/ip>, <url>, <private port>
func Add(docker *dockerclient.DockerClient, c redis.Conn, container_name string, url string, port int) {
	new_url := url
	if strings.HasPrefix(url, "http://") {
		new_url = strings.Replace(url, "http://", "", -1)
	}

	if strings.HasPrefix(url, "https://") {
		new_url = strings.Replace(url, "https://", "", -1)
	}

	ip, err := getContainerIpByName(docker, container_name)
	if err != nil {
		ip = container_name // user pass ip value in container_name parameter
	}

	var query RedisType
	query.port = port
	query.ip = ip
	query.url = new_url

	_Add(c, query)

	fmt.Println(" => Added!")
}

func Show(c redis.Conn) {
	data, err := getRedisDataAsList(c)
	if err != nil {
		log.Fatal(err)
	}

	rows := [][]string{}
	for _, item := range data {
		i := []string{}
		i = append(i, item.url)

		if item.ip != "" {
			i = append(i, item.ip)
		} else {
			i = append(i, "IDENTIFIER("+item.value+")")
		}

		if item.port != 0 {
			i = append(i, strconv.Itoa(item.port))
		} else {
			i = append(i, "-")
		}

		rows = append(rows, i)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Url", "Ip", "Port"})

	for _, v := range rows {
		table.Append(v)
	}
	table.SetColumnSeparator(" ")
	table.Render()
}

func Help() {
	fmt.Println("Hic (Hipache cli) ver. 2.0\n")
	fmt.Println("   Usage:\n")

	fmt.Println("   = List")
	fmt.Println("   ==============================")
	fmt.Println("   List all mappings.")
	fmt.Println("         hic\n")

	fmt.Println("   + Add")
	fmt.Println("   ==============================")
	fmt.Println("   No need to insert http or https for url...")
	fmt.Println("   ..., the protocol is detected by port number. If port...")
	fmt.Println("   ... number not given, default is 80.\n")

	fmt.Println("   Add a mapping by container name")
	fmt.Println("         hic add <container name> <url> <private port>")
	fmt.Println("    e.g. hic add blog_web_1 mywebsite.com 80\n")

	fmt.Println("   Add a mapping by ip")
	fmt.Println("         hic add <ip> <url> <private port>")
	fmt.Println("    e.g. hic add 192.168.1.6 mywebsite.com 80\n")

	fmt.Println("   - Remove")
	fmt.Println("   ==============================")
	fmt.Println("   Remove url(s).")
	fmt.Println("         hic rm <url>")
	fmt.Println("    e.g. hic rm mywebsite.com\n")

	fmt.Println("   Remove url(s) by ip. This will result...")
	fmt.Println("   ...removing all ports numbers (80, 443, etc.).")
	fmt.Println("         hic rm <url> <ip>")
	fmt.Println("    e.g. hic rm mywebsite.com 192.168.1.6\n")

	fmt.Println("   Remove an url mapping by container name and port.")
	fmt.Println("         hic rm <url> <ip> <private port>")
	fmt.Println("    e.g. hic rm mywebsite.com 192.168.1.6 80\n")

	fmt.Println("   <> Synchronization")
	fmt.Println("   ==============================")
	fmt.Println("   Sync all Ips between containers and Redis in Haraka.")
	fmt.Println("         hic sync\n")
}
func Sync(c redis.Conn) {

}
func main() {
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		log.Fatal("Cannot connect to Docker API.")
	}

	hipache_ip, err := getContainerIpByPort(docker, 6379)
	if err != nil {
		log.Fatal(err)
	}

	// Yeah, let's connect to Hipache server!
	c, err := redis.Dial("tcp", hipache_ip+":6379")
	if err != nil {
		log.Fatal(err)
	}

	argc := len(os.Args)
	argv := os.Args

	if argc == 1 {
		Show(c)
	} else if argc == 2 {
		if argv[1] == "sync" {
			Sync(c)
			Show(c)
		} else if argv[1] == "clear" {
			Clear(c)
			Show(c)
		} else {
			Help()
		}
	} else if argc == 3 {
		if argv[1] == "rm" {
			// remove(docker, c, <url>, "", 0)
			Remove(docker, c, argv[2], "", 0)
			Show(c)
		} else {
			Help()
		}
	} else if argc == 4 {
		if argv[1] == "add" {
			// add(docker, c, <container name/ip>, <url>, 80)
			Add(docker, c, argv[2], argv[3], 80)
			Show(c)
		} else if argv[1] == "rm" {
			// remove(docker, c, <url>, <ip>, 0)
			Remove(docker, c, argv[2], argv[3], 0)
			Show(c)
		} else {
			Help()
		}
	} else if argc == 5 {
		if argv[1] == "add" {
			port, err := strconv.Atoi(argv[4])
			if err != nil {
				log.Fatal(err)
			}
			// add(docker, c, <container name/ip>, <url>, <private port>)
			Add(docker, c, argv[2], argv[3], port)
			Show(c)
		} else if argv[1] == "rm" {
			port, err := strconv.Atoi(argv[4])
			if err != nil {
				log.Fatal(err)
			}
			// remove(docker, c, <url>, <ip>, <private port>)
			Remove(docker, c, argv[2], argv[3], port)
			Show(c)
		} else {
			Help()
		}
	} else {
		Help()
	}

}
