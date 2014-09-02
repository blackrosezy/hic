Hic
===========

Hic is a cli program written in Go to add and remove url in Hipache.

```
+-------------------------+--------------------------------------+------+
            URL                              IP                    PORT
+-------------------------+--------------------------------------+------+
  morzsoftware.com          IDENTIFIER(_morzsoftware.com)          -
  morzsoftware.com          192.168.1.4                            80
  morzsoftware.com          192.168.1.4                            443
  morzproject.com           IDENTIFIER(_morzproject.com)           -
  morzproject.com           192.168.1.5                            80
  morzproject.com           192.168.1.5                            443
  test.morzproject.com      IDENTIFIER(_test.morzproject.com)      -
  test.morzproject.com      172.17.0.74                            3186
  billing.morzproject.com   IDENTIFIER(_billing.morzproject.com)   -
  billing.morzproject.com   172.17.0.74                            5000
+-------------------------+--------------------------------------+------+
```

### Usage:

```
Hic (Hipache cli) ver. 2.0

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
         hic add <container name> <url> <private port>
    e.g. hic add blog_web_1 mywebsite.com 80

   Add a mapping by ip
         hic add <ip> <url> <private port>
    e.g. hic add 192.168.1.6 mywebsite.com 80

   - Remove
   ==============================
   Remove url(s).
         hic rm <url>
    e.g. hic rm mywebsite.com

   Remove url(s) by ip. This will result...
   ...removing all ports numbers (80, 443, etc.).
         hic rm <url> <ip>
    e.g. hic rm mywebsite.com 192.168.1.6

   Remove an url mapping by container name and port.
         hic rm <url> <ip> <private port>
    e.g. hic rm mywebsite.com 192.168.1.6 80

   <> Synchronization
   ==============================
   Sync all Ips between containers and Redis in Haraka.
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
