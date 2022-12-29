# BiliBili-Terminal

![](https://404name.oss-cn-shanghai.aliyuncs.com/project/bilibili.gif)

> 本项目长期维护~目前仅开发分支
> 
> 等待V 1.x.x版本应该才是可以使用的

## Features
> **系统进度** `[==>20%------------------------]` 
> **设计进度** `[========>35%------------------]` 

* [x] 命令行播放视频
* [x] 音频播放
* [ ] 直接解析bv
* [ ] 完成组件构建
* [ ] 日志收集器
* [ ] 弹幕解析 
* [ ] ...


## Build
> 1. **运行环境需要【[ffmpeg](https://juejin.cn/post/6992181270960685086)】**：命令行输入ffmpeg有输出即可
> 
> 2. **编译环境需要【golang】**：我的是go1.9.4，理论上都行

(win推荐使用[choco](https://juejin.cn/post/7009855843260710948)，任何环境直接`choco install xxx`搞定，完全不需要自己动手下载和配置环境变量)
```bash
# mod load
$ go mod tidy
# r启动
$ go run main.go
# 打包
$ go build main.go
```

##  TODO

**近期**
1. 缓存视频再播放
2. 自定义日志,日志可以输出到面板
3. 对接音频
4. 继续重构架构


**后期**

1. 对接b站api
2. 对接弹幕

**已完成**
1. ffmpeg 获取视频长度，做个进度条 
2. 让视频刷新更快
3. 尝试先缓存视频到本地再播放
4. 包装ffmpeg出来
二期
1. 记录播放上次时间
2. 根据resize设置不同的布局 长宽比小于1 就是手机layout 大于1就是电脑布局

##  PATH

```
├─ffmpeg        存放视频处理相关的逻辑
├─global        存放全局变量
├─model         存放对象
├─resource      存储资源(视频,图片,日志)
│  ├─images
│  └─output
├─test          测试目录
│  └─keyboard
├─utils         工具
└─widget        UI刷新相关

```


## ChangeLog
> **feat** 新功能、**fix** 修补 bug、**docs** 文档、**style** 格式、**refactor** 重构、**test** 增加测试、**chore** 构建过程、辅助工具、**perf** 提高性能

| Date      | Versiton | Describe |Changelog                                  |
| --------- | -------- | -------- |------------------------------------------ |
| 2022.12.29 | v0.0.3   | 支持音频了~ |  ☀️ **[perf]** 优化播放逻辑 <br/>💬 **[fix]** 修复线上问题及添加gitignore  <br/>🚀 **[feat]** 【支持音频播放！】/【同步拉取图片及音频】/ 【初步实现音视频同步播放】   |
| 2022.12.28 | v0.0.2   | 优化架构，优化播放  | 📺 **[refactor]** 重构项目框架<br/>☀️ **[perf]** 优化播放逻辑 <br/>📑 **[docs]** 更新README文档  <br/>🚀 **[feat]** 【实现响应式布局】/【实现进度条】/ 【实现预加载播放】   |
| 2022.12.27  | v0.0.1   | 页面布局、支持视频播放  |  📺 **[init]** 初始化界面<br/>🚀 **[feat]** 完成页面布局 <br/>🚀 **[feat]** 实现视频播放  |


## Dependence
- github.com/gizak/termui/v3 **Go命令行UI组件库**
- github.com/hajimehoshi/go-mp3  **音频播放组件**
- github.com/hajimehoshi/oto/v2 **音频播放组件**
- go.uber.org/zap **高性能日志库zap**

## License

* [MIT](https://github.com/404name/bilibili-terminal/blob/master/LICENSE)

## Copyright

* **Author:** [404name](https://github.com/404name)
* **AboutMe:** [here](https://yuque.com/404name)
