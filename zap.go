package log

import (
	"encoding/xml"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
)

var Logger *zap.Logger

var CFG_LIST LoggerCfgList

var DEFAULT_LOGGER_CFG = LoggerCfg{
	XMLName:    xml.Name{},
	CfgName:    "",
	FilePath:   "./logs/zap.log",
	Level:      "InfoLevel",
	MaxSize:    512,
	MaxBackups: 30,
	MaxAge:     10,
	Compress:   true,
}

type LoggerCfgList struct {
	XMLName xml.Name    `xml:"Loggers"`
	Zapcfgs []LoggerCfg `xml:"Logger"`
}

type LoggerCfg struct {
	XMLName    xml.Name `xml:"Logger"`
	CfgName    string   `xml:"name,attr"`
	FilePath   string   `xml:"FilePath"`
	Level      string   `xml:"Level"`
	MaxSize    int      `xml:"MaxSize"`
	MaxBackups int      `xml:"MaxBackups"`
	MaxAge     int      `xml:"MaxAge"`
	Compress   bool     `xml:"Compress"`
}

func InitLog(logFile string) {
	loadCfg(logFile)

	Logger = newLogger("Logger")
}

func loadCfg(logFile string) {
	_, err := os.Stat(logFile)
	if err != nil {
		fmt.Println("zap.xml not found,use default config!")
	} else {
		cfgContent, err := ioutil.ReadFile(logFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = xml.Unmarshal(cfgContent, &CFG_LIST)
		if err != nil {
			fmt.Println("zap.xml load err!")
			return
		}
	}
}
func getConfig(loggerName string) (logCfg *LoggerCfg, err error) {

	if CFG_LIST.Zapcfgs == nil {
		return &DEFAULT_LOGGER_CFG, nil
	}
	for _, tmpCfg := range CFG_LIST.Zapcfgs {
		if tmpCfg.CfgName == loggerName {
			logCfg = &tmpCfg
			break
		}

	}
	if logCfg == nil {
		err = errors.New("zap.xml, Logger name[" + loggerName + "] config err!")
		return nil, err
	}
	return

}

func newLogger(loggerName string) *zap.Logger {
	loggerCfg, err := getConfig(loggerName)
	if err != nil {
		println(err.Error())
		return nil
	}

	core := newCore(*loggerCfg)
	//return zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
	return zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields())

}

func newCore(loggerCfg LoggerCfg) zapcore.Core {
	hook := lumberjack.Logger{
		Filename:   loggerCfg.FilePath,   // 日志文件路径
		MaxSize:    loggerCfg.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: loggerCfg.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     loggerCfg.MaxAge,     // 文件最多保存多少天
		Compress:   loggerCfg.Compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(getLogLevel(loggerCfg.Level))
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

}

func getLogLevel(level string) zapcore.Level {
	if level == "DebugLevel" {
		return zapcore.DebugLevel
	} else if level == "InfoLevel" {
		return zapcore.InfoLevel
	} else if level == "ErrorLevel" {
		return zapcore.ErrorLevel
	} else if level == "WarnLevel" {
		return zapcore.WarnLevel
	} else if level == "PanicLevel" {
		return zapcore.PanicLevel
	} else if level == "FatalLevel" {
		return zapcore.FatalLevel
	} else if level == "DPanicLevel" {
		return zapcore.DPanicLevel
	}
	return zapcore.InfoLevel

}
