# Vardius - pubsub

[![Build Status](https://travis-ci.org/vardius/pubsub.svg?branch=master)](https://travis-ci.org/vardius/pubsub)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/pubsub)](https://goreportcard.com/report/github.com/vardius/pubsub)
[![codecov](https://codecov.io/gh/vardius/pubsub/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/pubsub)
[![](https://godoc.org/github.com/vardius/pubsub?status.svg)](http://godoc.org/github.com/vardius/pubsub)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/pubsub/blob/master/LICENSE.md)

pubsub - gRPC message-oriented middleware on top of [message-bus](https://github.com/vardius/message-bus), event ingestion and delivery system.

<details>
  <summary>Table of Content</summary>

<!-- toc -->

- [About](#about)
- [How to use](#how-to-use)
  - [Docker](#docker)
    - [How to use this image](#how-to-use-this-image)
    - [Environment Variables](#environment-variables)
    - [Makefile](#makefile)
  - [Client](#client) - [Use in your Go project](#use-in-your-go-project)
<!-- tocstop -->

</details>

# ABOUT

Contributors:

- [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/pubsub/issues) to manage them.

# HOW TO USE

1. [GoDoc](http://godoc.org/github.com/vardius/pubsub)
2. [Examples](http://godoc.org/github.com/vardius/pubsub#pkg-examples)

## Docker

### How to use this image

Starting a pubsub instance:

```bash
docker run --name my-pubsub -e QUEUE_BUFFER_SIZE=100 -d vardius/pubsub:tag
```

### Environment Variables

#### `HOST`

This is optional variable, sets gRPC server host value. Default to `0.0.0.0`

#### `PORT`

This is optional variable, sets gRPC server port value. Default to `9090`

#### `QUEUE_BUFFER_SIZE`

This is optional variable, sets message bus queue buffer size. Default to number of CPUs`

### Makefile

```sh
version                        Show version
docker-build                   Build given container. Example: `make docker-build`
docker-run                     Run container on given port. Example: `make docker-run PORT=3000`
docker-stop                    Stop docker container. Example: `make docker-stop`
docker-rm                      Stop and then remove docker container. Example: `make docker-rm`
docker-publish                 Docker publish. Example: `make docker-publish REGISTRY=https://your-registry.com`
docker-tag                     Tag current container. Example: `make docker-tag REGISTRY=https://your-registry.com`
docker-release                 Docker release - build, tag and push the container. Example: `make docker-release REGISTRY=https://your-registry.com`
```

## Client

### Use in your Go project

#### Publish

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	pubsub_proto "github.com/vardius/pubsub/proto"
)

func main() {
    host:= "0.0.0.0"
    port:= 9090
    ctx := context.Background()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             20 * time.Second, // wait 20 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}),
    }

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		os.Exit(1)
    }
    defer conn.Close()

	client := pubsub_proto.NewMessageBusClient(pubsubConn)

    client.Publish(ctx, &pubsub_proto.PublishRequest{
		Topic:   "my-topic",
		Payload: []byte("Hello you!"),
    })
}
```

#### Subscribe

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	pubsub_proto "github.com/vardius/pubsub/proto"
)

func main() {
    host:= "0.0.0.0"
    port:= 9090
    ctx := context.Background()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             20 * time.Second, // wait 20 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}),
    }

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		os.Exit(1)
    }
    defer conn.Close()

	client := pubsub_proto.NewMessageBusClient(pubsubConn)

	stream, err := client.Subscribe(ctx, &pubsub_proto.SubscribeRequest{
		Topic: "my-topic",
	})
	if err != nil {
		os.Exit(1)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
		    os.Exit(1) // stream closed or error
		}

		fmt.Println(resp.GetPayload())
	}
}
```

## License

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
