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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_Logger_Enabled(t *testing.T) {
	logger1 := GetLogger("PKG", "Log")
	assert.False(t, logger1.Enabled(DebugLevel))
	assert.True(t, logger1.Enabled(WarnLevel))
	assert.True(t, logger1.Enabled(FatalLevel))
}

func Test_Logger(t *testing.T) {
	logger1 := GetLogger("PKG", "Log")
	RunningAtomicLevel.SetLevel(zapcore.DebugLevel)

	fmt.Println(White.Add("white"))
	logger1.Warn("warn for test", String("count", "1"), Reflect("v1", map[string]string{"a": "1"}))
	logger1.Info("info for test", Uint16("value", 1), Int32("v1", 2),
		Int64("v2", 2), Any("v3", 3), Int("v", 1))
	logger1.Debug("debug for test", Uint32("value", 2))
	logger1.Error("error for test", Error(fmt.Errorf("error")))

	logger3 := GetLogger("PKG", "")
	logger3.Error("error test")
}

func Test_Level_String(t *testing.T) {
	defer func() {
		isTerminal = false
	}()
	isTerminal = true
	assert.Equal(t, "\x1b[35mDEBUG\x1b[0m", LevelString(zapcore.DebugLevel))
	assert.Equal(t, "\x1b[31mDPANIC\x1b[0m", LevelString(zapcore.DPanicLevel))
	assert.Equal(t, "\x1b[32mINFO\x1b[0m", LevelString(zapcore.InfoLevel))
	assert.Equal(t, "\x1b[33mWARN\x1b[0m", LevelString(zapcore.WarnLevel))
	assert.Equal(t, "\x1b[31mERROR\x1b[0m", LevelString(zapcore.ErrorLevel))
	isTerminal = false
	assert.Equal(t, "ERROR", LevelString(zapcore.ErrorLevel))
}

func Test_IsTerminal(t *testing.T) {
	defer func() {
		isWindowsFn = isWindows
	}()
	assert.False(t, IsTerminal(os.Stdout))
	isWindowsFn = func() bool {
		return true
	}
	assert.False(t, IsTerminal(os.Stdout))
	fmt.Println(isWindows())
}

func Test_foramt_msg(t *testing.T) {
	defer func() {
		isTerminal = false
	}()
	isTerminal = true
	log := &logger{}
	assert.NotEmpty(t, log.formatMsg("test"))
	isTerminal = false
	assert.Empty(t, log.formatMsg(""))
	isTerminal = true
	log.role = "role"
	assert.NotEmpty(t, log.formatMsg(""))
}

func Test_SimpleAccessLevelEncoder(t *testing.T) {
	isTerminal = true
	defer func() {
		isTerminal = IsTerminal(os.Stdout)
	}()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = SimpleTimeEncoder
	encoderConfig.EncodeLevel = SimpleAccessLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		os.Stdout,
		RunningAtomicLevel)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	log.Info("hello", Stack())
}
