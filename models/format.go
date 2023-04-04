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

package models

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// TableFormatter represents table formatter for displaying result in terminal.
type TableFormatter interface {
	// ToTable returns string value/row size for displaying result in terminal.
	ToTable() (rows int, tableStr string)
}

// NewTableFormatter creates a writer for table format.
func NewTableFormatter() table.Writer {
	writer := table.NewWriter()
	style := table.StyleDefault
	style.Format.Header = text.FormatDefault
	writer.SetStyle(style)
	return writer
}
