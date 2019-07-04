package main

import (
	"runtime"
	"time"

	"github.com/caarlos0/env/v6"
)

// Env stores environment values
var Env *environment

type environment struct {
	// Host gRPC tcp host value. Default 0.0.0.0
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	// gRPC tcp port value. Default 9090
	Port int `env:"PORT" envDefault:"9090"`
	// QueueSize sets buffered channel length per subscriber. Default 0, which evaluates to runtime.NumCPU().
	QueueSize int `env:"QUEUE_BUFFER_SIZE" envDefault:"0"`
	// KeepaliveEnforcementPolicyMinTime (nanoseconds) if a client pings more than once every 5 minutes (default), terminate the connection.
	KeepaliveEnforcementPolicyMinTime time.Duration `env:"KEEPALIVE_MIN_TIME" envDefault:"5m"`
	// KeepaliveParamsTime (nanoseconds)  Ping the client if it is idle for 2 hours (default) to ensure the connection is still active.
	KeepaliveParamsTime time.Duration `env:"KEEPALIVE_TIME" envDefault:"2h"`
	// KeepaliveParamsTimeout (nanoseconds)  Wait 20 second (default) for the ping ack before assuming the connection is dead.
	KeepaliveParamsTimeout time.Duration `env:"KEEPALIVE_TIMEOUT" envDefault:"20s"`
	// VerboseLevel level. -1 = Disabled, 0 = Critical, 1 = Error, 2 = Warning, 3 = Info, 4 = Debug. Default 4.
	VerboseLevel int `env:"LOG_VERBOSE_LEVEL" envDefault:"3"`
}

func init() {
	Env = &environment{}
	env.Parse(Env)

	if Env.QueueSize == 0 {
		Env.QueueSize = runtime.NumCPU()
	}
}
