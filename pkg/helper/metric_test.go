package helper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testMetric struct {
	Name  string  `prometheus:"name"`
	Value float64 `prometheus:",value"`
}

func TestGetMetric(t *testing.T) {
	expectedValue := 1.0
	expectedLabels := map[string]string{
		"name": "test",
	}
	input := testMetric{"test", expectedValue}
	labels, value := GetMetric(input)
	if value != expectedValue {
		t.Fatalf("want %v got %v", expectedValue, value)
	}
	if diff := cmp.Diff(expectedLabels, labels); diff != "" {
		t.Fatal(diff)
	}
}
