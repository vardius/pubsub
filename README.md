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
  - [Client](https://github.com/vardius/pubsub/tree/master/proto#client)
  	- [Use in your Go project](https://github.com/vardius/pubsub/tree/master/proto#use-in-your-go-project)
	  - [Publish](https://github.com/vardius/pubsub/tree/master/proto#publish)
	  - [Subscribe](https://github.com/vardius/pubsub/tree/master/proto#subscribe)
  - [Protocol Buffers](https://github.com/vardius/pubsub/tree/master/proto#protocol-buffers)
	- [Generating client and server code](https://github.com/vardius/pubsub/tree/master/proto#generating-client-and-server-code)
<!-- tocstop -->
</details>

# ABOUT

Contributors:

- [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/pubsub/issues) to manage them.

# HOW TO USE

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

See [proto package](https://github.com/vardius/pubsub/blob/master/proto) for details.

## License

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
