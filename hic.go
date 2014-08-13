package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"strings"
)

func remove(c redis.Conn, domain string, ip string) {
	// get all keys
	var domains, _ = redis.Strings(c.Do("KEYS", "*"))

	// for each keys(frontend:*)
	for _, tmp_domain := range domains {

		protocol := "http://"
		new_domain := domain

		if strings.Contains(domain, "http://") {
			new_domain = strings.Replace(domain, "http://", "", -1)
		}

		if strings.Contains(domain, "https://") {
			new_domain = strings.Replace(domain, "https://", "", -1)
			protocol = "https://"
		}

		key := fmt.Sprintf("frontend:%s", new_domain)

		// check requested domain same with current domain or not
		if key == tmp_domain {

			if ip != "" { //if specific ip was specify
				// for each value in current domain
				ips, _ := redis.Strings(c.Do("LRANGE", key, 0, -1))
				for _, local_ip := range ips {
					value := fmt.Sprintf("%s%s", protocol, ip) // 'http://' + ip

					// if value match with requested ip, then delete
					if value == local_ip {
						c.Do("LREM", key, 0, value)
					}
				}
				//if no specific ip specify
			} else {
				// for each value in current domain
				ips, _ := redis.Strings(c.Do("LRANGE", key, 0, -1))
				for _, local_ip := range ips {
					c.Do("LREM", key, 0, local_ip)
				}
			}

			break // break early if domain found
		}

	}
}

func add(c redis.Conn, domain string, ip string) {
	protocol := "http://"
	new_domain := domain

	if strings.Contains(domain, "http://") {
		new_domain = strings.Replace(domain, "http://", "", -1)
	}

	if strings.Contains(domain, "https://") {
		new_domain = strings.Replace(domain, "https://", "", -1)
		protocol = "https://"
	}

	remove(c, domain, ip)

	key := fmt.Sprintf("frontend:%s", new_domain)

	title, _ := redis.Strings(c.Do("LRANGE", key, 0, 0))
	if len(title) != 0 {
		fmt.Println("adsfasfsdaf")
		if title[0] != "_"+new_domain {
			c.Do("RPUSH", key, "_"+new_domain)
		}
	} else {
		c.Do("RPUSH", key, "_"+new_domain)
	}

	new_ip := protocol + ip
	c.Do("RPUSH", key, new_ip)

	fmt.Println(" => Added!")
}

func show(c redis.Conn) {
	var domains, _ = redis.Strings(c.Do("KEYS", "*"))
	fmt.Println("|--------------------------------------------")
	fmt.Printf("| %-20s        %-20s\n", "Domain", "IP")
	fmt.Println("|--------------------------------------------")
	for _, domain := range domains {
		clean_domain := strings.Replace(domain, "frontend:", "", -1)

		ips, _ := redis.Strings(c.Do("LRANGE", domain, 0, -1))
		for _, ip := range ips {
			fmt.Printf("| %-20s   --   %-20s\n", clean_domain, ip)
		}

	}
	fmt.Println("|--------------------------------------------")
}

func main() {

	var (
		c, err = redis.Dial("tcp", "127.0.0.1:6379")
	)

	if err != nil {
		log.Fatal(err)
	}

	argc := len(os.Args)
	argv := os.Args

	if argc == 2 {
		if argv[1] == "show" {
			show(c)
		}
	} else if argc == 3 {
		if argv[1] == "remove" {
			remove(c, argv[2], "")
			show(c)
		}
	} else if argc == 4 {
		if argv[1] == "add" {
			add(c, argv[2], argv[3])
			show(c)
		} else if argv[1] == "remove" {
			remove(c, argv[2], argv[3])
			show(c)
		}
	} else {
		fmt.Println("Hipache cli ver. 1.0")
		fmt.Println("   Usage:\n")

		fmt.Println("   List all mappings")
		fmt.Println("   hic show\n")

		fmt.Println("   Add a url mapping")
		fmt.Println("   hic add http://mywebsite.com 192.168.1.1:3475\n")

		fmt.Println("   Remove 192.168.1.1:3475 from mywebsite.com")
		fmt.Println("   hic remove http://mywebsite.com 192.168.1.1:3475\n")

		fmt.Println("   Remove mywebsite.com")
		fmt.Println("   hic remove http://mywebsite.com\n")
	}

}
