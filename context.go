package x

import (
	"github.com/gin-gonic/gin"
	"github.com/goantor/logs"
)

type Roboter interface {
	Send(message string) error
}

type Context interface {
	logs.Logger
	Set(key string, value interface{})
	Get(key string, def interface{}) interface{}
}

type GinContext struct {
	logs.Logger
	ctx *gin.Context
}

func (g GinContext) Set(key string, value interface{}) {
	g.ctx.Set(key, value)
}

func (g GinContext) Get(key string, def interface{}) interface{} {
	if value, exists := g.ctx.Get(key); exists {
		return value
	}

	return def
}

func NewContextWithGin(ctx *gin.Context, log logs.Logger) Context {
	return &GinContext{ctx: ctx, Logger: log}
}
