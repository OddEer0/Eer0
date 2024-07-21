package elogger

import "log/slog"

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

func mapSlogLevelToLevel(lvl slog.Level) Level {
	switch lvl {
	case slog.LevelDebug:
		return DebugLevel
	case slog.LevelInfo:
		return InfoLevel
	case slog.LevelWarn:
		return WarnLevel
	default:
		return ErrorLevel
	}
}
