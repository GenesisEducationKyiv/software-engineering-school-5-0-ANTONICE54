package metrics

type (
	MailRecorder interface {
		RecordEmailSuccess(subject string)
		RecordEmailFail(subject string)
	}
)
