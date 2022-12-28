# BiliBili-Terminal

![](https://404name.oss-cn-shanghai.aliyuncs.com/project/bilibili.gif)

> 本项目长期维护~目前仅开发分支
> 
> 等待V 1.x.x版本应该才是可以使用的
## Build

```bash
# mod load
$ go mod tidy
# r启动
$ go run main.go
# 打包
$ go build main.go
```

## Features
> **系统进度** `[==>20%------------------------]` 
> **设计进度** `[========>35%------------------]` 

* [x] 命令行播放视频
* [ ] 音频播放
* [ ] 直接解析bv
* [ ] 完成组件构建
* [ ] 日志收集器
* [ ] 弹幕解析 
* [ ] ...
# TODO

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

| Date      | Versiton | Changelog                                  |
| --------- | -------- | ------------------------------------------ |
| 2022.12.28 | v0.0.2   | 📺 **[refactor]** 重构项目框架<br/>☀️ **[perf]** 优化播放逻辑 <br/>☀️ **[docs]** 更新README文档  <br/>🚀 **[feat]** 【实现响应式布局】/【实现进度条】/ 【实现预加载播放】   |
| 2022.12.27  | v0.0.1   | 📺 **[init]** 初始化界面<br/>📺 **[feat]** 完成页面布局 <br/>🚀 **[feat]** 实现视频播放  |




## License

* [MIT](https://github.com/404name/bilibili-terminal/blob/master/LICENSE)

## Copyright

* **Author:** [404name](https://github.com/404name)
* **AboutMe:** [here](https://yuque.com/404name)
