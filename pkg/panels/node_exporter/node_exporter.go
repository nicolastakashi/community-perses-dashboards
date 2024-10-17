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
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func NodeCPUUsagePercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{cpu}}"),
			),
		),
	)
}

// ClusterNodeCPUUsagePercentage creates a panel option for displaying the CPU usage percentage of cluster nodes.
// It takes a datasource name and an optional list of Prometheus label matchers as arguments.
// The panel displays a time series chart with the CPU usage percentage, formatted as a percentage unit.
// The legend is positioned at the bottom and displayed in table mode, showing the last calculated value.
// The PromQL query used calculates the CPU usage percentage by dividing the rate of CPU utilization by the total number of CPUs.
//
// The following Prometheus metrics are used:
// - instance:node_cpu_utilisation:rate5m: The rate of CPU utilization over a 5-minute interval.
// - instance:node_num_cpu:sum: The total number of CPUs.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeCPUUsagePercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeCPUSaturationPercentage creates a panel option for displaying the CPU saturation percentage
// (Load1 per CPU) for cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a percentage format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used divides the instance load by the count of instances, filtering out zero values.
//
// The following Prometheus metrics are used:
// - instance:node_load1_per_cpu:ratio: The ratio of the 1-minute load average per CPU.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeCPUSaturationPercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Saturation (Load1 per CPU)",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeMemoryUsagePercentage creates a panel option for displaying the memory usage percentage of cluster nodes.
// It takes a datasource name and an optional list of Prometheus label matchers as arguments.
// The panel displays a time series chart with the memory usage percentage, formatted as a percentage unit.
// The legend is positioned at the bottom and displayed in table mode, showing the last calculated value.
// The PromQL query used calculates the memory usage percentage by dividing the used memory by the total memory.
//
// The following Prometheus metrics are used:
// - instance:node_memory_utilisation:ratio: The ratio of memory utilization.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeMemoryUsagePercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Utilisation",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeMemorySaturationPercentage creates a panel option for displaying the memory saturation percentage
// (Major Page Faults) of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a reads per second format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used calculates the rate of major page faults over a 5-minute interval.
//
// The following Prometheus metrics are used:
// - instance:node_vmstat_pgmajfault:rate5m: The rate of major page faults over a 5-minute interval.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeMemorySaturationPercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Saturation (Major Page Faults)",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
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
// The following Prometheus metrics are used:
// - instance_device:node_disk_io_time_seconds:rate5m: The rate of disk I/O time over a 5-minute interval.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeDiskUsagePercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk IO Utilisation",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeDiskSaturationPercentage creates a panel option for displaying the disk I/O saturation
// of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a percentage format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used calculates the disk I/O saturation by dividing the weighted disk I/O time by the total number of disks.
//
// The following Prometheus metrics are used:
// - instance_device:node_disk_io_time_weighted_seconds:rate5m: The rate of weighted disk I/O time over a 5-minute interval.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeDiskSaturationPercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk IO Saturation",
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
					"(instance_device:node_disk_io_time_weighted_seconds:rate5m{job='node'} / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate5m{job='node'}))) != 0",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeDiskSpacePercentage creates a panel option for displaying the disk space utilization
// of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a percentage format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used calculates the disk space utilization by dividing the used disk space by the total disk space.
//
// The following Prometheus metrics are used:
// - node_filesystem_size_bytes: The total size of the filesystem.
// - node_filesystem_avail_bytes: The available size of the filesystem.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeDiskSpacePercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk Space Utilisation",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeNetworkSaturationBytes creates a panel option for displaying the network saturation
// (Drops Receive/Transmit) of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a bytes per second format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// query used calculates the rate of network drops over a 5-minute interval.
//
// The following Prometheus metrics are used:
// - instance:node_network_receive_drop_excluding_lo:rate5m: The rate of network receive drops over a 5-minute interval.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeNetworkSaturationBytes(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Saturation (Drops Receive/Transmit)",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }}"),
			),
		),
	)
}

// ClusterNodeNetworkUsageBytes creates a panel option for displaying the network utilization
// (Bytes Receive/Transmit) of cluster nodes. It takes a datasource name and an optional list of Prometheus label
// matchers as arguments. The panel includes a time series chart with a bytes per second format on the Y-axis
// and a legend positioned at the bottom in table mode, showing the last calculated value. The PromQL
// queries used calculate the rate of network bytes received and transmitted over a 5-minute interval.
//
// The following Prometheus metrics are used:
// - instance:node_network_receive_bytes_excluding_lo:rate5m: The rate of network bytes received over a 5-minute interval.
// - instance:node_network_transmit_bytes_excluding_lo:rate5m: The rate of network bytes transmitted over a 5-minute interval.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func ClusterNodeNetworkUsageBytes(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Utilisation (Bytes Receive/Transmit)",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }} Receive"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"instance:node_network_transmit_bytes_excluding_lo:rate5m{job='node'} != 0",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ instance }} Transmit"),
			),
		),
	)
}

// NodeAverage creates a panel group option for displaying CPU usage metrics.
// It adds a time series panel with multiple Prometheus queries to visualize
// the 1-minute, 5-minute, and 15-minute load averages, as well as the count
// of logical CPU cores.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the queries.
//   - labelMathers: Optional Prometheus label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel group option.
func NodeAverage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("1m load average"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_load5{job='node', instance='$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("5m load average"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_load15{job='node', instance='$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("15m load average"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count(node_cpu_seconds_total{job='node', instance='$instance', mode='idle'})",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("logical cores"),
			),
		),
	)
}

// NodeMemoryUsageBytes creates a panel group option for displaying node memory usage in bytes.
// It generates a time series panel with the following queries:
// - Memory used: Total memory minus free, buffers, and cached memory.
// - Memory buffers: Memory used for buffers.
// - Memory cached: Memory used for caching.
// - Memory free: Free memory.
//
// The panel includes a Y-axis formatted in bytes and a legend positioned at the bottom in table mode.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the queries.
// - labelMathers: Optional Prometheus label matchers to filter the queries.
//
// Returns:
// - panelgroup.Option: The configured panel group option.
func NodeMemoryUsageBytes(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
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
					"(node_memory_MemTotal_bytes{job='node',instance='$instance'}-node_memory_MemFree_bytes{job='node',instance='$instance'}-node_memory_Buffers_bytes{job='node',instance='$instance'}-node_memory_Cached_bytes{job='node',instance='$instance'})",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("memory used"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_memory_Buffers_bytes{job='node', instance='$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("memory buffers"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_memory_Cached_bytes{job='node', instance='$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("memory cached"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"node_memory_MemFree_bytes{job='node', instance='$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("memory free"),
			),
		),
	)
}

// NodeMemoryUsagePercentage creates a panel option for displaying the memory usage percentage of a node.
// The panel is a gauge chart that shows the memory usage with thresholds set at 80% (orange) and 90% (red).
// The function takes a datasource name and an optional list of Prometheus label matchers.
//
// The following Prometheus metrics are used:
// - node_memory_MemAvailable_bytes: The amount of memory available for starting new applications, without swapping.
// - node_memory_MemTotal_bytes: The total amount of physical memory.
//
// Parameters:
// - datasourceName: The name of the Prometheus datasource.
// - labelMathers: Optional Prometheus label matchers to filter the metrics.
//
// Returns:
// - panelgroup.Option: The panel option for the memory usage gauge chart.
func NodeMemoryUsagePercentage(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("memory used"),
			),
		),
	)
}

// NodeDiskIOBytes creates a panel group option for displaying Disk I/O metrics in a time series chart.
// The panel includes two Prometheus queries:
// 1. "rate(node_disk_read_bytes_total{job='node', instance='$instance',device!=”}[5m])" - This metric tracks the rate of bytes read from disk.
// 2. "rate(node_disk_io_time_seconds_total{job='node', instance='$instance',device!=”}[5m])" - This metric tracks the rate of I/O time in seconds.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the Prometheus queries.
// - labelMathers: Optional Prometheus label matchers to filter the metrics.
//
// Returns:
// - panelgroup.Option: A configured panel group option for the Disk I/O metrics.
func NodeDiskIOBytes(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk I/O Bytes",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} read"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(node_disk_io_time_seconds_total{job='node', instance='$instance',device!=''}[5m])",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} written"),
			),
		),
	)
}

// NodeDiskIOSeconds creates a panel option for displaying Disk I/O time series data.
//
// The panel queries the Prometheus metric `node_disk_io_time_seconds_total`,
// which measures the total time spent on I/O operations for each disk device.
// The query applies a rate function over a 5-minute interval and filters the data
// based on the provided label matchers.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the query.
// - labelMathers: Optional Prometheus label matchers to filter the query results.
//
// Returns:
// - A panelgroup.Option that adds the configured Disk I/O panel to a panel group.
func NodeDiskIOSeconds(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Disk I/O Seconds",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}} io time"),
			),
		),
	)
}

// NodeNetworkReceivedBytes creates a panel option for displaying the rate of network received bytes
// for nodes, excluding the loopback device. The panel is configured with a time series chart that
// shows the data in bytes per second, and includes a legend positioned at the bottom in table mode,
// displaying the last value.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: Optional Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func NodeNetworkReceivedBytes(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Received",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}}"),
			),
		),
	)
}

// NodeNetworkTransmitedBytes creates a panel option for displaying the network transmitted bytes
// for nodes. It configures a time series panel with specific Y-axis formatting, legend settings,
// and a Prometheus query to fetch the data.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: Optional Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func NodeNetworkTransmitedBytes(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Network Transmitted",
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
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{device}}"),
			),
		),
	)
}
