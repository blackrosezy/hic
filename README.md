Hic
===========

Hic is a standalone cli program written in Go to add and remove url in Hipache. The motivation why I wrote this program is because the lack of cli program to add/remove url in [Hipache](https://github.com/hipache/hipache). My focus to this program is to make sure it have **small memory usage**, **fast** and **no dependencies**.

```
+-------------------------+--------------+------+--------------------------+
            URL                  IP        PORT        CONTAINER NAME
+-------------------------+--------------+------+--------------------------+
  www.example.com           172.17.1.244   80     example_web_1
  www.mywebsite.com         172.17.2.35    80     mywebsite_web_1
  example.com               172.17.1.244   80     example_web_1
  fun.morzproject.com       172.17.2.55    80     funproject_web_1
  mywebsite.com             172.17.2.35    80     mywebsite_web_1
+-------------------------+--------------+------+--------------------------+
```

### Usage:

```
Hic (Hipache cli) ver. 2.1

   Usage:

   = List
   ==============================
   List all mappings.
         hic

   + Add
   ==============================
   No need to insert http or https for url...
   ..., the protocol is detected by port number. If port...
   ... number not given, default is 80.

   Add a mapping by container name
         hic add <url> <container name> <private port>
    e.g. hic add mywebsite.com blog_web_1 80

   Add a mapping by ip
         hic add <url> <ip> <private port>
    e.g. hic add mywebsite.com 192.168.1.6 80

   - Remove
   ==============================
   Remove url(s).
         hic rm <url>
    e.g. hic rm mywebsite.com

   Remove url(s) by ip. This will result...
   ...removing all ports numbers (80, 443, etc.).
         hic rm <url> <ip>
    e.g. hic rm mywebsite.com 192.168.1.6

   Remove url(s) mapping by container name and port.
         hic rm <url> <ip> <private port>
    e.g. hic rm mywebsite.com 192.168.1.6 80

   <> Synchronization
   ==============================
   Sync all IPs between Docker containers and Redis server.
         hic sync
```

### Install

First, download and compile Hic source code:
```bash
go get -u github.com/blackrosezy/hic
go build github.com/blackrosezy/hic
```

Then, copy the Hic binary to /usr/local/bin
```bash
cp hic /usr/local/bin
chmod +x /usr/local/bin/hic
```

Now you can access Hic from any path.


### License

Hic is licensed under the MIT license. (http://opensource.org/licenses/MIT)


### Contributing

I you have better idea or something to share, your pull request are welcome!
