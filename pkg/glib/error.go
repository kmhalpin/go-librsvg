package glib

type GError struct {
	Domain  uint32
	Code    int
	Message string
}

var _ error = GError{}

func (e GError) Error() string {
	return e.Message
}

func NewGError(domain uint32, code int, message string) GError {
	return GError{
		Domain:  domain,
		Code:    code,
		Message: message,
	}
}
