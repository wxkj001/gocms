package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var MiddlewareModule = fx.Module("middlewareModule", fx.Provide(NewMiddleware))

type MiddlewareResult struct {
	fx.Out
	*Middlewares
}
type MiddlewareParams struct {
	fx.In
	*Middlewares
}
type Middlewares struct {
	Jwt *JwtMiddleware
}

func NewMiddleware(config *viper.Viper, e *casbin.Enforcer) MiddlewareResult {
	//
	return MiddlewareResult{Middlewares: &Middlewares{
		Jwt: NewJwtMiddleware(config, e),
	}}
}
