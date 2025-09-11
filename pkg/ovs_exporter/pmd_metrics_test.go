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
	"testing"
)

func TestParsePmdPerfOutput(t *testing.T) {
	// Sample output from dpif-netdev/pmd-perf-show
	sampleOutput := `pmd thread numa_id 0 core_id 2:
  iterations:        12345678 (123.45 us/it)
  busy cycles:       75.2% (2345.67 Mcycles, 1234 us/it)
  cycles/it:         1234567.8 (123.45 Mcycles)
  pkts/it:           234.5
  cycles/pkt:        5678.9
  avg pkts/batch:    32.1
  avg max vhost qlen: 128
  upcalls:           1234 (567.8 us 89.0 Mcycles)
  vhost tx retries:  567
  vhost tx contention: 89
  vhost tx irqs:     123

pmd thread numa_id 1 core_id 4:
  iterations:        87654321 (234.56 us/it)
  busy cycles:       82.3% (3456.78 Mcycles, 2345 us/it)
  cycles/it:         2345678.9 (234.56 Mcycles)
  pkts/it:           345.6
  cycles/pkt:        6789.0
  avg pkts/batch:    64.2
  avg max vhost qlen: 256
  upcalls:           2345 (678.9 us 123.4 Mcycles)
  vhost tx retries:  678
  vhost tx contention: 123
  vhost tx irqs:     234`

	metrics, err := parsePmdPerfOutput(sampleOutput)
	if err != nil {
		t.Fatalf("Failed to parse PMD output: %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("Expected 2 PMD metrics, got %d", len(metrics))
	}

	// Check first PMD
	pmd1 := metrics[0]
	if pmd1.NumaID != "0" {
		t.Errorf("Expected NumaID=0, got %s", pmd1.NumaID)
	}
	if pmd1.PmdID != "2" {
		t.Errorf("Expected PmdID=2, got %s", pmd1.PmdID)
	}
	if pmd1.Iterations != 12345678 {
		t.Errorf("Expected Iterations=12345678, got %d", pmd1.Iterations)
	}
	if pmd1.CyclesPerIteration != 1234567.8 {
		t.Errorf("Expected CyclesPerIteration=1234567.8, got %f", pmd1.CyclesPerIteration)
	}
	if pmd1.PacketsPerIteration != 234.5 {
		t.Errorf("Expected PacketsPerIteration=234.5, got %f", pmd1.PacketsPerIteration)
	}
	if pmd1.CyclesPerPacket != 5678.9 {
		t.Errorf("Expected CyclesPerPacket=5678.9, got %f", pmd1.CyclesPerPacket)
	}
	if pmd1.PacketsPerBatch != 32.1 {
		t.Errorf("Expected PacketsPerBatch=32.1, got %f", pmd1.PacketsPerBatch)
	}
	if pmd1.MaxVhostQueueLength != 128 {
		t.Errorf("Expected MaxVhostQueueLength=128, got %d", pmd1.MaxVhostQueueLength)
	}
	if pmd1.Upcalls != 1234 {
		t.Errorf("Expected Upcalls=1234, got %d", pmd1.Upcalls)
	}
	if pmd1.TxRetries != 567 {
		t.Errorf("Expected TxRetries=567, got %d", pmd1.TxRetries)
	}
	if pmd1.TxContention != 89 {
		t.Errorf("Expected TxContention=89, got %d", pmd1.TxContention)
	}
	if pmd1.TxIrqs != 123 {
		t.Errorf("Expected TxIrqs=123, got %d", pmd1.TxIrqs)
	}

	// Check second PMD
	pmd2 := metrics[1]
	if pmd2.NumaID != "1" {
		t.Errorf("Expected NumaID=1, got %s", pmd2.NumaID)
	}
	if pmd2.PmdID != "4" {
		t.Errorf("Expected PmdID=4, got %s", pmd2.PmdID)
	}
	if pmd2.Iterations != 87654321 {
		t.Errorf("Expected Iterations=87654321, got %d", pmd2.Iterations)
	}
}

func TestParsePmdPerfOutputEmpty(t *testing.T) {
	// Test with empty output (non-DPDK deployment)
	metrics, err := parsePmdPerfOutput("")
	if err != nil {
		t.Fatalf("Failed to parse empty PMD output: %v", err)
	}

	if len(metrics) != 0 {
		t.Fatalf("Expected 0 PMD metrics for empty output, got %d", len(metrics))
	}
}

func TestParsePmdPerfOutputInvalid(t *testing.T) {
	// Test with invalid output
	invalidOutput := `This is not valid PMD output
Some random text here`

	metrics, err := parsePmdPerfOutput(invalidOutput)
	if err != nil {
		t.Fatalf("Failed to parse invalid PMD output: %v", err)
	}

	if len(metrics) != 0 {
		t.Fatalf("Expected 0 PMD metrics for invalid output, got %d", len(metrics))
	}
}