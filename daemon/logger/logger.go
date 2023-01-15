package logger

import (
	"bytes"
	"fmt"
	"time"

	conf "WBA/daemon/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 전역 로거
var lg *zap.Logger

// 로거 초기화 컨피그 파라메터
func InitLogger(cfg *conf.Config) (err error) {
	cf := cfg.Log

	now := time.Now()
	lPath := fmt.Sprintf("%s_%s.log", cf.Fpath, now.Format("2006-01-02"))
	// 설정 옵션
	writeSyncer := getLogWriter(lPath, cf.Msize, cf.Mbackup, cf.Mage)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cf.Level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	// lg 생성
	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg)
	return
}

func Debug(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(str.(string))
		b.WriteString(" ")
	}

	lg.Debug("debug", zap.String("-", b.String()))
}

// Info is a convenient alias for Root().Info
func Info(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(str.(string))
		b.WriteString(" ")
	}
	lg.Info("info", zap.String("-", b.String()))
}

// Info is a convenient alias for Root().Event
func Event(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(str.(string))
		b.WriteString(" ")
	}
	lg.Info("Event", zap.String("-", b.String()))
}

// Warn is a convenient alias for Root().Warn
func Warn(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(str.(string))
		b.WriteString(" ")
	}
	lg.Warn("warn", zap.String("-", b.String()))
}

// Error is a convenient alias for Root().Error
func Error(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(str.(string))
		b.WriteString(" ")
	}
	lg.Error("error", zap.String("-", b.String()))
}

// encoder 옵션 설정
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 로그파일 명 지정
		MaxSize:    maxSize,   // 로그 파일 최대 사이즈
		MaxBackups: maxBackup, // 최대 로그파일 수
		MaxAge:     maxAge,    // 로그파일 저장 일수
	}
	return zapcore.AddSync(lumberJackLogger)
}
