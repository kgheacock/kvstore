package ctx

type ContextKey string

const (
	ContextSourceKey ContextKey = "source"
	EXTERNAL         string     = "EXTERNAL"
	INTERNAL         string     = "INTERNAL"

	ContextCausalContextKey ContextKey = "vector_clock"
)