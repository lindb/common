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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSetting_initLevel(t *testing.T) {
	initLogLevel("")
	initLogLevel("NO")
	initLogLevel("info")
}

func Test_IsDebug(t *testing.T) {
	RunningAtomicLevel.SetLevel(zapcore.InfoLevel)
	assert.False(t, IsDebug())
	RunningAtomicLevel.SetLevel(zapcore.DebugLevel)
	assert.True(t, IsDebug())
}

func Test_InitLogger(t *testing.T) {
	defaultSetting := NewDefaultSetting()
	encoderConfig := zap.NewProductionEncoderConfig()

	cfg1 := Setting{Level: "LLL"}
	log, err := InitLogger("test.log", cfg1, &encoderConfig)
	assert.Error(t, err)
	assert.Nil(t, log)

	log, err = InitLogger("test.log", *defaultSetting, &encoderConfig)
	assert.NoError(t, err)
	assert.NotNil(t, log)

	cfg3 := Setting{Level: "info"}
	log, err = InitLogger("test.log", cfg3, &encoderConfig)
	assert.NoError(t, err)
	assert.NotNil(t, log)

	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	cfg4 := Setting{Level: "debug"}
	log, err = InitLogger("test.log", cfg4, &encoderConfig)
	assert.NoError(t, err)
	assert.NotNil(t, log)
}

func TestRegisterLogger(t *testing.T) {
	RegisterLogger("test", defaultLogger, false)
	log, ok := loggers["test"]
	assert.True(t, ok)
	assert.Equal(t, defaultLogger, log.logger)
	assert.False(t, log.ignoreModuleAndRole)

	assert.NotNil(t, GetLogger("test", "test"))
}

func TestRegisterLogger_Default(t *testing.T) {
	DefaultLogger.Store(defaultLogger)
	assert.NotNil(t, GetLogger("test11", "test"))
}
