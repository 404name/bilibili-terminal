package global

import (
	"github.com/404name/termui-demo/global/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// UI组件
var (
	CONFIG config.Config
	VIPER  *viper.Viper
	LOG    *zap.SugaredLogger
	PATH   string
)
