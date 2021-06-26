package alg

import (
	"context"
)

const (
	GPU = 1
	CPU = 2
)

type Alg interface {
	Name() string
	Type() int
	GetServer() []Server
	Run(ctx context.Context, userAddr string, withSSL bool) error
	Stop(ctx context.Context) error
}

type Server struct {
	Addr 	string
	StratumPort 	int
	AltStratumPort 	int
	SSLPort 		int
}