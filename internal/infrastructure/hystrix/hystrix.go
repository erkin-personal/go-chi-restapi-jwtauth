package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
)

func ConfigureCircuitBreaker() {
	hystrix.ConfigureCommand("postgres", hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  10,
		ErrorPercentThreshold:  50,
		RequestVolumeThreshold: 5,
		SleepWindow:            5000,
	})
}
