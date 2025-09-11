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
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// PmdPerformanceMetrics represents PMD performance statistics
type PmdPerformanceMetrics struct {
	PmdID                 string
	NumaID                string
	Iterations           uint64
	BusyCycles           uint64
	CyclesPerIteration   float64
	PacketsPerIteration  float64
	CyclesPerPacket      float64
	PacketsPerBatch      float64
	MaxVhostQueueLength  uint64
	Upcalls              uint64
	UpcallCycles         uint64
	TxRetries            uint64
	TxContention         uint64
	TxIrqs               uint64
}

// GetPmdPerfMetrics retrieves PMD performance metrics using ovs-appctl
func (e *Exporter) GetPmdPerfMetrics() ([]PmdPerformanceMetrics, error) {
	cmd := exec.Command("ovs-appctl", "dpif-netdev/pmd-perf-show")
	output, err := cmd.Output()
	if err != nil {
		// Check if the command is not available (e.g., non-DPDK deployment)
		if strings.Contains(err.Error(), "exit status") {
			// Command not available, return empty metrics
			return []PmdPerformanceMetrics{}, nil
		}
		return nil, fmt.Errorf("failed to execute pmd-perf-show: %w", err)
	}

	return parsePmdPerfOutput(string(output))
}

// parsePmdPerfOutput parses the output of dpif-netdev/pmd-perf-show
func parsePmdPerfOutput(output string) ([]PmdPerformanceMetrics, error) {
	var metrics []PmdPerformanceMetrics
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	// Regular expressions for parsing different sections
	pmdHeaderRe := regexp.MustCompile(`pmd thread numa_id (\d+) core_id (\d+):`)
	iterationsRe := regexp.MustCompile(`iterations:\s+(\d+)\s+\([\d.]+\s+us/it\)`)
	busyCyclesRe := regexp.MustCompile(`busy cycles:\s+([\d.]+)%.*\(([\d.]+) Mcycles.*\)`)
	cyclesPerItRe := regexp.MustCompile(`cycles/it:\s+([\d.]+)\s+\(([\d.]+) Mcycles\)`)
	pktsPerItRe := regexp.MustCompile(`pkts/it:\s+([\d.]+)`)
	cyclesPerPktRe := regexp.MustCompile(`cycles/pkt:\s+([\d.]+)`)
	pktsPerBatchRe := regexp.MustCompile(`avg pkts/batch:\s+([\d.]+)`)
	maxVhostQRe := regexp.MustCompile(`avg max vhost qlen:\s+(\d+)`)
	upcallsRe := regexp.MustCompile(`upcalls:\s+(\d+)\s+\(([\d.]+) us\s+([\d.]+) Mcycles\)`)
	txRetriesRe := regexp.MustCompile(`vhost tx retries:\s+(\d+)`)
	txContentionRe := regexp.MustCompile(`vhost tx contention:\s+(\d+)`)
	txIrqsRe := regexp.MustCompile(`vhost tx irqs:\s+(\d+)`)
	
	var currentMetric *PmdPerformanceMetrics
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check for PMD thread header
		if matches := pmdHeaderRe.FindStringSubmatch(line); matches != nil {
			if currentMetric != nil {
				metrics = append(metrics, *currentMetric)
			}
			currentMetric = &PmdPerformanceMetrics{
				NumaID: matches[1],
				PmdID:  matches[2],
			}
			continue
		}
		
		if currentMetric == nil {
			continue
		}
		
		// Parse iterations
		if matches := iterationsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.Iterations = val
			}
		}
		
		// Parse busy cycles
		if matches := busyCyclesRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.BusyCycles = uint64(val * 1000000) // Convert Mcycles to cycles
			}
		}
		
		// Parse cycles per iteration
		if matches := cyclesPerItRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.CyclesPerIteration = val
			}
		}
		
		// Parse packets per iteration
		if matches := pktsPerItRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.PacketsPerIteration = val
			}
		}
		
		// Parse cycles per packet
		if matches := cyclesPerPktRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.CyclesPerPacket = val
			}
		}
		
		// Parse packets per batch
		if matches := pktsPerBatchRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.PacketsPerBatch = val
			}
		}
		
		// Parse max vhost queue length
		if matches := maxVhostQRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.MaxVhostQueueLength = val
			}
		}
		
		// Parse upcalls
		if matches := upcallsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.Upcalls = val
			}
			if val, err := strconv.ParseFloat(matches[3], 64); err == nil {
				currentMetric.UpcallCycles = uint64(val * 1000000) // Convert Mcycles to cycles
			}
		}
		
		// Parse vhost tx retries
		if matches := txRetriesRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.TxRetries = val
			}
		}
		
		// Parse vhost tx contention
		if matches := txContentionRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.TxContention = val
			}
		}
		
		// Parse vhost tx IRQs
		if matches := txIrqsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.TxIrqs = val
			}
		}
	}
	
	// Add the last metric if exists
	if currentMetric != nil {
		metrics = append(metrics, *currentMetric)
	}
	
	return metrics, nil
}

// GetPmdStatsMetrics retrieves PMD statistics using ovs-appctl dpif-netdev/pmd-stats-show
func (e *Exporter) GetPmdStatsMetrics() ([]PmdPerformanceMetrics, error) {
	cmd := exec.Command("ovs-appctl", "dpif-netdev/pmd-stats-show")
	output, err := cmd.Output()
	if err != nil {
		// Check if the command is not available
		if strings.Contains(err.Error(), "exit status") {
			return []PmdPerformanceMetrics{}, nil
		}
		return nil, fmt.Errorf("failed to execute pmd-stats-show: %w", err)
	}

	// The pmd-stats-show output is similar to pmd-perf-show
	// We can reuse the same parser or create a specialized one if needed
	return parsePmdPerfOutput(string(output))
}