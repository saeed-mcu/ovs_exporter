// Copyright 2025 OVS Exporter Authors
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

package ovs_exporter

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-kit/log/level"
)

// GetSystemIDFromDatabase attempts to retrieve the system-id from the OVS database
// using ovs-vsctl. This is the preferred method for newer OVS versions.
func (e *Exporter) GetSystemIDFromDatabase() (string, error) {
	cmd := exec.Command("ovs-vsctl", "get", "Open_vSwitch", ".", "external-ids:system-id")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get system-id from database: %w", err)
	}

	// The output is typically quoted, so we need to trim quotes and whitespace
	systemID := strings.TrimSpace(string(output))
	systemID = strings.Trim(systemID, "\"")

	if systemID == "" {
		return "", fmt.Errorf("system-id is empty in database")
	}

	level.Debug(e.logger).Log(
		"msg", "Retrieved system-id from database",
		"system_id", systemID,
	)

	return systemID, nil
}

// GetSystemIDFromFile reads the system-id from a file.
// This is the fallback method for older OVS versions or when the database doesn't have it.
func (e *Exporter) GetSystemIDFromFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open system-id file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var systemID string
	for scanner.Scan() {
		systemID = strings.TrimSpace(scanner.Text())
		if systemID != "" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading system-id file: %w", err)
	}

	if systemID == "" {
		return "", fmt.Errorf("system-id is empty in file")
	}

	level.Debug(e.logger).Log(
		"msg", "Retrieved system-id from file",
		"system_id", systemID,
		"file", filepath,
	)

	return systemID, nil
}

// GetSystemID attempts to retrieve the system-id, first from the database,
// then falling back to the file if necessary.
func (e *Exporter) GetSystemID() error {
	// First, try to get system-id from the database (newer OVS versions)
	systemID, err := e.GetSystemIDFromDatabase()
	if err == nil && systemID != "" {
		e.Client.System.ID = systemID
		level.Info(e.logger).Log(
			"msg", "System ID retrieved from database",
			"system_id", systemID,
		)
		return nil
	}

	level.Debug(e.logger).Log(
		"msg", "Failed to get system-id from database, trying file",
		"error", err,
	)

	// Fallback to reading from file (older OVS versions or when not in database)
	systemIDPath := e.Client.Database.Vswitch.File.SystemID.Path
	if systemIDPath == "" {
		systemIDPath = "/etc/openvswitch/system-id.conf"
	}

	systemID, err = e.GetSystemIDFromFile(systemIDPath)
	if err != nil {
		return fmt.Errorf("failed to get system-id from both database and file: %w", err)
	}

	e.Client.System.ID = systemID
	level.Info(e.logger).Log(
		"msg", "System ID retrieved from file",
		"system_id", systemID,
		"file", systemIDPath,
	)

	return nil
}