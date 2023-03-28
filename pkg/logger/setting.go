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
	"os"
	"path/filepath"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	// IsCli represents if command-line.
	IsCli      = false
	isTerminal = IsTerminal(os.Stdout)
	// max length of all modules
	maxModuleNameLen uint32
	// RunningAtomicLevel supports changing level on the fly
	RunningAtomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	loggers            = make(map[string]*zap.Logger)
	// uninitialized logger for default usage
	defaultLogger = newDefaultLogger()
)

func init() {
	// get log level from evn
	level := os.Getenv("LOG_LEVEL")
	initLogLevel(level)
}

func RegisterLogger(module string, logger *zap.Logger) {
	loggers[module] = logger
}

func initLogLevel(level string) {
	if level != "" {
		var zapLevel zapcore.Level
		if err := zapLevel.Set(level); err == nil {
			RunningAtomicLevel.SetLevel(zapLevel)
		}
	}
}

// newDefaultLogger creates a default logger for uninitialized usage
func newDefaultLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = SimpleTimeEncoder
	encoderConfig.EncodeLevel = SimpleLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		os.Stdout,
		RunningAtomicLevel)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

func IsDebug() bool {
	return RunningAtomicLevel.Level() == zapcore.DebugLevel
}

// GetLogger return logger with module name
func GetLogger(module, role string) Logger {
	length := len(module)
	for {
		currentMaxModuleLen := atomic.LoadUint32(&maxModuleNameLen)
		if uint32(length) <= currentMaxModuleLen {
			break
		}
		if atomic.CompareAndSwapUint32(&maxModuleNameLen, currentMaxModuleLen, uint32(length)) {
			break
		}
	}
	log, ok := loggers[module]
	if !ok {
		log = defaultLogger
	}
	return &logger{
		module: module,
		role:   role,
		log:    log,
	}
}

// InitLogger initializes a zap logger from user config
func InitLogger(fileName string, setting Setting, cfg *zapcore.EncoderConfig, options ...zap.Option) (*zap.Logger, error) {
	logger, err := initLogger(fileName, setting, cfg, options...)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// initLogger initializes a zap logger for different module
func initLogger(logFilename string, setting Setting, cfg *zapcore.EncoderConfig, options ...zap.Option) (*zap.Logger, error) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(setting.Dir, logFilename),
		MaxSize:    int(setting.MaxSize / 1024 / 1024), // because in lumberjack will * megabyte
		MaxBackups: int(setting.MaxBackups),
		MaxAge:     int(setting.MaxAge),
	})
	// check if it is terminal
	if !IsCli && isTerminal {
		w = os.Stdout
	}
	// parse logging level
	if err := RunningAtomicLevel.UnmarshalText([]byte(setting.Level)); err != nil {
		return nil, err
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(*cfg),
		w,
		RunningAtomicLevel)
	return zap.New(core, options...), nil
}
