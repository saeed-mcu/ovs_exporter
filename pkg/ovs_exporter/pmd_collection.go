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
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// CollectPMDMetrics collects all PMD-related metrics
func (e *Exporter) CollectPMDMetrics() {
	level.Debug(e.logger).Log(
		"msg", "Collecting enhanced PMD performance metrics",
		"system_id", e.Client.System.ID,
	)
	
	// Collect enhanced PMD metrics
	enhancedMetrics, err := e.GetEnhancedPmdMetrics()
	if err != nil {
		level.Debug(e.logger).Log(
			"msg", "Enhanced PMD metrics collection failed",
			"system_id", e.Client.System.ID,
			"error", err.Error(),
		)
		// Fall back to basic metrics
		e.collectBasicPMDMetrics()
		return
	}
	
	if len(enhancedMetrics) == 0 {
		level.Debug(e.logger).Log(
			"msg", "No PMD metrics available (likely non-DPDK deployment)",
			"system_id", e.Client.System.ID,
		)
		return
	}
	
	for _, pmd := range enhancedMetrics {
		// CPU Utilization (convert from percentage to ratio)
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdCPUUtilization,
			prometheus.GaugeValue,
			pmd.CPUUtilization / 100.0, // Convert percentage to ratio (0-1)
			e.Client.System.ID, pmd.PmdID, pmd.NumaID, pmd.CoreID,
		))
		
		// Idle and Sleep metrics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdIdleCycles,
			prometheus.CounterValue,
			float64(pmd.IdleCycles),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdSleepIterations,
			prometheus.CounterValue,
			float64(pmd.SleepIterations),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// Core performance metrics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdCyclesPerIteration,
			prometheus.GaugeValue,
			pmd.CyclesPerIteration,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdPacketsPerIteration,
			prometheus.GaugeValue,
			pmd.PacketsPerIteration,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdCyclesPerPacket,
			prometheus.GaugeValue,
			pmd.CyclesPerPacket,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdPacketsPerBatch,
			prometheus.GaugeValue,
			pmd.PacketsPerBatch,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// RX Batch Statistics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdRxBatches,
			prometheus.CounterValue,
			float64(pmd.RxBatches),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdRxPackets,
			prometheus.CounterValue,
			float64(pmd.RxPackets),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdAvgRxBatchSize,
			prometheus.GaugeValue,
			pmd.AvgRxBatchSize,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdMaxRxBatchSize,
			prometheus.GaugeValue,
			float64(pmd.MaxRxBatchSize),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// TX Batch Statistics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdTxBatches,
			prometheus.CounterValue,
			float64(pmd.TxBatches),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdTxPackets,
			prometheus.CounterValue,
			float64(pmd.TxPackets),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdAvgTxBatchSize,
			prometheus.GaugeValue,
			pmd.AvgTxBatchSize,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// vHost Queue Metrics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdMaxVhostQueueLength,
			prometheus.GaugeValue,
			float64(pmd.MaxVhostQueueLength),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdAvgVhostQueueLength,
			prometheus.GaugeValue,
			pmd.AvgVhostQueueLength,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdVhostQueueFull,
			prometheus.CounterValue,
			float64(pmd.VhostQueueFull),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// Upcalls
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdUpcalls,
			prometheus.CounterValue,
			float64(pmd.Upcalls),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdUpcallCycles,
			prometheus.CounterValue,
			float64(pmd.UpcallCycles),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// vHost TX metrics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			vhostTxRetries,
			prometheus.CounterValue,
			float64(pmd.VhostTxRetries),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			vhostTxContention,
			prometheus.CounterValue,
			float64(pmd.VhostTxContention),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			vhostTxIrqs,
			prometheus.CounterValue,
			float64(pmd.VhostTxIrqs),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// Iterations
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdIterations,
			prometheus.CounterValue,
			float64(pmd.Iterations),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdBusyCycles,
			prometheus.CounterValue,
			float64(pmd.BusyCycles),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// Hit/Miss Statistics
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdExactMatchHit,
			prometheus.CounterValue,
			float64(pmd.ExactMatchHit),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdMaskedHit,
			prometheus.CounterValue,
			float64(pmd.MaskedHit),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdMiss,
			prometheus.CounterValue,
			float64(pmd.Miss),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdLost,
			prometheus.CounterValue,
			float64(pmd.Lost),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		// Suspicious Iterations
		if pmd.SuspiciousIterations > 0 {
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				pmdSuspiciousIterations,
				prometheus.CounterValue,
				float64(pmd.SuspiciousIterations),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))

			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				pmdSuspiciousPercent,
				prometheus.GaugeValue,
				pmd.SuspiciousPercent / 100.0, // Convert percentage to ratio (0-1)
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
		}
		
		// Flow Cache Metrics
		if pmd.EMCHitRate > 0 || pmd.EMCHits > 0 {
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				emcHitRate,
				prometheus.GaugeValue,
				pmd.EMCHitRate / 100.0, // Convert percentage to ratio (0-1)
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
			
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				emcHits,
				prometheus.CounterValue,
				float64(pmd.EMCHits),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
			
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				emcInserts,
				prometheus.CounterValue,
				float64(pmd.EMCInserts),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
		}
		
		if pmd.SMCHitRate > 0 || pmd.SMCHits > 0 {
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				smcHitRate,
				prometheus.GaugeValue,
				pmd.SMCHitRate / 100.0, // Convert percentage to ratio (0-1)
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
			
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				smcHits,
				prometheus.CounterValue,
				float64(pmd.SMCHits),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
		}
		
		if pmd.MegaflowHitRate > 0 || pmd.MegaflowHits > 0 || pmd.MegaflowMisses > 0 {
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				megaflowHitRate,
				prometheus.GaugeValue,
				pmd.MegaflowHitRate / 100.0, // Convert percentage to ratio (0-1)
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
			
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				megaflowHits,
				prometheus.CounterValue,
				float64(pmd.MegaflowHits),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
			
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				megaflowMisses,
				prometheus.CounterValue,
				float64(pmd.MegaflowMisses),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
		}
		
		if pmd.FlowCacheLookups > 0 {
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				flowCacheLookups,
				prometheus.CounterValue,
				float64(pmd.FlowCacheLookups),
				e.Client.System.ID, pmd.PmdID, pmd.NumaID,
			))
		}
	}
	
	level.Debug(e.logger).Log(
		"msg", "Enhanced PMD metrics collected successfully",
		"system_id", e.Client.System.ID,
		"pmd_count", len(enhancedMetrics),
	)
	
	// Collect specific drop counters
	e.collectDropCounters()
}

// collectBasicPMDMetrics falls back to basic PMD metrics collection
func (e *Exporter) collectBasicPMDMetrics() {
	pmdMetrics, err := e.GetPmdPerfMetrics()
	if err != nil || len(pmdMetrics) == 0 {
		return
	}
	
	for _, pmd := range pmdMetrics {
		// Add basic metrics as before
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdCyclesPerIteration,
			prometheus.GaugeValue,
			pmd.CyclesPerIteration,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdPacketsPerIteration,
			prometheus.GaugeValue,
			pmd.PacketsPerIteration,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdCyclesPerPacket,
			prometheus.GaugeValue,
			pmd.CyclesPerPacket,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdPacketsPerBatch,
			prometheus.GaugeValue,
			pmd.PacketsPerBatch,
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdMaxVhostQueueLength,
			prometheus.GaugeValue,
			float64(pmd.MaxVhostQueueLength),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdUpcalls,
			prometheus.CounterValue,
			float64(pmd.Upcalls),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdUpcallCycles,
			prometheus.CounterValue,
			float64(pmd.UpcallCycles),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			vhostTxRetries,
			prometheus.CounterValue,
			float64(pmd.TxRetries),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			vhostTxContention,
			prometheus.CounterValue,
			float64(pmd.TxContention),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			vhostTxIrqs,
			prometheus.CounterValue,
			float64(pmd.TxIrqs),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdIterations,
			prometheus.CounterValue,
			float64(pmd.Iterations),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
		
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			pmdBusyCycles,
			prometheus.CounterValue,
			float64(pmd.BusyCycles),
			e.Client.System.ID, pmd.PmdID, pmd.NumaID,
		))
	}
}

// collectDropCounters collects specific drop counter metrics
func (e *Exporter) collectDropCounters() {
	dropCounters, err := e.GetDropCounters()
	if err != nil {
		level.Debug(e.logger).Log(
			"msg", "Failed to collect drop counters",
			"system_id", e.Client.System.ID,
			"error", err.Error(),
		)
		return
	}
	
	for dropReason, count := range dropCounters {
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			datapathDrops,
			prometheus.CounterValue,
			float64(count),
			e.Client.System.ID, dropReason,
		))
	}
	
	if len(dropCounters) > 0 {
		level.Debug(e.logger).Log(
			"msg", "Drop counters collected",
			"system_id", e.Client.System.ID,
			"counter_count", len(dropCounters),
		)
	}
}