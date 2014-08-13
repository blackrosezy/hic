Hipache Cli
===========

Hipache cli (hic) is a program written in Go to add and remove url in Hipache.


### Basic Usage:

Show help file
```bash
hic
```

List all mappings
```bash
hic show
```

Add a url mapping
```bash
hic add http://mywebsite.com 192.168.1.1:3475
```

Remove 192.168.1.1:3475 from mywebsite.com
```bash
hic remove http://mywebsite.com 192.168.1.1:3475
```

Remove mywebsite.com
```bash
hic remove http://mywebsite.com
```


### Prerequisite

You need to install Go to compile this program. You may skip this section if you have working Go installation in your machine. E.g in Ubuntu :
```bash
apt-get install golang
```
Create directory for Go source codes and binaries:
```bash
mkdir -p ~/.gocode
```
set GOPATH environment:
```bash
vi ~/.bashrc
```
...and add this export command at the bottom of the file:
```bash
export GOPATH=~/.gocode
```
Reload .bashrc
```bash
. ~/.bashrc
```


### Compile

To compile Hipache Cli, you need to download dependency library first:
```bash
go get github.com/garyburd/redigo/redis
```
Once you finish download you can checkout Hipache Cli source code and compile:
```bash
git clone git@github.com:blackrosezy/hipache-cli.git
cd hipache-cli
go build hic.go
chmod +x hic
```

You can test by typing:
```bash
./hic
```
If you can see a help instruction, then you build is success.


### Install

Copy the Hipache Cli to /usr/local/bin
```bash
cp hic /usr/local/bin
```
Now you can access Hipache Cli from any path.


### License

Hipache Cli is licensed under the MIT license. (http://opensource.org/licenses/MIT)


### Contributing

I you have better idea or something to share, your pull request are welcome!
