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
	"os"
	"path/filepath"
	"testing"

	"github.com/greenpau/ovsdb"
)

func TestGetSystemIDFromFile(t *testing.T) {
	// Create a temporary file with a system ID
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "system-id.conf")
	expectedID := "test-system-id-12345"

	if err := os.WriteFile(tmpFile, []byte(expectedID+"\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	logger, err := NewLogger("debug")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	exporter := &Exporter{
		Client: ovsdb.NewOvsClient(),
		logger: logger,
	}

	systemID, err := exporter.GetSystemIDFromFile(tmpFile)
	if err != nil {
		t.Errorf("GetSystemIDFromFile() returned error: %v", err)
	}

	if systemID != expectedID {
		t.Errorf("GetSystemIDFromFile() = %v, want %v", systemID, expectedID)
	}
}

func TestGetSystemIDFromFileEmpty(t *testing.T) {
	// Create an empty file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "system-id.conf")

	if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	logger, err := NewLogger("debug")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	exporter := &Exporter{
		Client: ovsdb.NewOvsClient(),
		logger: logger,
	}

	_, err = exporter.GetSystemIDFromFile(tmpFile)
	if err == nil {
		t.Error("GetSystemIDFromFile() should return error for empty file")
	}
}

func TestGetSystemIDFromFileNotExists(t *testing.T) {
	logger, err := NewLogger("debug")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	exporter := &Exporter{
		Client: ovsdb.NewOvsClient(),
		logger: logger,
	}

	_, err = exporter.GetSystemIDFromFile("/non/existent/file")
	if err == nil {
		t.Error("GetSystemIDFromFile() should return error for non-existent file")
	}
}

func TestGetSystemIDFromFileWithWhitespace(t *testing.T) {
	// Create a file with whitespace around the system ID
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "system-id.conf")
	expectedID := "test-system-id-67890"

	if err := os.WriteFile(tmpFile, []byte("  "+expectedID+"  \n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	logger, err := NewLogger("debug")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	exporter := &Exporter{
		Client: ovsdb.NewOvsClient(),
		logger: logger,
	}

	systemID, err := exporter.GetSystemIDFromFile(tmpFile)
	if err != nil {
		t.Errorf("GetSystemIDFromFile() returned error: %v", err)
	}

	if systemID != expectedID {
		t.Errorf("GetSystemIDFromFile() = %v, want %v", systemID, expectedID)
	}
}