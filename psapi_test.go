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
	"slices"
	"syscall"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestGetPerformanceInfo(t *testing.T) {
	pi, err := GetPerformanceInfo()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the structure was populated correctly
	assert.Equal(t, uint32(unsafe.Sizeof(pi)), pi.CB, "CB field should contain the size of the structure")

	// Verify that physical memory values are reasonable
	assert.True(t, pi.PhysicalTotal > 0, "PhysicalTotal should be greater than 0")
	assert.True(t, pi.PhysicalAvailable <= pi.PhysicalTotal, "PhysicalAvailable should not exceed PhysicalTotal")

	// Verify that commit values are reasonable
	assert.True(t, pi.CommitTotal > 0, "CommitTotal should be greater than 0")
	assert.True(t, pi.CommitLimit > 0, "CommitLimit should be greater than 0")
	assert.True(t, pi.CommitTotal <= pi.CommitLimit, "CommitTotal should not exceed CommitLimit")

	// Verify that page size is reasonable (typically at least 4KB on x86/x64)
	assert.True(t, pi.PageSize >= 4096, "PageSize should be at least 4KB")

	// Verify that process and thread counts are reasonable
	assert.True(t, pi.ProcessCount > 0, "ProcessCount should be greater than 0")
	assert.True(t, pi.ThreadCount > 0, "ThreadCount should be greater than 0")
	assert.True(t, pi.ThreadCount >= pi.ProcessCount, "ThreadCount should be at least equal to ProcessCount")

	// Verify that handle count is reasonable
	assert.True(t, pi.HandleCount > 0, "HandleCount should be greater than 0")

	// Verify that kernel memory values are reasonable
	assert.True(t, pi.KernelTotal > 0, "KernelTotal should be greater than 0")
	assert.True(t, pi.KernelPaged <= pi.KernelTotal, "KernelPaged should not exceed KernelTotal")
	assert.True(t, pi.KernelNonpaged <= pi.KernelTotal, "KernelNonpaged should not exceed KernelTotal")

	// Verify that system cache is reasonable
	assert.True(t, pi.SystemCache > 0, "SystemCache should be greater than 0")
}

func TestGetProcessMemoryInfo(t *testing.T) {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		t.Fatal(err)
	}

	info, err := GetProcessMemoryInfo(h)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the structure was populated correctly
	assert.Equal(t, uint32(unsafe.Sizeof(info)), info.cb, "cb field should contain the size of the structure")

	// Verify that memory values are reasonable
	assert.True(t, info.WorkingSetSize > 0, "WorkingSetSize should be greater than 0")
	assert.True(t, info.PeakWorkingSetSize >= info.WorkingSetSize, "PeakWorkingSetSize should be at least equal to WorkingSetSize")
	assert.True(t, info.PrivateUsage > 0, "PrivateUsage should be greater than 0")

	// Verify that page fault count is reasonable
	assert.True(t, info.PageFaultCount > 0, "PageFaultCount should be greater than 0")
}

func TestGetProcessImageFileName(t *testing.T) {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		t.Fatal(err)
	}

	filename, err := GetProcessImageFileName(h)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that we got a valid device path
	assert.True(t, len(filename) > 0, "Filename should not be empty")
	assert.Contains(t, filename, "\\Device\\", "Filename should be a device path starting with \\Device\\")

	// The path should contain the executable name
	assert.Contains(t, filename, ".exe", "Filename should contain .exe extension")
}

func TestEnumProcesses(t *testing.T) {
	pids, err := EnumProcesses()
	if err != nil {
		t.Fatal(err)
	}

	// Verify that we got a reasonable number of processes
	assert.True(t, len(pids) > 0, "Should return at least one process")
	assert.True(t, len(pids) < 10000, "Should not return an unreasonable number of processes")

	// Verify that the current process is in the list
	currentPID := uint32(syscall.Getpid())
	found := slices.Contains(pids, currentPID)
	assert.True(t, found, "Current process PID should be in the returned list")
}
