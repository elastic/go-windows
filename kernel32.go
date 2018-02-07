// Copyright 2018 Elasticsearch Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build windows

package windows

import (
	"fmt"
	"syscall"
)

// Syscalls
//sys   _GetTickCount64() (millis uint64, err error) = kernel32.GetTickCount64

// Version identifies a Windows version by major, minor, and build number.
type Version struct {
	Major int
	Minor int
	Build int
}

// GetWindowsVersion returns the Windows version information. Applications not
// manifested for Windows 8.1 or Windows 10 will return the Windows 8 OS version
// value (6.2).
//
// For a table of version numbers see:
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724833(v=vs.85).aspx
func GetWindowsVersion() Version {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724439(v=vs.85).aspx
	ver, err := syscall.GetVersion()
	if err != nil {
		// GetVersion should never return an error.
		panic(fmt.Errorf("GetVersion failed: %v", err))
	}

	return Version{
		Major: int(ver & 0xFF),
		Minor: int(ver >> 8 & 0xFF),
		Build: int(ver >> 16),
	}
}

// IsWindowsVistaOrGreater returns true if the Windows version is Vista or
// greater.
func (v Version) IsWindowsVistaOrGreater() bool {
	// Vista is 6.0.
	return v.Major >= 6 && v.Minor >= 0
}

// GetTickCount64 retrieves the number of milliseconds that have elapsed since
// the system was started.
// This function is available on Windows Vista and newer.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724411(v=vs.85).aspx
func GetTickCount64() (uint64, error) {
	return _GetTickCount64()
}
