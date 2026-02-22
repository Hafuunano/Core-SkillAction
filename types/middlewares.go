package types

type MiddlewareEngine struct {
	MiddlewareID   string `yaml:"middleware_id" json:"middleware_id"`     // Only and only one
	MiddlewareName string `yaml:"middleware_name" json:"middleware_name"` // what actually called in the code
}

// DefaultMiddlewareEngine returns a MiddlewareEngine with minimal default values (like Python __init__ defaults).
// Set MiddlewareName (and optionally MiddlewareID) when using with config.Init.
func DefaultMiddlewareEngine() MiddlewareEngine {
	return MiddlewareEngine{}
}
