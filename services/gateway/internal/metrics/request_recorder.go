package metrics

import "time"

type MetricRecorder interface {
	RecordRequest(path, method string, duration time.Duration)
}
