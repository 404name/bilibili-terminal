package widget

import (
	"github.com/404name/termui-demo/model"
)

// 初始化启动要渲染的组件
func Init() {
	go VideoRender(model.Video)

}
