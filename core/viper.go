package core

import (
	"flag"
	"fmt"

	"os"

	"github.com/404name/termui-demo/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
// Author [SliverHorn](https://github.com/SliverHorn)
func Viper(path ...string) *viper.Viper {
	var config string = "config.yaml"

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(global.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv
				config = global.ConfigDefaultFile
				fmt.Printf("您正在使用%s默认变量,config的路径为%s\n", global.ConfigEnv, config)
			} else {
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", global.ConfigEnv, config)
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {

		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	// global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}
