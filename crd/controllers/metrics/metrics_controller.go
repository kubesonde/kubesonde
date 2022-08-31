package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Counter        = 0
	MetricsSummary = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kubesonde",
			Name:      "link",
			Help:      "Link between two pods",
		},
		[]string{"from", "to", "label", "exists"},
	)

	DurationSummary = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kubesonde",
			Name:      "duration",
			Help:      "Probe duration",
		},
		[]string{"start", "end", "counter"},
	)

	TargetedMetricsSummary = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kubesonde",
			Name:      "targetedLink",
			Help:      "Link between two pods when running in targeted mode",
		},
		[]string{"from", "to", "label", "exists", "shouldExist"},
	)
)
