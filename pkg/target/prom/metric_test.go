package prom

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testMetric struct {
	Name  string  `prometheus:"name"`
	Other int     `prometheus:"other"`
	Value float64 `prometheus:",value"`
}

func TestGetMetric(t *testing.T) {
	expectedValue := 1.0
	expectedLabels := map[string]string{
		"name":  "test",
		"other": "1",
	}
	input := testMetric{"test", 1, expectedValue}
	labels, value := GetMetric(input)
	if value != expectedValue {
		t.Fatalf("want %v got %v", expectedValue, value)
	}
	if diff := cmp.Diff(expectedLabels, labels); diff != "" {
		t.Fatal(diff)
	}
}

func TestSetMetric(t *testing.T) {
	if _, err := SetMetric(testMetric{}); err != nil {
		t.Fatal(err)
	}
}
