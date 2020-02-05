# zaputils

zaputils是使用zap的日志使用工具类

## 使用方法

```
//初始化日志组件
log.InitLog("./zap.xml")
//打印日志
log.Logger.Error("GetIndex", zap.String("Security", fd.Security), zap.Error(err))
```

## zap.xml配置说明

```
<?xml version="1.0" encoding="UTF-8"?>
<Loggers>
    <Logger name="Logger">
        <!--日志文件路径-->
        <FilePath>./logs/fundqin.log</FilePath>
        <!--日志级别-->
        <Level>DebugLevel</Level>
        <!--每个日志文件保存的最大尺寸 单位：M-->
        <MaxSize>128</MaxSize>
        <!--日志文件最多保存多少个备份-->
        <MaxBackups>30</MaxBackups>
        <!--文件最多保存多少天-->
        <MaxAge>7</MaxAge>
        <!--是否压缩-->
        <Compress>true</Compress>
    </Logger>

</Loggers>

```
## 日志框架说明

zap
zap是uber开源的Go高性能日志库
https://github.com/uber-go/zap

lumberjack
Lumberjack用于将日志写入滚动文件。zap 不支持文件归档，如果要支持文件按大小或者时间归档，需要使用lumberjack，lumberjack也是zap官方推荐的。
https://github.com/natefinch/lumberjack
