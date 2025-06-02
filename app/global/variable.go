package variable

import (
	"douyin-backend/app/utils/yml_config/interf"
	"go.uber.org/zap"
)

var (
	BasePath        string                    // 定义项目的根目录
	ConfigYml       interf.YmlConfigInterface //  程序退出时需要销毁的事件前缀
	ConfigGormv2Yml interf.YmlConfigInterface //  配置文件键值缓存时，键的前缀
	DateFormat      = "2024-01-02 15:04:05"

	ZapLog *zap.Logger // 全局日志指针
)
