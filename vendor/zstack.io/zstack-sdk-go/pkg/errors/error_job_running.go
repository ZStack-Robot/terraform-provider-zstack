package errors

type JobRunningError struct {
	msg string
}

func (f *JobRunningError) Error() string {
	return f.msg
}

func IsJobRunningError(err error) bool {
	_, ok := err.(*JobRunningError)
	return ok
}

func NewJobRunningError(msg string) error {
	return &JobRunningError{
		msg: msg,
	}
}
