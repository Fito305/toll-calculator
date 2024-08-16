# toll-calculator

```

focker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAKFA_CFG_AUTO_CREATE_TOPIC
kafka:latest


```
## The one I have been using

docker-compose up -d --force-recreate

In docker -d means deamon which runs the container in the background.
```

```

```

```

## Istalling protobu compiler
For linux user or (WSL2)

sudo apt install -y protobuf-compiler

```

For Mac users you can use Brew for this
```

brew install protobuff
```

## Installing GRPC and Protobuffer plugins for Golang
1. Protobuffers
```

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

2. GRPC
```

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```


3. NOTE that you need to set the /go/bin directory in your path
```

PATH="${PATH}:${HOME}/go/bin"
```

4. Install the package dependencies
```
4.1 protobuffer package
```

go get google.golang.org/protobuf

```

4.2 
grpc package

```
go get google.golang.org/grpc

```

## Installing Prometheus
Install Prometheus in a Docker container
```

docker run -p 9090:9090 -v ./.config/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

Installing prometheus golang client
```
go get github.com/prometheus/client_golang/prometheus
```


Installing Prometheus natively on you system
1. Clone the repository
```

git clone https://github.com/prometheus/prometheus.git
```

2. Insatll
```
cd prometheus
make build
```

3. Run the Prometheus daemon
```
./promtheus --config.file<your_config_file>yml
```

4. In the project case that would be (running from inside the project directory)
```

../prometheus/prometheus --config.file=.config/prometheus.yml
```

Use .env enviroment variables
```

https:// <--- delete the http part when using `go get`
go get github.com/joho/godotenv

```
