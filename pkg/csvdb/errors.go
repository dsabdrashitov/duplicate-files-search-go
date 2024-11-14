package csvdb

type IOError struct {
	comment string
}

func (e IOError) Error() string {
	return e.comment
}
