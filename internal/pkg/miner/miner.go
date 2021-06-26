package miner

import (
	"github.com/mineway/excavator/internal/pkg/alg"
)

type Core struct {
	ID            int
	Alg           alg.Alg
	CurrPoolID    int
	IsRunning     bool
	IsWin         bool
	RigName       string
	Port          int
	ExpectedHr    int
	IsLogEnabled  bool
	StatsInterval int
}

func New() (c *Core) {
	return
}