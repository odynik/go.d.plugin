package suricata

import (
	"math/rand"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("suricata", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery:        module.UpdateEvery,
			AutoDetectionRetry: module.AutoDetectionRetry,
			Priority:           module.Priority,
			Disabled:           true,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Suricata {
	return &Suricata{
		Config: Config{
			Charts: ConfigCharts{
				Num:  1,
				Dims: 4,
			},
			HiddenCharts: ConfigCharts{
				Num:  0,
				Dims: 4,
			},
		},

		randInt:       func() int64 { return rand.Int63n(100) },
		collectedDims: make(map[string]bool),
	}
}

type (
	Config struct {
		Charts       ConfigCharts `yaml:"charts"`
		HiddenCharts ConfigCharts `yaml:"hidden_charts"`
	}
	ConfigCharts struct {
		Type string `yaml:"type"`
		Num  int    `yaml:"num"`
		Dims int    `yaml:"dimensions"`
	}
)

type Suricata struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`

	randInt       func() int64
	charts        *module.Charts
	collectedDims map[string]bool
}

func (e *Suricata) Init() bool {
	err := e.validateConfig()
	if err != nil {
		e.Errorf("config validation: %v", err)
		return false
	}

	charts, err := e.initCharts()
	if err != nil {
		e.Errorf("charts init: %v", err)
		return false
	}
	e.charts = charts
	return true
}

func (e *Suricata) Check() bool {
	return len(e.Collect()) > 0
}

func (e *Suricata) Charts() *module.Charts {
	return e.charts
}

func (e *Suricata) Collect() map[string]int64 {
	mx, err := e.collect()
	if err != nil {
		e.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (Suricata) Cleanup() {}
