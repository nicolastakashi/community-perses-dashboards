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
