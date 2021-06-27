package suricata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// We want to ensure that module is a reference type, nothing more.

	assert.IsType(t, (*Suricata)(nil), New())
}

func TestSuricata_Init(t *testing.T) {
	// 'Init() bool' initializes the module with an appropriate config, so to test it we need:
	// - provide the config.
	// - set module.Config field with the config.
	// - call Init() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success on default config": {
			config: New().Config,
		},
		"success when only 'charts' set": {
			config: Config{
				Charts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
			},
		},
		"success when only 'hidden_charts' set": {
			config: Config{
				HiddenCharts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
			},
		},
		"success when 'charts' and 'hidden_charts' set": {
			config: Config{
				Charts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
				HiddenCharts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
			},
		},
		"fails when 'charts' and 'hidden_charts' set, but 'num' == 0": {
			wantFail: true,
			config: Config{
				Charts: ConfigCharts{
					Num:  0,
					Dims: 2,
				},
				HiddenCharts: ConfigCharts{
					Num:  0,
					Dims: 2,
				},
			},
		},
		"fails when only 'charts' set, 'num' > 0, but 'dimensions' == 0": {
			wantFail: true,
			config: Config{
				Charts: ConfigCharts{
					Num:  1,
					Dims: 0,
				},
			},
		},
		"fails when only 'hidden_charts' set, 'num' > 0, but 'dimensions' == 0": {
			wantFail: true,
			config: Config{
				HiddenCharts: ConfigCharts{
					Num:  1,
					Dims: 0,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			suricata := New()
			suricata.Config = test.config

			if test.wantFail {
				assert.False(t, suricata.Init())
			} else {
				assert.True(t, suricata.Init())
			}
		})
	}
}

func TestSuricata_Check(t *testing.T) {
	// 'Check() bool' reports whether the module is able to collect any data, so to test it we need:
	// - provide the module with a specific config.
	// - initialize the module (call Init()).
	// - call Check() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		prepare  func() *Suricata
		wantFail bool
	}{
		"success on default":                            {prepare: prepareSuricataDefault},
		"success when only 'charts' set":                {prepare: prepareSuricataOnlyCharts},
		"success when only 'hidden_charts' set":         {prepare: prepareSuricataOnlyHiddenCharts},
		"success when 'charts' and 'hidden_charts' set": {prepare: prepareSuricataChartsAndHiddenCharts},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			suricata := test.prepare()
			require.True(t, suricata.Init())

			if test.wantFail {
				assert.False(t, suricata.Check())
			} else {
				assert.True(t, suricata.Check())
			}
		})
	}
}

func TestSuricata_Charts(t *testing.T) {
	// We want to ensure that initialized module does not return 'nil'.
	// If it is not 'nil' we are ok.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		prepare func(t *testing.T) *Suricata
		wantNil bool
	}{
		"not initialized collector": {
			wantNil: true,
			prepare: func(t *testing.T) *Suricata {
				return New()
			},
		},
		"initialized collector": {
			prepare: func(t *testing.T) *Suricata {
				suricata := New()
				require.True(t, suricata.Init())
				return suricata
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			suricata := test.prepare(t)

			if test.wantNil {
				assert.Nil(t, suricata.Charts())
			} else {
				assert.NotNil(t, suricata.Charts())
			}
		})
	}
}

func TestSuricata_Cleanup(t *testing.T) {
	// Since this module has nothing to clean up,
	// we want just to ensure that Cleanup() not panics.

	assert.NotPanics(t, New().Cleanup)
}

func TestSuricata_Collect(t *testing.T) {
	// 'Collect() map[string]int64' returns collected data, so to test it we need:
	// - provide the module with a specific config.
	// - initialize the module (call Init()).
	// - call Collect() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		prepare       func() *Suricata
		wantCollected map[string]int64
	}{
		"default config": {
			prepare: prepareSuricataDefault,
			wantCollected: map[string]int64{
				"random_0_random0": 1,
				"random_0_random1": -1,
				"random_0_random2": 1,
				"random_0_random3": -1,
			},
		},
		"only 'charts' set": {
			prepare: prepareSuricataOnlyCharts,
			wantCollected: map[string]int64{
				"random_0_random0": 1,
				"random_0_random1": -1,
				"random_0_random2": 1,
				"random_0_random3": -1,
				"random_0_random4": 1,
				"random_1_random0": 1,
				"random_1_random1": -1,
				"random_1_random2": 1,
				"random_1_random3": -1,
				"random_1_random4": 1,
			},
		},
		"only 'hidden_charts' set": {
			prepare: prepareSuricataOnlyHiddenCharts,
			wantCollected: map[string]int64{
				"hidden_random_0_random0": 1,
				"hidden_random_0_random1": -1,
				"hidden_random_0_random2": 1,
				"hidden_random_0_random3": -1,
				"hidden_random_0_random4": 1,
				"hidden_random_1_random0": 1,
				"hidden_random_1_random1": -1,
				"hidden_random_1_random2": 1,
				"hidden_random_1_random3": -1,
				"hidden_random_1_random4": 1,
			},
		},
		"'charts' and 'hidden_charts' set": {
			prepare: prepareSuricataChartsAndHiddenCharts,
			wantCollected: map[string]int64{
				"hidden_random_0_random0": 1,
				"hidden_random_0_random1": -1,
				"hidden_random_0_random2": 1,
				"hidden_random_0_random3": -1,
				"hidden_random_0_random4": 1,
				"hidden_random_1_random0": 1,
				"hidden_random_1_random1": -1,
				"hidden_random_1_random2": 1,
				"hidden_random_1_random3": -1,
				"hidden_random_1_random4": 1,
				"random_0_random0":        1,
				"random_0_random1":        -1,
				"random_0_random2":        1,
				"random_0_random3":        -1,
				"random_0_random4":        1,
				"random_1_random0":        1,
				"random_1_random1":        -1,
				"random_1_random2":        1,
				"random_1_random3":        -1,
				"random_1_random4":        1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			suricata := test.prepare()
			require.True(t, suricata.Init())

			collected := suricata.Collect()

			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, suricata, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, e *Suricata, collected map[string]int64) {
	for _, chart := range *e.Charts() {
		if chart.Obsolete {
			continue
		}
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok,
				"collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok,
				"collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareSuricataDefault() *Suricata {
	return prepareSuricata(New().Config)
}

func prepareSuricataOnlyCharts() *Suricata {
	return prepareSuricata(Config{
		Charts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
	})
}

func prepareSuricataOnlyHiddenCharts() *Suricata {
	return prepareSuricata(Config{
		HiddenCharts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
	})
}

func prepareSuricataChartsAndHiddenCharts() *Suricata {
	return prepareSuricata(Config{
		Charts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
		HiddenCharts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
	})
}

func prepareSuricata(cfg Config) *Suricata {
	suricata := New()
	suricata.Config = cfg
	suricata.randInt = func() int64 { return 1 }
	return suricata
}
