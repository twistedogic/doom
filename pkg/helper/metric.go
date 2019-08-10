package helper

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/prometheus/client_golang/prometheus"
)

const metricTag = "prometheus"

func GetMetric(i interface{}) (map[string]string, float64) {
	s := structs.New(i)
	labels := make(map[string]string)
	value := 0.0
	for _, field := range s.Fields() {
		tag := field.Tag(metricTag)
		if len(tag) == 0 {
			continue
		}
		tokens := strings.Split(tag, ",")
		name, values := tokens[0], tokens[1:]
		if len(values) != 0 {
			value = field.Value().(float64)
		} else {
			labels[name] = field.Value().(string)
		}
	}
	return labels, value
}

func Update(metric *prometheus.GaugeVec, i interface{}) {
	labels, value := GetMetric(i)
	metric.With(prometheus.Labels(labels)).Set(value)
}
