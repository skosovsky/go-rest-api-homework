package errs

type Error string

const (
	ErrNoBD       = Error("no BD")
	ErrNoAdded    = Error("task not added")
	ErrNotFound   = Error("task not found")
	ErrFiledWrite = Error("filed to write")
	ErrFiledRead  = Error("filed to read")
	ErrMarshal    = "filed to marshal response: %w"
	ErrUnmarshal  = "filed to unmarshal response: %w"
)
