// Copyright 2018 henrylee2cn. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aster

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"os"
	"path/filepath"

	"github.com/henrylee2cn/goutil"
)

// Store formats the module codes and writes to the local files.
func (m *Module) Store() (first error) {
	codes, first := m.Format()
	if first != nil {
		return first
	}
	for _, v := range codes {
		for kk, vv := range v {
			first = writeFile(kk, vv)
			if first != nil {
				return first
			}
		}
	}
	return
}

// Store formats the package codes and writes to the local files.
func (p *Package) Store() (first error) {
	codes, first := p.Format()
	if first != nil {
		return
	}
	for k, v := range codes {
		first = writeFile(k, v)
		if first != nil {
			return first
		}
	}
	return
}

// Store formats the file codes and writes to the local file.
func (f *File) Store() (err error) {
	code, err := f.Format()
	if err != nil {
		return
	}
	return writeFile(f.Filename, code)
}

// Format format the package and returns the string.
// @codes <packageName,<fileName,code>>
func (m *Module) Format() (codes map[string]map[string]string, first error) {
	codes = make(map[string]map[string]string, len(m.Packages))
	for k, v := range m.Packages {
		subcodes, err := v.Format()
		if err != nil {
			first = err
			return
		}
		codes[k] = subcodes
	}
	return
}

// Format format the package and returns the string.
// @codes <fileName,code>
func (p *Package) Format() (codes map[string]string, first error) {
	codes = make(map[string]string, len(p.Files))
	var code string
	for k, v := range p.Files {
		code, first = v.Format()
		if first != nil {
			return
		}
		codes[k] = code
	}
	return
}

// Format formats the file and returns the string.
func (f *File) Format() (string, error) {
	return f.FormatNode(f.File)
}

// String returns the formated file text.
func (f *File) String() string {
	s, err := f.Format()
	if err != nil {
		return fmt.Sprintf("// Formatting error: %s", err.Error())
	}
	return s
}

// FormatNode formats the node and returns the string.
func (m *Module) FormatNode(node ast.Node) (string, error) {
	var dst bytes.Buffer
	err := format.Node(&dst, m.FileSet, node)
	if err != nil {
		return "", err
	}
	return goutil.BytesToString(dst.Bytes()), nil
}

// FormatNode formats the node and returns the string.
func (p *Package) FormatNode(node ast.Node) (string, error) {
	var dst bytes.Buffer
	err := format.Node(&dst, p.FileSet, node)
	if err != nil {
		return "", err
	}
	return goutil.BytesToString(dst.Bytes()), nil
}

// FormatNode formats the node and returns the string.
func (f *File) FormatNode(node ast.Node) (string, error) {
	var dst bytes.Buffer
	err := format.Node(&dst, f.FileSet, node)
	if err != nil {
		return "", err
	}
	return goutil.BytesToString(dst.Bytes()), nil
}

// TryFormatNode formats the node and returns the string,
// returns the default string if fail.
func (f *File) TryFormatNode(node ast.Node, defaultValue ...string) string {
	code, err := f.FormatNode(node)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return code
}

func writeFile(filename, text string) error {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(goutil.StringToBytes(text))
	return err
}
