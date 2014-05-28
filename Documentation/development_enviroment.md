# How to install the development enviroment

#####1. Install GO 1.2+  
You can find help in the official documentation: [http://golang.org/doc/install#tarball](http://golang.org/doc/install#tarball)
```
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PATH:/usr/local/go/bin
cd
mkdir go
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
```

#####2. Install zeromq 3.2  
CentOS: [http://zeromq.org/distro:centos](http://zeromq.org/distro:centos)  
Ubuntu: [https://launchpad.net/~chris-lea/+archive/zeromq](https://launchpad.net/~chris-lea/+archive/zeromq)

#####3. Get etcd
```
go get github.com/coreos/etcd
cd $GOPATH/src/github.com/coreos/etcd
git checkout v0.3.0
./build
```

#####4. Get goven
```
go get github.com/kr/goven
```

#####5. Get Hydra
```
go get github.com/innotech/hydra
cd /home/innotechdev/go/src/github.com/innotech/hydra
git checkout 3.0.0
cd vendors
goven goven github.com/coreos/etcd
```

#####6. Build Hydra
```
# Go to hydra parent directory
cd ..
./build2
```
