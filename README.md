# Vardius - pubsub

[![Build Status](https://travis-ci.org/vardius/pubsub.svg?branch=master)](https://travis-ci.org/vardius/pubsub)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/pubsub)](https://goreportcard.com/report/github.com/vardius/pubsub)
[![codecov](https://codecov.io/gh/vardius/pubsub/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/pubsub)
[![](https://godoc.org/github.com/vardius/pubsub?status.svg)](http://godoc.org/github.com/vardius/pubsub)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/pubsub/blob/master/LICENSE.md)

pubsub - gRPC message-oriented middleware on top of [message-bus](https://github.com/vardius/message-bus), event ingestion and delivery system.

<details>
  <summary>Table of Contents</summary>

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

## [Docker](https://hub.docker.com/r/vardius/pubsub)

### How to use this image

Starting a pubsub instance:

```bash
docker run --name my-pubsub -e QUEUE_BUFFER_SIZE=100 -d vardius/pubsub:tag
```

### Environment Variables

#### `HOST` (string)

This is optional variable, sets gRPC server host value. **Default `0.0.0.0`**

#### `PORT` (int)

This is optional variable, sets gRPC server port value. **Default `9090`**

#### `QUEUE_BUFFER_SIZE` (int)

This is optional variable, sets buffered channel length per subscriber. **Default 0**, which evaluates to `runtime.NumCPU()`.

#### `KEEPALIVE_MIN_TIME`

This is optional variable, if a client pings more than once every **5 minutes (default)**, terminate the connection.
ParseDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"

#### `KEEPALIVE_TIME` (nanoseconds)

This is optional variable, ping the client if it is idle for **2 hours (default)** to ensure the connection is still active.
ParseDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms",

#### `KEEPALIVE_TIMEOUT` (nanoseconds)

This is optional variable, wait **20 second (default)** for the ping ack before assuming the connection is dead.
ParseDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms",

#### `LOG_VERBOSE_LEVEL` (int)

This is optional variable, Verbose level. `-1` = Disabled, `0` = Critical, `1` = Error, `2` = Warning, `3` = Info, `4` = Debug. **Default 3 (Info)**.

### Makefile

```sh
➜  pubsub git:(master) make help
version                        Show version
docker-build                   Build given container. Example: `make docker-build`
docker-run                     Run container on given port. Example: `make docker-run PORT=9090`
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
