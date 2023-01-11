package x

import (
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"strings"
)

const (
	APP_MODE     = "APP_MODE"
	APP_SERVICES = "APP_SERVICES"
)

var (
	mode           string = "debug"
	ErrNotice             = errors.New("not have able to services, can use -mode debug -srv api/ export APP_MODE APP_SERVICES setting")
	ErrNotFoundSrv        = errors.New("not found services")
)

type IServiceHandler func() IService

type IService interface {
	IInitialize
	IBootstrap
}

type IServices map[string]IService
type IServiceHandlers map[string]IServiceHandler

func (s IServiceHandlers) Register(name string, handler IServiceHandler) {
	s[name] = handler
}

func (s IServiceHandlers) Take(name string) (service IService) {
	if handler, ok := s[name]; ok {
		return handler()
	}

	PrintRed("%s handler not registered", name)
	return
}

type Bootstrap struct {
	record IServiceHandlers
	will   []string
}

func NewBootstrap() *Bootstrap {
	return &Bootstrap{
		record: make(IServiceHandlers),
		will:   make([]string, 0),
	}
}

func (b *Bootstrap) Register(name string, handler IServiceHandler) {
	b.record.Register(name, handler)
}

func (b *Bootstrap) Boot() error {
	if err := b.params(); err != nil {
		return err
	}

	services, err := b.services()
	if err != nil {
		return err
	}

	b.boot(services)
	return nil
}

func (b *Bootstrap) boot(services IServices) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill)

	for name, service := range services {
		_ = service.Initialize()
		PrintGreen("service %s starting", name)
		go service.Boot()
	}

	<-shutdown
	PrintGreen("services did shutdown")
}

func (b *Bootstrap) services() (map[string]IService, error) {
	var services = make(map[string]IService)
	for _, name := range b.will {
		if service := b.record.Take(name); service != nil {
			services[name] = service
		}
	}

	if len(services) == 0 {
		return nil, ErrNotFoundSrv
	}

	return services, nil
}

func (b *Bootstrap) params() error {
	// 判断flag
	mode, services := b.flag()
	if len(services) > 0 {
		gin.SetMode(mode)
		b.formatService(services)
		return nil
	}

	mode, services = b.env()
	if len(services) > 0 {
		gin.SetMode(mode)
		b.formatService(services)
		return nil
	}

	return ErrNotice
}

func (b *Bootstrap) formatService(services string) {
	if -1 == strings.LastIndex(services, ",") {
		b.will = []string{services}
		return
	}

	will := strings.Split(services, ",")
	for _, v := range will {
		b.will = append(b.will, strings.Trim(v, " "))
	}
}

func (b *Bootstrap) flag() (string, string) {
	mode := flag.String("mode", "debug", "请输入环境")
	services := flag.String("srv", "", "请输入服务名,多服务请以逗号分割")
	flag.Parse()

	return *mode, *services
}

func (b *Bootstrap) env() (string, string) {
	envs := os.Environ()
	var ms, i = map[string]string{
		APP_MODE:     "",
		APP_SERVICES: "",
	}, 0

	for _, v := range envs {
		if -1 == strings.LastIndex(v, "=") {
			continue
		}

		item := strings.Split(v, "=")
		if _, ok := ms[item[0]]; !ok {
			continue
		}

		ms[item[0]] = item[1]
		i++

		if i > 1 {
			break
		}
	}

	return ms[APP_MODE], ms[APP_SERVICES]
}
