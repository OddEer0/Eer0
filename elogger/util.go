package elogger

import (
	"encoding/json"
	"log/slog"
)

func mapLevelToSlogLevel(lvl Level) slog.Level {
	switch lvl {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

//func mapSlogLevelToLevel(lvl slog.Level) Level {
//	switch lvl {
//	case slog.LevelDebug:
//		return DebugLevel
//	case slog.LevelInfo:
//		return InfoLevel
//	case slog.LevelWarn:
//		return WarnLevel
//	default:
//		return ErrorLevel
//	}
//}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	}
	return "UNKNOWN"
}

func (v Value) MarshalJSON() ([]byte, error) {
	r, ok := v.Val.(LogValuer)
	if ok {
		return json.Marshal(r.LogValue())
	}

	return json.Marshal(v.Val)
}
