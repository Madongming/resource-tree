package global

func init() {
	//初始化config
	initConfig()

	//初始化日志系统
	initLog()

	//初始化DB
	initDb()
}
