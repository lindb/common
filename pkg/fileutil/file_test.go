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

package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testPath = "./file"

func TestMkDirIfNotExist(t *testing.T) {
	defer func() {
		mkdirAllFunc = os.MkdirAll
		_ = RemoveDir(testPath)
	}()

	mkdirAllFunc = func(path string, perm os.FileMode) error {
		return fmt.Errorf("err")
	}
	err := MkDirIfNotExist(testPath)
	assert.Error(t, err)

	err = MkDir(testPath)
	assert.Error(t, err)
	mkdirAllFunc = os.MkdirAll
	err = MkDir(testPath)
	assert.NoError(t, err)
}

func TestRemoveDir(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		removeAllFunc = os.RemoveAll
		_ = RemoveDir(testPath)
	}()
	removeAllFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	err := RemoveDir(testPath)
	assert.Error(t, err)
}

func TestFileUtil(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
	}()

	assert.True(t, Exist(testPath))
}

func TestFileUtil_errors(t *testing.T) {
	// not existent directory
	_, err := ListDir(filepath.Join(os.TempDir(), "tmp", "tmp", "tmp", "tmp"))

	// encode toml failure
	assert.NotNil(t, err)
}

func TestGetExistPath(t *testing.T) {
	temp := t.TempDir()
	_ = MkDirIfNotExist(temp)
	assert.Equal(t, temp, GetExistPath(filepath.Join(temp, "a", "b", "c")))
	assert.Equal(t, "", GetExistPath(filepath.Join("tmp", "test1", "test333")))
}

func TestListDir(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
	}()
	f, err := os.Create(filepath.Join(testPath, "file1"))
	assert.NoError(t, err)
	assert.NotNil(t, f)
	err = f.Close()
	assert.NoError(t, err)
	files, err := ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
}

func TestRemoveFile(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
		removeFunc = os.Remove
	}()
	f, err := os.Create(filepath.Join(testPath, "file1"))
	assert.NoError(t, err)
	assert.NotNil(t, f)
	err = f.Close()
	assert.NoError(t, err)

	err = RemoveFile(filepath.Join(testPath, "file1"))
	assert.NoError(t, err)
	files, err := ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 0)

	f, err = os.Create(filepath.Join(testPath, "file1"))
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.NoError(t, f.Close())
	removeFunc = func(name string) error {
		return fmt.Errorf("err")
	}
	err = RemoveFile(filepath.Join(testPath, "file1"))
	assert.Error(t, err)
	err = RemoveFile(filepath.Join(testPath, "file2"))
	assert.NoError(t, err)
	files, err = ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
}

func TestFileUtil_GetDirectoryList(t *testing.T) {
	rs, err := GetDirectoryList("../")
	assert.NoError(t, err)
	assert.NotEmpty(t, rs)

	rs, err = GetDirectoryList("not_exist")
	assert.Error(t, err)
	assert.Empty(t, rs)
}
