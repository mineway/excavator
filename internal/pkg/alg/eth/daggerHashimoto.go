package eth

import (
	"context"
	"github.com/mineway/excavator/internal/pkg/alg"
)

type DaggerHashimoto struct {}

func (DaggerHashimoto) Name() string {
	return "Dagger Hashimoto"
}

func (DaggerHashimoto) Type() int {
	return alg.GPU
}

func (DaggerHashimoto) GetServer() []alg.Server {
	return []alg.Server{
		{
			Addr: "eu1.ethermine.org",
			StratumPort: 4444,
			AltStratumPort: 14444,
			SSLPort: 5555,
		},
		{
			Addr: "asia1.ethermine.org",
			StratumPort: 4444,
			AltStratumPort: 14444,
			SSLPort: 5555,
		},
		{
			Addr: "us1.ethermine.org",
			StratumPort: 4444,
			AltStratumPort: 14444,
			SSLPort: 5555,
		},
		{
			Addr: "us2.ethermine.org",
			StratumPort: 4444,
			AltStratumPort: 14444,
			SSLPort: 5555,
		},
	}
}

func (DaggerHashimoto) Run(ctx context.Context, userAddr string, withSSL bool) error {
	return nil
}

func (DaggerHashimoto) Stop(ctx context.Context) error {
	return nil
}
