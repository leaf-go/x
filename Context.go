package x

/*
	x包定义框架基础能力,也可以配合使用
	DefaultConfigs 配置
	Context 中间上下文
	Error 错误
	defaultLog 日志配置
	Validate 验证
*/

import (
	"github.com/gin-gonic/gin"
)

type Roboter interface {
	Send(message string) error
}

type Context interface {
	Logger
	Set(key string, value interface{})
	Get(key string, def interface{}) interface{}
	Response(code int, h interface{})
}

type GinContext struct {
	Logger
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

func (g GinContext) Response(code int, h interface{}) {
	if g.ctx.Writer.Written() {
		return
	}

	g.ctx.AbortWithStatusJSON(code, h)
}

func NewContextWithGin(ctx *gin.Context, log Logger) Context {
	return &GinContext{ctx: ctx, Logger: log}
}
