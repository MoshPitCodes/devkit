package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// Init initializes the global zap logger with console-friendly defaults.
// debug sets the log level to Debug, otherwise Info.
func Init(debug bool) {
	// Check NO_COLOR environment variable
	_, noColor := os.LookupEnv("NO_COLOR")

	// Configure encoder
	encoderCfg := zap.NewProductionEncoderConfig()
	// Use short file paths
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	// Use simpler time format
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	// Encode Level with color
	if noColor {
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder // No color
	} else {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Color based on level
	}

	// Create the base console encoder
	baseEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// Wrap it with our icon encoder
	encoder := &iconEncoder{Encoder: baseEncoder}

	// Set log level
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	// Configure core
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level)

	// Build logger
	// AddCallerSkip(1) adjusts the caller info to show the call site into logging.L()
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	log = logger.Sugar()

	log.Debug("Logger initialized")
}

// L returns the global SugaredLogger instance.
// Panics if Init() has not been called.
func L() *zap.SugaredLogger {
	if log == nil {
		panic("Logger not initialized. Call logging.Init() first.")
	}
	return log
}

// iconEncoder wraps a zapcore.Encoder to prepend an icon based on level.
type iconEncoder struct {
	zapcore.Encoder // Embed the interface
}

// EncodeEntry implements the zapcore.Encoder interface.
func (e *iconEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	icon := ""
	switch entry.Level {
	case zapcore.DebugLevel:
		icon = "🐛 "
	case zapcore.InfoLevel:
		icon = "ℹ️ "
	case zapcore.WarnLevel:
		icon = "⚠️ "
	case zapcore.ErrorLevel:
		icon = "❌ "
	case zapcore.FatalLevel, zapcore.PanicLevel, zapcore.DPanicLevel:
		icon = "💀 "
	}

	// Prepend icon to the message
	entry.Message = icon + entry.Message

	// Call the embedded encoder's EncodeEntry
	return e.Encoder.EncodeEntry(entry, fields)
}

// Clone duplicates the encoder.
func (e *iconEncoder) Clone() zapcore.Encoder {
	// Clone the embedded encoder and wrap it again
	return &iconEncoder{Encoder: e.Encoder.Clone()}
}
