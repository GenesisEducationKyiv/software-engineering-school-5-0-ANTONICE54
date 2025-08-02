package metrics

import "time"

type BroadcastRecorder interface {
	RecordBroadcastDuration(frequency string, duration time.Duration)
}
