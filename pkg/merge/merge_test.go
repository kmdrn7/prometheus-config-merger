package merge

import (
	"prometheus-config-merger/pkg/config"
	"sync"
	"testing"
)

func TestMerge_Run(t *testing.T) {
	type field struct {
		testname     string
		promeconfigs []config.PrometheusConfig
	}
	tests := []field{
		{
			testname: "test-1",
			promeconfigs: []config.PrometheusConfig{
				{
					Id:     "1",
					Path:   "../../example/prometheus_current.yaml",
					Weight: 1,
				},
				{
					Id:     "2",
					Path:   "../../example/prometheus_operator.yaml",
					Weight: 2,
				},
			},
		},
		{
			testname: "test-2",
			promeconfigs: []config.PrometheusConfig{
				{
					Id:     "1",
					Path:   "../../example/prometheus_current.yaml",
					Weight: 1,
				},
				{
					Id:     "2",
					Path:   "../../example/prometheus_operator.yaml",
					Weight: 2,
				},
			},
		},
		{
			testname: "test-3",
			promeconfigs: []config.PrometheusConfig{
				{
					Id:     "1",
					Path:   "../../example/prometheus_current.yaml",
					Weight: 1,
				},
				{
					Id:     "2",
					Path:   "../../example/prometheus_operator.yaml",
					Weight: 2,
				},
			},
		},
		{
			testname: "test-4",
			promeconfigs: []config.PrometheusConfig{
				{
					Id:     "1",
					Path:   "../../example/prometheus_current.yaml",
					Weight: 1,
				},
				{
					Id:     "2",
					Path:   "../../example/prometheus_operator.yaml",
					Weight: 2,
				},
			},
		},
	}

	tmpDir := t.TempDir()

	cfg := &config.Config{
		TargetPrometheusConfig: tmpDir + "/prometheus_all.yaml",
		Debug:                  true,
	}

	wg := &sync.WaitGroup{}

	for _, tt := range tests {
		wg.Add(1)
		go func(tt field) {
			cfg.PrometheusConfigs = tt.promeconfigs
			t.Run(tt.testname, func(t *testing.T) {
				Run(cfg)
			})
			wg.Done()
		}(tt)
	}

	wg.Wait()
}
