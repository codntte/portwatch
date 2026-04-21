package history

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintMetrics_Basic(t *testing.T) {
	store, path := tempMetricStore(t)
	_ = store.Append(MetricEntry{Host: "host1", Name: "latency", Value: 5.5, Unit: "ms"})
	_ = store.Append(MetricEntry{Host: "host2", Name: "cpu", Value: 72.0, Unit: "%"})

	var buf bytes.Buffer
	if err := printMetricsTo(&buf, path, "", ""); err != nil {
		t.Fatalf("printMetricsTo: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "host1") {
		t.Error("expected host1 in output")
	}
	if !strings.Contains(out, "latency") {
		t.Error("expected latency in output")
	}
	if !strings.Contains(out, "ms") {
		t.Error("expected unit ms in output")
	}
}

func TestPrintMetrics_FilterHost(t *testing.T) {
	store, path := tempMetricStore(t)
	_ = store.Append(MetricEntry{Host: "hostA", Name: "cpu", Value: 10})
	_ = store.Append(MetricEntry{Host: "hostB", Name: "cpu", Value: 20})

	var buf bytes.Buffer
	if err := printMetricsTo(&buf, path, "hostA", ""); err != nil {
		t.Fatalf("printMetricsTo: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "hostB") {
		t.Error("hostB should not appear when filtering by hostA")
	}
	if !strings.Contains(out, "hostA") {
		t.Error("expected hostA in filtered output")
	}
}

func TestPrintMetrics_NoResults(t *testing.T) {
	_, path := tempMetricStore(t)

	var buf bytes.Buffer
	if err := printMetricsTo(&buf, path, "", ""); err != nil {
		t.Fatalf("printMetricsTo: %v", err)
	}

	if !strings.Contains(buf.String(), "no metric entries found") {
		t.Error("expected no-results message")
	}
}

func TestPrintMetrics_NoUnit(t *testing.T) {
	store, path := tempMetricStore(t)
	_ = store.Append(MetricEntry{Host: "h", Name: "n", Value: 1})

	var buf bytes.Buffer
	if err := printMetricsTo(&buf, path, "", ""); err != nil {
		t.Fatalf("printMetricsTo: %v", err)
	}

	if !strings.Contains(buf.String(), "-") {
		t.Error("expected dash placeholder for missing unit")
	}
}
