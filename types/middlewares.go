package types

type MiddlewareEngine struct {
	MiddlewareID   string `yaml:"middleware_id" json:"middleware_id"`     // Only and only one
	MiddlewareName string `yaml:"middleware_name" json:"middleware_name"` // what actually called in the code
}
