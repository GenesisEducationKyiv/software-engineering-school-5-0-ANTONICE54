package metrics

type SubscriptionRecorder interface {
	RecordSubscriptionCreated()
	RecordSubscriptionConfirmed()
	RecordSubscriptionDeleted()
}
