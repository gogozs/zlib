package http

type (
	RequestOptions struct {
		headers map[string]string
	}
	RequestOptionFunc func(o *RequestOptions)
)

// WithHeaders customized request headers
func WithHeaders(m map[string]string) RequestOptionFunc {
	return func(o *RequestOptions) {
		o.headers = m
	}
}
