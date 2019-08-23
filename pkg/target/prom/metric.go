package prom

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/prometheus/client_golang/prometheus"
)

const metricTag = "prometheus"

func SetMetric(i interface{}, reg *prometheus.Registry) (*prometheus.GaugeVec, error) {
	name := structs.Name(i)
	labelMap, _ := GetMetric(i)
	labels := make([]string, 0, len(labelMap))
	for label := range labelMap {
		labels = append(labels, label)
	}
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: strings.ToLower(name),
			Help: name,
		},
		labels,
	)
	if err := reg.Register(metric); err != nil {
		return nil, err
	}
	return metric, nil
}

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
			labels[name] = fmt.Sprintf("%v", field.Value())
		}
	}
	return labels, value
}

func Update(metric *prometheus.GaugeVec, i interface{}) error {
	labels, value := GetMetric(i)
	if len(labels) == 0 && value == 0.0 {
		return fmt.Errorf("invalid input %#v", i)
	}
	metric.With(prometheus.Labels(labels)).Set(value)
	return nil
}
