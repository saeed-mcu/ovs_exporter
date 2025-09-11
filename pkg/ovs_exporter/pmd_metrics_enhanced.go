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

// EnhancedPmdMetrics represents comprehensive PMD performance statistics
type EnhancedPmdMetrics struct {
	// Core identification
	PmdID   string
	NumaID  string
	CoreID  string
	
	// CPU Utilization
	CPUUtilization float64
	IdleCycles     uint64
	BusyCycles     uint64
	TotalCycles    uint64
	
	// Iteration metrics
	Iterations           uint64
	SleepIterations     uint64
	BusyIterations      uint64
	CyclesPerIteration  float64
	UsPerIteration      float64
	
	// Packet processing
	PacketsPerIteration  float64
	CyclesPerPacket      float64
	PacketsPerBatch      float64
	TotalPackets        uint64
	
	// RX Batch Statistics
	RxBatches           uint64
	RxPackets           uint64
	AvgRxBatchSize      float64
	MaxRxBatchSize      uint64
	
	// TX Batch Statistics
	TxBatches           uint64
	TxPackets           uint64
	AvgTxBatchSize      float64
	
	// vHost Queue Metrics
	MaxVhostQueueLength  uint64
	AvgVhostQueueLength  float64
	VhostQueueFull       uint64
	
	// Upcall metrics
	Upcalls              uint64
	UpcallCycles         uint64
	AvgUpcallCycles      float64
	
	// vHost specific
	VhostTxRetries       uint64
	VhostTxContention    uint64
	VhostTxIrqs          uint64
	
	// Hit/Miss Statistics
	ExactMatchHit        uint64
	MaskedHit           uint64
	Miss                uint64
	Lost                uint64
	
	// Histogram data (if available)
	CyclesHistogram      map[string]uint64
	PacketsHistogram     map[string]uint64
	BatchSizeHistogram   map[string]uint64
	
	// Suspicious iterations
	SuspiciousIterations uint64
	SuspiciousPercent    float64
}

// GetEnhancedPmdMetrics retrieves comprehensive PMD metrics
func (e *Exporter) GetEnhancedPmdMetrics() ([]EnhancedPmdMetrics, error) {
	// First try detailed metrics
	cmd := exec.Command("ovs-appctl", "dpif-netdev/pmd-perf-show")
	output, err := cmd.Output()
	if err != nil {
		// If not available, return empty
		if strings.Contains(err.Error(), "exit status") {
			return []EnhancedPmdMetrics{}, nil
		}
		return nil, fmt.Errorf("failed to execute pmd-perf-show: %w", err)
	}

	metrics := parseEnhancedPmdOutput(string(output))
	
	// Also get pmd-stats-show for additional metrics
	statsCmd := exec.Command("ovs-appctl", "dpif-netdev/pmd-stats-show")
	statsOutput, err := statsCmd.Output()
	if err == nil {
		enrichWithStats(metrics, string(statsOutput))
	}
	
	return metrics, nil
}

// parseEnhancedPmdOutput parses comprehensive PMD performance output
func parseEnhancedPmdOutput(output string) []EnhancedPmdMetrics {
	var metrics []EnhancedPmdMetrics
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	// Enhanced regex patterns
	pmdHeaderRe := regexp.MustCompile(`pmd thread numa_id (\d+) core_id (\d+):`)
	
	// CPU and iteration patterns
	cpuUtilRe := regexp.MustCompile(`(?:cpu|processor) utilization:\s+([\d.]+)%`)
	idleCyclesRe := regexp.MustCompile(`idle cycles:\s+([\d.]+)%.*\(([\d.]+) Mcycles`)
	busyCyclesRe := regexp.MustCompile(`busy cycles:\s+([\d.]+)%.*\(([\d.]+) Mcycles`)
	iterationsRe := regexp.MustCompile(`iterations:\s+(\d+)\s+\(([\d.]+) us/it\)`)
	sleepIterRe := regexp.MustCompile(`sleep iterations:\s+(\d+)\s+\(([\d.]+)%`)
	
	// Packet processing patterns
	cyclesPerItRe := regexp.MustCompile(`cycles/it:\s+([\d.]+)\s+\(([\d.]+) Mcycles\)`)
	pktsPerItRe := regexp.MustCompile(`pkts/it:\s+([\d.]+)`)
	cyclesPerPktRe := regexp.MustCompile(`cycles/pkt:\s+([\d.]+)`)
	pktsPerBatchRe := regexp.MustCompile(`avg pkts/batch:\s+([\d.]+)`)
	
	// RX/TX batch patterns
	rxBatchRe := regexp.MustCompile(`rx batches:\s+(\d+).*avg:\s+([\d.]+).*max:\s+(\d+)`)
	txBatchRe := regexp.MustCompile(`tx batches:\s+(\d+).*avg:\s+([\d.]+)`)
	rxPacketsRe := regexp.MustCompile(`rx packets:\s+(\d+)`)
	txPacketsRe := regexp.MustCompile(`tx packets:\s+(\d+)`)
	
	// vHost queue patterns
	maxVhostQRe := regexp.MustCompile(`(?:avg )?max vhost qlen:\s+(\d+)`)
	avgVhostQRe := regexp.MustCompile(`avg vhost qlen:\s+([\d.]+)`)
	vhostFullRe := regexp.MustCompile(`vhost queue full:\s+(\d+)`)
	
	// Upcall patterns
	upcallsRe := regexp.MustCompile(`upcalls:\s+(\d+)\s+\(([\d.]+) us\s+([\d.]+) Mcycles\)`)
	avgUpcallRe := regexp.MustCompile(`avg upcall cycles:\s+([\d.]+)`)
	
	// vHost TX patterns
	txRetriesRe := regexp.MustCompile(`vhost tx retries:\s+(\d+)`)
	txContentionRe := regexp.MustCompile(`vhost tx contention:\s+(\d+)`)
	txIrqsRe := regexp.MustCompile(`vhost tx irqs:\s+(\d+)`)
	
	// Hit/Miss patterns
	exactHitRe := regexp.MustCompile(`exact match hit:\s+(\d+)`)
	maskedHitRe := regexp.MustCompile(`masked hit:\s+(\d+)`)
	missRe := regexp.MustCompile(`miss:\s+(\d+)`)
	lostRe := regexp.MustCompile(`lost:\s+(\d+)`)
	
	// Suspicious iterations
	suspiciousRe := regexp.MustCompile(`suspicious iterations:\s+(\d+)\s+\(([\d.]+)%\)`)
	
	// Histogram patterns
	histogramRe := regexp.MustCompile(`(\w+) histogram:`)
	histogramEntryRe := regexp.MustCompile(`\s+(\d+-\d+|\d+\+):\s+(\d+)`)
	
	var currentMetric *EnhancedPmdMetrics
	var currentHistogram *map[string]uint64
	var histogramType string
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check for PMD thread header
		if matches := pmdHeaderRe.FindStringSubmatch(line); matches != nil {
			if currentMetric != nil {
				metrics = append(metrics, *currentMetric)
			}
			currentMetric = &EnhancedPmdMetrics{
				NumaID:             matches[1],
				CoreID:             matches[2],
				PmdID:              matches[2],
				CyclesHistogram:    make(map[string]uint64),
				PacketsHistogram:   make(map[string]uint64),
				BatchSizeHistogram: make(map[string]uint64),
			}
			currentHistogram = nil
			continue
		}
		
		if currentMetric == nil {
			continue
		}
		
		// Parse CPU utilization
		if matches := cpuUtilRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.CPUUtilization = val
			}
		}
		
		// Parse idle cycles
		if matches := idleCyclesRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.IdleCycles = uint64(val * 1000000) // Convert Mcycles to cycles
			}
		}
		
		// Parse busy cycles
		if matches := busyCyclesRe.FindStringSubmatch(line); matches != nil {
			if percent, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.CPUUtilization = percent // Busy percentage is CPU utilization
			}
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.BusyCycles = uint64(val * 1000000)
			}
		}
		
		// Parse iterations
		if matches := iterationsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.Iterations = val
			}
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.UsPerIteration = val
			}
		}
		
		// Parse sleep iterations
		if matches := sleepIterRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.SleepIterations = val
				currentMetric.BusyIterations = currentMetric.Iterations - val
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
		
		// Parse RX batch statistics
		if matches := rxBatchRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.RxBatches = val
			}
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.AvgRxBatchSize = val
			}
			if val, err := strconv.ParseUint(matches[3], 10, 64); err == nil {
				currentMetric.MaxRxBatchSize = val
			}
		}
		
		// Parse RX packets
		if matches := rxPacketsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.RxPackets = val
			}
		}
		
		// Parse TX batch statistics
		if matches := txBatchRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.TxBatches = val
			}
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.AvgTxBatchSize = val
			}
		}
		
		// Parse TX packets
		if matches := txPacketsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.TxPackets = val
			}
		}
		
		// Parse vHost queue metrics
		if matches := maxVhostQRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.MaxVhostQueueLength = val
			}
		}
		
		if matches := avgVhostQRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.AvgVhostQueueLength = val
			}
		}
		
		if matches := vhostFullRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.VhostQueueFull = val
			}
		}
		
		// Parse upcalls
		if matches := upcallsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.Upcalls = val
			}
			if val, err := strconv.ParseFloat(matches[3], 64); err == nil {
				currentMetric.UpcallCycles = uint64(val * 1000000)
			}
		}
		
		if matches := avgUpcallRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				currentMetric.AvgUpcallCycles = val
			}
		}
		
		// Parse vHost TX metrics
		if matches := txRetriesRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.VhostTxRetries = val
			}
		}
		
		if matches := txContentionRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.VhostTxContention = val
			}
		}
		
		if matches := txIrqsRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.VhostTxIrqs = val
			}
		}
		
		// Parse hit/miss statistics
		if matches := exactHitRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.ExactMatchHit = val
			}
		}
		
		if matches := maskedHitRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.MaskedHit = val
			}
		}
		
		if matches := missRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.Miss = val
			}
		}
		
		if matches := lostRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.Lost = val
			}
		}
		
		// Parse suspicious iterations
		if matches := suspiciousRe.FindStringSubmatch(line); matches != nil {
			if val, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
				currentMetric.SuspiciousIterations = val
			}
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				currentMetric.SuspiciousPercent = val
			}
		}
		
		// Parse histograms
		if matches := histogramRe.FindStringSubmatch(line); matches != nil {
			histogramType = strings.ToLower(matches[1])
			switch histogramType {
			case "cycles":
				currentHistogram = &currentMetric.CyclesHistogram
			case "packets":
				currentHistogram = &currentMetric.PacketsHistogram
			case "batch":
				currentHistogram = &currentMetric.BatchSizeHistogram
			default:
				currentHistogram = nil
			}
			continue
		}
		
		// Parse histogram entries
		if currentHistogram != nil {
			if matches := histogramEntryRe.FindStringSubmatch(line); matches != nil {
				if val, err := strconv.ParseUint(matches[2], 10, 64); err == nil {
					(*currentHistogram)[matches[1]] = val
				}
			}
		}
	}
	
	// Add the last metric if exists
	if currentMetric != nil {
		// Calculate total cycles and packets if not set
		if currentMetric.TotalCycles == 0 {
			currentMetric.TotalCycles = currentMetric.BusyCycles + currentMetric.IdleCycles
		}
		if currentMetric.TotalPackets == 0 && currentMetric.Iterations > 0 {
			currentMetric.TotalPackets = uint64(float64(currentMetric.Iterations) * currentMetric.PacketsPerIteration)
		}
		metrics = append(metrics, *currentMetric)
	}
	
	return metrics
}

// enrichWithStats adds additional statistics from pmd-stats-show
func enrichWithStats(metrics []EnhancedPmdMetrics, statsOutput string) {
	// Parse additional stats and enrich the metrics
	// This would parse pmd-stats-show output for additional details
	// Implementation depends on the specific format of stats output
}

// GetDropCounters retrieves specific drop counters from coverage
func (e *Exporter) GetDropCounters() (map[string]uint64, error) {
	cmd := exec.Command("ovs-appctl", "coverage/show")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get coverage: %w", err)
	}
	
	dropCounters := make(map[string]uint64)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	
	// List of specific drop counters to extract
	dropTypes := []string{
		"datapath_drop_upcall_error",
		"datapath_drop_lock_error",
		"datapath_drop_rx_invalid_packet",
		"datapath_drop_meter",
		"datapath_drop_userspace_action_error",
		"datapath_drop_tunnel_push_error",
		"datapath_drop_tunnel_pop_error",
		"datapath_drop_recirc_error",
		"datapath_drop_invalid_port",
		"datapath_drop_invalid_tnl_port",
		"datapath_drop_sample_error",
		"datapath_drop_nsh_decap_error",
		"drop_action_of_pipeline",
		"drop_action_bridge_not_found",
		"drop_action_recursion_too_deep",
		"drop_action_too_many_resubmit",
		"drop_action_stack_too_deep",
		"drop_action_no_recirculation",
		"drop_action_recirculation_conflict",
		"drop_action_too_many_mpls_labels",
		"drop_action_invalid_tunnel_metadata",
		"drop_action_unsupported_packet_type",
		"drop_action_congestion",
		"drop_action_forwarding_disabled",
	}
	
	// Create a map for quick lookup
	dropTypeMap := make(map[string]bool)
	for _, dt := range dropTypes {
		dropTypeMap[dt] = true
	}
	
	// Pattern to match coverage lines
	coverageRe := regexp.MustCompile(`^(\S+)\s+(\d+)`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := coverageRe.FindStringSubmatch(line); matches != nil {
			eventName := matches[1]
			if dropTypeMap[eventName] {
				if val, err := strconv.ParseUint(matches[2], 10, 64); err == nil {
					dropCounters[eventName] = val
				}
			}
		}
	}
	
	return dropCounters, nil
}