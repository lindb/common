// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	isatty "github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	isWindowsFn = isWindows
)

// SimpleTimeEncoder serializes a time.Time to a simplified format without timezone
func SimpleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// SimpleLevelEncoder serializes a Level to a upper case string. For example,
// InfoLevel is serialized to "INFO".
func SimpleLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(LevelString(l))
}

// SimpleAccessLevelEncoder serializes access log level
func SimpleAccessLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if isTerminal {
		enc.AppendString(LevelString(l))
	}
}

// LevelString returns a upper-case ASCII representation of the log level.
func LevelString(l zapcore.Level) string {
	if !isTerminal {
		return l.CapitalString()
	}
	switch l {
	case zapcore.DebugLevel:
		return Magenta.Add(l.CapitalString())
	case zapcore.InfoLevel:
		return Green.Add(l.CapitalString())
	case zapcore.WarnLevel:
		return Yellow.Add(l.CapitalString())
	case zapcore.ErrorLevel:
		return Red.Add(l.CapitalString())
	default:
		return Red.Add(l.CapitalString())
	}
}

// IsTerminal checks if the stdOut is a terminal or not
func IsTerminal(f *os.File) bool {
	if isWindowsFn() {
		return false
	}
	fd := f.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

// logger is wrapper for zap logger with module.
type Logger interface {
	// Debug logs a message at DebugLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Debug(msg string, fields ...zap.Field)
	// Info logs a message at InfoLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Info(msg string, fields ...zap.Field)
	// Warn logs a message at WarnLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Warn(msg string, fields ...zap.Field)
	// Error logs a message at ErrorLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Error(msg string, fields ...zap.Field)
}

// logger implements Logger interface.
type logger struct {
	module              string
	role                string
	ignoreModuleAndRole bool

	log *zap.Logger
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(l.formatMsg(msg), fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(l.formatMsg(msg), fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(l.formatMsg(msg), fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(l.formatMsg(msg), fields...)
}

// formatMsg formats msg using module name
func (l *logger) formatMsg(msg string) string {
	if !isTerminal || l.ignoreModuleAndRole {
		return msg
	}
	moduleName := fmt.Sprintf("[%*s]", atomic.LoadUint32(&maxModuleNameLen), l.module)
	if isTerminal {
		moduleName = Cyan.Add(moduleName)
	}
	if l.role == "" {
		return fmt.Sprintf("%s: %s",
			moduleName, msg)
	}
	return fmt.Sprintf("%s [%s]: %s",
		moduleName, l.role, msg)
}

// String constructs a field with the given key and value.
func String(key, val string) zap.Field {
	return zap.Field{Key: key, Type: zapcore.StringType, String: val}
}

// Error is shorthand for the common idiom NamedError("error", err).
func Error(err error) zap.Field {
	return zap.NamedError("error", err)
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Uint16Type, Integer: int64(val)}
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Uint32Type, Integer: int64(val)}
}

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack() zap.Field {
	return zap.Stack("stack")
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) zap.Field {
	return zap.Reflect(key, val)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Int32Type, Integer: int64(val)}
}

// Int is a shortcut for int.
func Int(key string, val int) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Int32Type, Integer: int64(val)}
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Int64Type, Integer: val}
}
