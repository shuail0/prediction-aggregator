package opinion

import "github.com/shuail0/prediction-aggregator/pkg/exchange"

func init() {
	exchange.Register("opinion", func() (exchange.Exchange, error) {
		return New(Config{})
	})
}
