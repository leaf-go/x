package x

type IBootstrap interface {
	Boot() error
}

type IInitialize interface {
	Initialize() error
}

type IConfig interface {
	Config() error
}

type IApplications map[string]IApplication

type IApplication interface {
	IInitialize
	IBootstrap

	// Bootstrap 启动项
	//Bootstrap(boots ...IBootstrap) IApplication

	// Handler 获取实例
	//Handler() interface{}

	// Shutdown 关闭服务
	Shutdown()
}
