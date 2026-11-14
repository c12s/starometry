package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/c12s/metrics/internal/config"
)

const nodeExporterSample = `# HELP node_load1 1m load average.
# TYPE node_load1 gauge
node_load1 0.42
`

func newTestMetricsService() *MetricsService {
	return NewMetricsService(NewLocalFileService(), config.NewMetricsConfigWithPresetConfiguration(), "node-a", config.NewAppConfigFromEnv())
}

// A bootc node has node-exporter but no cAdvisor: it must still report.
func TestScrapeSourcesSkipsMissingCAdvisor(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(nodeExporterSample))
	}))
	defer srv.Close()
	addr := strings.TrimPrefix(srv.URL, "http://")

	m := newTestMetricsService()
	for name, cadvisor := range map[string]string{"unset": "", "unreachable": "127.0.0.1:1"} {
		metrics, err := m.scrapeSources([]metricsSource{
			{name: "cAdvisor", address: cadvisor},
			{name: "node-exporter", address: addr},
		})
		if err != nil {
			t.Fatalf("%s cAdvisor: %s", name, err.GetErrorMessage())
		}
		found := false
		for _, md := range metrics {
			if md.MetricName == "node_load1" {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s cAdvisor: node_load1 missing from %d metrics", name, len(metrics))
		}
	}
}

func TestScrapeSourcesFailsWhenNothingAnswers(t *testing.T) {
	m := newTestMetricsService()
	if _, err := m.scrapeSources([]metricsSource{{name: "node-exporter", address: "127.0.0.1:1"}}); err == nil {
		t.Fatal("want error when no source answers")
	}
	if _, err := m.scrapeSources([]metricsSource{{name: "cAdvisor"}}); err == nil {
		t.Fatal("want error when no source is configured")
	}
}
