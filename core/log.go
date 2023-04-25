package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

func InitLog() {
	//var coreArr []zapcore.Core

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()     //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log", //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    100,              //文件大小限制,单位MB
		MaxBackups: 7,                //最大保留日志文件数量
		MaxAge:     30,               //日志文件保留天数
		Compress:   false,            //是否压缩处理
	})
	fileCore := zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel)
	log = zap.New(fileCore, zap.AddCaller())
}

func Info(result string) {
	log.Info(result)
}
func Warn(result string) {
	log.Warn(result)
}
func Error(result string) {
	log.Info(result)
}
