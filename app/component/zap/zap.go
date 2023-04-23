package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type LogConfig struct {
	Level       string // info error warn debug
	Filename    string
	MaxSize     int  // 单文件最大尺寸 MB
	MaxAge      int  // 最多保存天数day
	MaxBackups  int  // 最大备份文件数量
	OutputFile  bool // 是否输出文件
	OutputStdio bool // 是否输出控制台
	Color       bool
}

// InitLogger 初始化Logger
func InitLogger(cfg *LogConfig) (*zap.SugaredLogger, error) {
	writeSyncer := getLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge, cfg.OutputFile, cfg.OutputStdio)
	encoder := getEncoder(cfg.Color)
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, err
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	lg := zap.New(core, zap.AddCaller())
	//zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return lg.Sugar(), nil
}

func getEncoder(color bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "T"
	if color {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	//return zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, outputFile, outputStdio bool) zapcore.WriteSyncer {
	var syncers []zapcore.WriteSyncer
	if outputStdio {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if outputFile {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackup,
			MaxAge:     maxAge,
		}
		syncers = append(syncers, zapcore.AddSync(lumberJackLogger))
	}

	return zap.CombineWriteSyncers(syncers...)
}
