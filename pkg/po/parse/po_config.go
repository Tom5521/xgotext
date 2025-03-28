package parse

import (
	"io"
	"log"
)

type PoConfig struct {
	lastCfg any // Any type to not refer itself.

	Logger          *log.Logger
	SkipHeader      bool
	CleanDuplicates bool
	IgnoreComments  bool
}

func (p *PoConfig) RestoreLastCfg() {
	if p.lastCfg != nil {
		*p = p.lastCfg.(PoConfig)
	}
}

func (p *PoConfig) ApplyOptions(opts ...PoOption) {
	p.lastCfg = *p
	for _, opt := range opts {
		opt(p)
	}
}

func DefaultPoConfig(opts ...PoOption) PoConfig {
	c := PoConfig{
		Logger:          log.New(io.Discard, "", 0),
		CleanDuplicates: true,
	}

	c.ApplyOptions(opts...)

	return c
}

type PoOption func(*PoConfig)

func PoWithConfig(cfg PoConfig) PoOption {
	return func(c *PoConfig) { *c = cfg }
}

func PoWithSkipHeader(s bool) PoOption {
	return func(c *PoConfig) { c.SkipHeader = s }
}

func PoWithCleanDuplicates(cd bool) PoOption {
	return func(c *PoConfig) { c.CleanDuplicates = cd }
}

func PoWithLogger(logger *log.Logger) PoOption {
	return func(c *PoConfig) { c.Logger = logger }
}
