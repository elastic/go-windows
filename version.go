// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
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

//go:build windows
// +build windows

package windows

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// FixedFileInfo contains version information for a file. This information is
// language and code page independent. This is an equivalent representation of
// VS_FIXEDFILEINFO.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms646997(v=vs.85).aspx
type FixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

// ProductVersion returns the ProductVersion value in string format.
func (info FixedFileInfo) ProductVersion() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(info.ProductVersionMS >> 16),
		(info.ProductVersionMS & 0xFFFF),
		(info.ProductVersionLS >> 16),
		(info.ProductVersionLS & 0xFFFF))
}

// FileVersion returns the FileVersion value in string format.
func (info FixedFileInfo) FileVersion() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(info.FileVersionMS >> 16),
		(info.FileVersionMS & 0xFFFF),
		(info.FileVersionLS >> 16),
		(info.FileVersionLS & 0xFFFF))
}

// VersionData is a buffer holding the data returned by GetFileVersionInfo.
type VersionData []byte

// QueryValue uses VerQueryValue to query version information from the a
// version-information resource. It returns responses using the first language
// and code point found in the resource. The accepted keys are listed in
// the VerQueryValue documentation (e.g. ProductVersion, FileVersion, etc.).
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647464(v=vs.85).aspx
func (d VersionData) QueryValue(key string) (string, error) {
	type LangAndCodePage struct {
		Language uint16
		CodePage uint16
	}

	var langCodePage *LangAndCodePage
	langCodeLen := uint32(unsafe.Sizeof(*langCodePage))
	if err := windows.VerQueryValue(unsafe.Pointer(&d[0]), `\VarFileInfo\Translation`, (unsafe.Pointer)(&langCodePage), &langCodeLen); err != nil || langCodeLen == 0 {
		return "", fmt.Errorf("failed to get list of languages: %w", err)
	}

	var dataPtr uintptr
	var size uint32
	subBlock := fmt.Sprintf(`\StringFileInfo\%04x%04x\%v`, langCodePage.Language, langCodePage.CodePage, key)
	if err := windows.VerQueryValue(unsafe.Pointer(&d[0]), subBlock, (unsafe.Pointer)(&dataPtr), &size); err != nil || langCodeLen == 0 {
		return "", fmt.Errorf("failed to query %v: %w", subBlock, err)
	}

	offset := int(dataPtr - (uintptr)(unsafe.Pointer(&d[0])))
	if offset <= 0 || offset > len(d)-1 {
		return "", errors.New("invalid address")
	}

	str, _, err := UTF16BytesToString(d[offset : offset+int(size)*2])
	if err != nil {
		return "", fmt.Errorf("failed to decode UTF16 data: %w", err)
	}

	return str, nil
}

// FixedFileInfo returns the fixed version information from a
// version-information resource. It queries the root block to get the
// VS_FIXEDFILEINFO value.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647464(v=vs.85).aspx
func (d VersionData) FixedFileInfo() (*FixedFileInfo, error) {
	if len(d) == 0 {
		return nil, errors.New("use GetFileVersionInfo to initialize VersionData")
	}

	var fixedInfo *FixedFileInfo
	fixedInfoLen := uint32(unsafe.Sizeof(*fixedInfo))
	if err := windows.VerQueryValue(unsafe.Pointer(&d[0]), `\`, (unsafe.Pointer)(&fixedInfo), &fixedInfoLen); err != nil {
		return nil, fmt.Errorf("VerQueryValue failed for \\: %w", err)
	}

	return fixedInfo, nil
}

// GetFileVersionInfo retrieves version information for the specified file.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647003(v=vs.85).aspx
func GetFileVersionInfo(filename string) (VersionData, error) {
	size, err := windows.GetFileVersionInfoSize(filename, nil)
	if err != nil {
		return nil, fmt.Errorf("GetFileVersionInfoSize failed: %w", err)
	}

	data := make(VersionData, size)
	err = windows.GetFileVersionInfo(filename, 0, uint32(len(data)), unsafe.Pointer(&data[0]))
	if err != nil {
		return nil, fmt.Errorf("GetFileVersionInfo failed: %w", err)
	}

	return data, nil
}
