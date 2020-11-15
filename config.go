package commonlogger

import (
	"github.com/johnmackenzie91/commonlogger/resolvers"
)

type Config struct {
	ContextResolver  resolvers.ContextResolver
	RequestResolver  resolvers.RequestResolver
	ResponseResolver resolvers.ResponseResolver
}

type Option func(*Config)

func WithContextResolver(f resolvers.ContextResolver) Option {
	return func(c *Config) {
		c.ContextResolver = f
	}
}

func WithRequestResolver(f resolvers.RequestResolver) Option {
	return func(c *Config) {
		c.RequestResolver = f
	}
}

func WithResponseResolver(f resolvers.ResponseResolver) Option {
	return func(c *Config) {
		c.ResponseResolver = f
	}
}
