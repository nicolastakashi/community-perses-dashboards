package nodeexporter

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/panel/gauge"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
	"github.com/perses/perses/go-sdk/prometheus/query"
)

// NodeCPUUsagePercentage creates a panel option for displaying the CPU usage percentage
// of nodes using Prometheus as the data source. It generates a time series panel with
// specific configurations for the Y-axis format and legend position.
//
// The panel uses the following Prometheus metrics:
// - node_cpu_seconds_total: CPU time spent in different modes
//
// The panel shows:
// - CPU usage percentage excluding idle, iowait, and steal time
// - Breakdown by CPU core
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func NodeCPUUsagePercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		panel.Description("Shows CPU utilization percentage across cluster nodes"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"((1 - sum without (mode) (rate(node_cpu_seconds_total{job='node', mode=~'idle|iowait|steal', instance='$instance'}[5m]))) / ignoring(cpu) group_left count without (cpu, mode) (node_cpu_seconds_total{job='node', mode='idle', instance='$instance'}))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} - CPU - Usage"),
			),
		),
	)
}

// ClusterNodeCPUUsagePercentage creates a panel option for displaying the CPU usage percentage of cluster nodes.
//
// The panel uses the following Prometheus metrics:
// - instance:node_cpu_utilisation:rate5m: Rate of CPU utilization
// - instance:node_num_cpu:sum: Total number of CPUs
//
// The panel shows:
// - CPU usage percentage per instance
// - Utilization relative to total available CPU cores
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func ClusterNodeCPUUsagePercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"((instance:node_cpu_utilisation:rate5m{job='node'} * instance:node_num_cpu:sum{job='node'}) != 0 ) / scalar(sum(instance:node_num_cpu:sum{job='node'}))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeCPUSaturationPercentage creates a panel option for displaying the CPU saturation percentage.
// (Load1 per CPU) for cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a percentage format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used divides the instance load by the count of instances, filtering out zero values.
//
// The panel uses the following Prometheus metrics:
// - instance:node_load1_per_cpu:ratio: Load average per CPU
//
// The panel shows:
// - CPU saturation based on 1-minute load average
// - Load per CPU across cluster nodes
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeCPUSaturationPercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Saturation (Load1 per CPU)",
		panel.Description("Shows CPU saturation metrics across cluster nodes"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(instance:node_load1_per_cpu:ratio{job='node'} / scalar(count(instance:node_load1_per_cpu:ratio{job='node'})))  != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeMemoryUsagePercentage creates a panel option for displaying the memory usage percentage of cluster nodes.
//
// The panel uses the following Prometheus metrics:
// - instance:node_memory_utilisation:ratio: Memory utilization ratio
//
// The panel shows:
// - Memory usage percentage per instance
// - Relative utilization across cluster nodes
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func ClusterNodeMemoryUsagePercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Utilisation",
		panel.Description("Shows memory utilization percentage across cluster nodes"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(instance:node_memory_utilisation:ratio{job='node'} / scalar(count(instance:node_memory_utilisation:ratio{job='node'}))) != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeMemorySaturationPercentage creates a panel option for displaying memory saturation.
// (Major Page Faults) of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a reads per second format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used calculates the rate of major page faults over a 5-minute interval.
//
// The panel uses the following Prometheus metrics:
// - instance:node_vmstat_pgmajfault:rate5m: Rate of major page faults
//
// The panel shows:
// - Memory saturation through page fault rate
// - Page faults per second per node
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeMemorySaturationPercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Saturation (Major Page Faults)",
		panel.Description("Shows memory saturation through page fault metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.ReadsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"instance:node_vmstat_pgmajfault:rate5m{job='node'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeDiskUsagePercentage creates a panel option for displaying the disk usage percentage
// of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a percentage format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used calculates the disk usage percentage by dividing the rate of disk I/O time by the total number of disks.
//
// The panel uses the following Prometheus metrics:
// - instance_device:node_disk_io_time_seconds:rate5m: Rate of disk I/O time
//
// The panel shows:
// - Disk I/O utilization percentage
// - Usage patterns across different disks
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeDiskUsagePercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk IO Utilisation",
		panel.Description("Shows disk I/O utilization across cluster nodes"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(instance_device:node_disk_io_time_seconds:rate5m{job='node'} / scalar(count(instance_device:node_disk_io_time_seconds:rate5m{job='node'}))) != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeDiskSaturationPercentage creates a panel option for displaying the disk I/O saturation
// of cluster nodes.
// The panel uses the following Prometheus metrics:
// - instance_device:node_disk_io_time_seconds:rate5m: Rate of disk I/O time
//
// The panel shows:
// - Disk I/O saturation per device
// - Saturation patterns across nodes
func ClusterNodeDiskSaturationPercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk IO Saturation",
		panel.Description("Shows disk I/O saturation metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(instance_device:node_disk_io_time_seconds:rate5m{job='node'} / scalar(count(instance_device:node_disk_io_time_seconds:rate5m{job='node'}))) != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeDiskSpacePercentage creates a panel option for displaying the disk space utilization
// of cluster nodes.
// The panel uses the following Prometheus metrics:
// - instance:node_filesystem_avail_bytes:sum: Available filesystem space
// - instance:node_filesystem_size_bytes:sum: Total filesystem size
//
// The panel shows:
// - Disk space utilization percentage
// - Available vs total space ratio
func ClusterNodeDiskSpacePercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk Space Utilisation",
		panel.Description("Shows disk space utilization metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentDecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum without (device) (max without (fstype, mountpoint) ((node_filesystem_size_bytes{job='node', fstype!='', mountpoint!=''} - node_filesystem_avail_bytes{job='node', fstype!='', mountpoint!=''}) != 0)) / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{job='node', fstype!='', mountpoint!=''})))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// ClusterNodeNetworkSaturationBytes creates a panel option for displaying the network saturation
// (Drops Receive/Transmit) of cluster nodes.
// The panel uses the following Prometheus metrics:
// - instance:node_network_receive_drop_excluding_lo:rate5m: Rate of network receive drops
//
// The panel shows:
// - Network packet drops per interface
// - Drop rates across nodes
func ClusterNodeNetworkSaturationBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Saturation (Drops Receive/Transmit)",
		panel.Description("Shows network saturation through drop metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"instance:node_network_receive_drop_excluding_lo:rate5m{job='node'} != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Network - Received"),
			),
		),
	)
}

// ClusterNodeNetworkUsageBytes creates a panel option for displaying the network utilization
// (Bytes Receive/Transmit) of cluster nodes.
// The panel uses the following Prometheus metrics:
// - instance:node_network_receive_bytes_excluding_lo:rate5m: Network bytes received
// - instance:node_network_transmit_bytes_excluding_lo:rate5m: Network bytes transmitted
//
// The panel shows:
// - Network throughput per interface
// - Receive and transmit rates
func ClusterNodeNetworkUsageBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Utilisation (Bytes Receive/Transmit)",
		panel.Description("Shows network utilization in bytes"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"instance:node_network_receive_bytes_excluding_lo:rate5m{job='node'} != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Network - Received"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"instance:node_network_transmit_bytes_excluding_lo:rate5m{job='node'} != 0",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Network - Transmitted"),
			),
		),
	)
}

// NodeAverage creates a panel group option for displaying CPU usage metrics.
//
// The panel uses the following Prometheus metrics:
// - node_load1: 1-minute load average
// - node_load5: 5-minute load average
// - node_load15: 15-minute load average
// - node_cpu_seconds_total: CPU time in different modes
//
// The panel shows:
// - Load averages over different time periods
// - Number of logical CPU cores
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the queries.
//   - labelMatchers: Optional Prometheus label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel group option.
func NodeAverage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		panel.Description("Shows CPU utilization metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_load1{job='node', instance='$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("CPU - 1m Average"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_load5{job='node', instance='$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("CPU - 5m Average"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_load15{job='node', instance='$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("CPU - 15m Average"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count(node_cpu_seconds_total{job='node', instance='$instance', mode='idle'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("CPU - Logical Cores"),
			),
		),
	)
}

// NodeMemoryUsageBytes creates a panel group option for displaying node memory usage in bytes.
// The panel uses the following Prometheus metrics:
// - node_memory_Buffers_bytes: Memory used for buffers
// - node_memory_Cached_bytes: Memory used for cache
// - node_memory_MemFree_bytes: Free memory available
//
// The panel shows:
// - Memory usage breakdown by type (buffers, cached, free)
// - Values in bytes
//
// Parameters:
// - datasourceName: The name of the data source to be used for the queries.
// - labelMatchers: Optional Prometheus label matchers to filter the queries.
//
// Returns:
// - panelgroup.Option: The configured panel group option.
func NodeMemoryUsageBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
		panel.Description("Shows memory utilization metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:        string(commonSdk.BytesUnit),
					ShortValues: true,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_memory_Buffers_bytes{job='node', instance='$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Memory - Buffers"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_memory_Cached_bytes{job='node', instance='$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Memory - Cached"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_memory_MemFree_bytes{job='node', instance='$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Memory - Free"),
			),
		),
	)
}

// NodeMemoryUsagePercentage creates a panel option for displaying memory usage percentage.
// The panel uses the following Prometheus metrics:
// - node_memory_MemAvailable_bytes: Available memory in bytes
// - node_memory_MemTotal_bytes: Total physical memory in bytes
//
// The panel shows:
// - Memory usage percentage with thresholds
// - Available vs total memory ratio
//
// Parameters:
// - datasourceName: The name of the Prometheus datasource.
// - labelMatchers: Optional Prometheus label matchers to filter the metrics.
//
// Returns:
// - panelgroup.Option: The panel option for the memory usage gauge chart.
func NodeMemoryUsagePercentage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
		panel.Description("Shows memory utilization across nodes"),
		gauge.Chart(
			gauge.Calculation(commonSdk.LastCalculation),
			gauge.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentMode),
			}),
			gauge.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "orange",
						Value: 80,
					},
					{
						Color: "red",
						Value: 90,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"100 - (avg(node_memory_MemAvailable_bytes{job='node', instance='$instance'}) / avg(node_memory_MemTotal_bytes{job='node', instance='$instance'}) * 100)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Memory - Usage"),
			),
		),
	)
}

// NodeDiskIOBytes creates a panel option for displaying Disk I/O metrics.
//
// The panel uses the following Prometheus metrics:
// - node_disk_read_bytes_total: Total number of bytes read from disk
// - node_disk_io_time_seconds_total: Total seconds spent on I/O operations
//
// The panel shows:
// - Rate of bytes read from disk per device
// - I/O time per device
//
// Parameters:
//   - datasourceName: The name of the data source.
//   - labelMatchers: Optional Prometheus label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func NodeDiskIOBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk I/O Bytes",
		panel.Description("Shows disk I/O metrics in bytes"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(node_disk_read_bytes_total{job='node', instance='$instance',device!=''}[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} - Disk - Usage"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(node_disk_io_time_seconds_total{job='node', instance='$instance',device!=''}[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} - Disk - Written"),
			),
		),
	)
}

// NodeDiskIOSeconds creates a panel option for displaying Disk I/O time series data.
//
// The panel uses the following Prometheus metrics:
// - node_disk_io_time_seconds_total: Total time spent on I/O operations
//
// The panel shows:
// - I/O operation duration per device
// - Time spent on disk operations
func NodeDiskIOSeconds(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk I/O Seconds",
		panel.Description("Shows disk I/O duration metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.SecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(node_disk_io_time_seconds_total{job='node', instance='$instance',device!=''}[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} - Disk - IO Time"),
			),
		),
	)
}

// NodeNetworkReceivedBytes creates a panel option for displaying the rate of network received bytes
// for nodes, excluding the loopback device.
//
// The panel uses the following Prometheus metrics:
// - node_network_receive_bytes_total: Total bytes received over network
//
// The panel shows:
// - Network receive rate per interface
// - Bandwidth utilization patterns
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMatchers: Optional Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func NodeNetworkReceivedBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Received",
		panel.Description("Shows network received bytes metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(node_network_receive_bytes_total{job='node', instance='$instance',device!='lo'}[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} - Network - Received"),
			),
		),
	)
}

// NodeNetworkTransmitedBytes creates a panel option for displaying the network transmitted bytes
// for nodes.
//
// The panel uses the following Prometheus metrics:
// - node_network_transmit_bytes_total: Total bytes transmitted over network
//
// The panel shows:
// - Network transmit rate per interface
// - Bandwidth utilization patterns
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMatchers: Optional Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func NodeNetworkTransmitedBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Transmitted",
		panel.Description("Shows network transmitted bytes metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(node_network_transmit_bytes_total{job='node', instance='$instance',device!='lo'}[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} - Network - Transmitted"),
			),
		),
	)
}
