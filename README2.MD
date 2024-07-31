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
