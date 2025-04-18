package model

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
	xlog "xorm.io/xorm/log"
)

var (
	SessionIDKey      = "__xorm_session_id"
	SessionKey        = "__xorm_session_key"
	SessionShowSQLKey = "__xorm_show_sql"
)

type LoggerAdapter struct {
	logger  *zap.Logger
	showSQL bool
	level   xlog.LogLevel
}

// AfterSQL implements log.ContextLogger.
func (l *LoggerAdapter) AfterSQL(context xlog.LogContext) {
	var sessionPart string
	v := context.Ctx.Value(SessionIDKey)
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
	}
	if context.ExecuteTime > 0 {
		l.Infof("[SQL]%s %s %v - %v", sessionPart, context.SQL, context.Args, context.ExecuteTime)
	} else {
		l.Infof("[SQL]%s %s %v", sessionPart, context.SQL, context.Args)
	}
}

// BeforeSQL implements log.ContextLogger.
func (l *LoggerAdapter) BeforeSQL(context xlog.LogContext) {
}

func NewLoggerAdapter(logger *zap.Logger) xlog.ContextLogger {
	return &LoggerAdapter{
		logger: logger,
	}
}

func (l *LoggerAdapter) Debug(v ...interface{}) {
	if l.level <= xlog.LOG_DEBUG {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Debug("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Debug("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Debugf(format string, v ...interface{}) {
	if l.level <= xlog.LOG_DEBUG {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Debug("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Debug("xorm", zap.Any("sql", v))
		}
	}
}

func (l *LoggerAdapter) Error(v ...interface{}) {
	if l.level <= xlog.LOG_ERR {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Error("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Error("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Errorf(format string, v ...interface{}) {
	if l.level <= xlog.LOG_ERR {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Error("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Error("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Info(v ...interface{}) {
	if l.level <= xlog.LOG_INFO {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Info("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Info("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Infof(format string, v ...interface{}) {
	if l.level <= xlog.LOG_INFO {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Info("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Info("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Warn(v ...interface{}) {
	if l.level <= xlog.LOG_WARNING {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Warn("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Warn("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Warnf(format string, v ...interface{}) {
	if l.level <= xlog.LOG_WARNING {
		pc, file, lineNo, ok := runtime.Caller(10)
		if ok {
			l.logger.Warn("xorm", zap.Any("sql", v), zap.Any("file", file), zap.Any("line", lineNo), zap.Any("func", runtime.FuncForPC(pc).Name()))
		} else {
			l.logger.Warn("xorm", zap.Any("sql", v))
		}
	}
}
func (l *LoggerAdapter) Level() xlog.LogLevel {
	return l.level
}
func (l *LoggerAdapter) SetLevel(lv xlog.LogLevel) {
	l.level = lv
}
func (l *LoggerAdapter) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = false
		return
	}
	l.showSQL = show[0]
}
func (l *LoggerAdapter) IsShowSQL() bool {
	return l.showSQL
}
