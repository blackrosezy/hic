Hic
===========

Hic is a simple tool to add, remove and sync urls in redis (Hipache server default backend)

```
URL                     IP            HOST PORT    CONT. PORT  CONTAINER NAME
----------------------  ----------  -----------  ------------  -----------------------------
http://a.example.com    172.17.0.3         8081          8081  container_A
http://b.example.com    172.17.0.b         4567            80  container_B
```

### Usage:

**Add**

Add an url to redis server.
```
hic add <url> <container_name> [<container_port>]
```
e.g: `hic add http://a.example.com container_A`

e.g: `hic add http://a.example.com container_A 8081`

**Remove**

Remove an url from redis server.
```
hic rm <url>
```
e.g: `hic rm http://a.example.com`

**Sync**

Sync redis server with configuration file. The default location for configuration file is `~/.hic.json`.
```
  hic add sync
```

### Installation

```bash
pip install https://github.com/blackrosezy/hic/archive/master.zip
```

### License

Hic is licensed under the MIT license. (http://opensource.org/licenses/MIT)


### Contributing

I you have better idea or something to share, your pull request are welcome!
