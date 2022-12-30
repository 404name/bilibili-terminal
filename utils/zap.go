package utils

import (
	"github.com/404name/termui-demo/resource"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func InitLogger() {

	writeSyncer := getLogWriter()

	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	Log = logger.Sugar()
}

func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)

}

func getLogWriter() zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{

		Filename: resource.BaseLogPath,

		MaxSize: 1,

		MaxBackups: 5,

		MaxAge: 30,

		Compress: false,
	}

	return zapcore.AddSync(lumberJackLogger)

}
