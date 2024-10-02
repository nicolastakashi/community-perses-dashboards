package prometheus

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

// PrometheusStatsTable creates a panel group option for displaying Prometheus statistics in a table format.
// The table includes columns for job, instance, and version, and hides the value and timestamp columns.
// It uses the Prometheus metric `prometheus_build_info` to count instances by job, instance, and version.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the Prometheus query.
// - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
// - panelgroup.Option: An option to add the configured panel to a panel group.
func PrometheusStatsTable(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Prometheus Stats",
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "job",
					Header: "Job",
				},
				{
					Name:   "instance",
					Header: "Instance",
				},
				{
					Name:   "version",
					Header: "Version",
				},
				{
					Name: "value",
					Hide: true,
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("count by (job, instance, version) (prometheus_build_info{job=~'$job', instance=~'$instance'})", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

// PrometheusTargetSync creates a panel option for monitoring Prometheus target synchronization.
// It adds a time series panel with specific configurations for the Y-axis and legend, and includes a PromQL query.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the query.
// - labelMathers: A variadic parameter for PromQL label matchers.
//
// The function uses the following Prometheus metric:
// - prometheus_target_sync_length_seconds_sum: This metric represents the total time taken for target synchronization in seconds.
//
// The panel displays the sum rate of the target synchronization length over a 5-minute interval, grouped by job, scrape_job, and instance.
func PrometheusTargetSync(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Target Sync",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: "seconds",
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum(rate(prometheus_target_sync_length_seconds_sum{job=~'$job',instance=~'$instance'}[5m])) by (job, scrape_job, instance)", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}:{{instance}}:{{scrape_job}}"),
			),
		),
	)
}

// PrometheusTargets creates a panel group option for displaying Prometheus targets.
// It adds a time series panel with a legend positioned at the bottom in table mode.
// The panel includes a PromQL query that sums discovered targets by job and instance,
// with optional label matchers for filtering.
//
// Parameters:
// - datasourceName: The name of the Prometheus datasource.
// - labelMathers: Optional variadic parameter for PromQL label matchers.
//
// Metrics Used:
// - prometheus_sd_discovered_targets: This metric provides information about the targets discovered by Prometheus service discovery.
//
// Returns:
// - panelgroup.Option: An option to add the configured panel to a panel group.
func PrometheusTargets(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Targets",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (job, instance) (prometheus_sd_discovered_targets{job=~'$job',instance=~'$instance'})", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}:{{instance}}"),
			),
		),
	)
}

// PrometheusAverageScrapeIntervalDuration creates a panel option for displaying the average scrape interval duration
// for Prometheus targets. It uses the following Prometheus metrics:
// - prometheus_target_interval_length_seconds_sum: The sum of the target interval lengths in seconds.
// - prometheus_target_interval_length_seconds_count: The count of the target interval lengths.
//
// The function accepts a datasource name and an optional list of PromQL label matchers to filter the metrics.
//
// Parameters:
// - datasourceName: The name of the Prometheus datasource.
// - labelMathers: Optional PromQL label matchers to filter the metrics.
//
// Returns:
// - panelgroup.Option: A panel option configured to display the average scrape interval duration.
func PrometheusAverageScrapeIntervalDuration(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Scrape Interval Duration",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: "seconds",
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("rate(prometheus_target_interval_length_seconds_sum{job=~'$job',instance=~'$instance'}[5m]) / rate(prometheus_target_interval_length_seconds_count{job=~'$job',instance=~'$instance'}[5m])",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}:{{instance}} {{interval}} configured"),
			),
		),
	)
}

// PrometheusScrapeFailures creates a panel group option for displaying Prometheus scrape failure metrics.
// It generates a time series panel with multiple queries to visualize different types of scrape failures.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the queries.
// - labelMathers: Optional PromQL label matchers to filter the metrics.
//
// The following Prometheus metrics are used:
// - prometheus_target_scrapes_exceeded_body_size_limit_total: Number of times a scrape exceeded the body size limit.
// - prometheus_target_scrapes_exceeded_sample_limit_total: Number of times a scrape exceeded the sample limit.
// - prometheus_target_scrapes_sample_duplicate_timestamp_total: Number of times a scrape had duplicate timestamps.
// - prometheus_target_scrapes_sample_out_of_bounds_total: Number of times a scrape had samples out of bounds.
// - prometheus_target_scrapes_sample_out_of_order_total: Number of times a scrape had samples out of order.
//
// Each metric is aggregated by job and instance, and the rate is calculated over a 1-minute interval.
func PrometheusScrapeFailures(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Scrape failures",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (job, instance) (rate(prometheus_target_scrapes_exceeded_body_size_limit_total{job=~'$job',instance=~'$instance'}[1m]))", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("exceeded body size limit: {{job}} {{instance}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (job, instance) (rate(prometheus_target_scrapes_exceeded_sample_limit_total{job=~'$job',instance=~'$instance'}[1m]))", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("exceeded sample limit: {{job}} {{instance}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (job, instance) (rate(prometheus_target_scrapes_sample_duplicate_timestamp_total{job=~'$job',instance=~'$instance'}[1m]))", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("duplicate timestamp: {{job}} {{instance}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (job, instance) (rate(prometheus_target_scrapes_sample_out_of_bounds_total{job=~'$job',instance=~'$instance'}[1m]))", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("out of bounds: {{job}} {{instance}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (job, instance) (rate(prometheus_target_scrapes_sample_out_of_order_total{job=~'$job',instance=~'$instance'}[1m]))", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("out of order: {{job}} {{instance}}"),
			),
		),
	)
}

// PrometheusAppendedSamples creates a panel option for visualizing the rate of samples appended to Prometheus' TSDB head over a 5-minute interval.
// It uses the Prometheus metric `prometheus_tsdb_head_samples_appended_total` and allows for custom label matchers to filter the data.
// The panel includes a time series chart with a legend positioned at the bottom in table mode.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the query.
// - labelMathers: A variadic parameter for Prometheus label matchers to filter the metric data.
//
// Returns:
// - panelgroup.Option: A configured panel option for the appended samples visualization.
func PrometheusAppendedSamples(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Appended Samples",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("rate(prometheus_tsdb_head_samples_appended_total{job=~'$job',instance=~'$instance'}[5m])", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{instance}}"),
			),
		),
	)
}

// PrometheusHeadSeries creates a panel option for displaying the head series metric from Prometheus.
// The panel will show a time series chart with a legend positioned at the bottom in table mode.
//
// Parameters:
// - datasourceName: The name of the Prometheus datasource to be used for the query.
// - labelMathers: A variadic parameter of Prometheus label matchers to filter the query.
//
// The function queries the Prometheus metric `prometheus_tsdb_head_series` with the provided label matchers
// and formats the series name as "{{job}} {{instance}} head series".
//
// Returns:
// - panelgroup.Option: An option to add the configured panel to a panel group.
func PrometheusHeadSeries(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Head Series",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("prometheus_tsdb_head_series{job=~'$job',instance=~'$instance'}", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{instance}} head series"),
			),
		),
	)
}

// PrometheusHeadChunks creates a panel option for displaying the "Head Chunks" metric from Prometheus.
// It uses the `prometheus_tsdb_head_chunks` metric to show the number of head chunks in the TSDB.
// The panel includes a time series chart with a legend positioned at the bottom in table mode.
// The function accepts a datasource name and an optional list of PromQL label matchers.
//
// Parameters:
//   - datasourceName: The name of the Prometheus datasource.
//   - labelMathers: Optional PromQL label matchers to filter the metric.
//
// Returns:
//   - panelgroup.Option: An option to add the "Head Chunks" panel to a panel group.
func PrometheusHeadChunks(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Head Chunks",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("prometheus_tsdb_head_chunks{job=~'$job',instance=~'$instance'}", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{instance}} head chunks"),
			),
		),
	)
}

// PrometheusQueryRate creates a panel option for displaying the query rate of Prometheus engine query duration.
// It adds a time series panel with a legend positioned at the bottom in table mode.
// The panel includes a PromQL query that calculates the rate of the metric `prometheus_engine_query_duration_seconds_count`
// filtered by the provided label matchers and a fixed slice of 'inner_eval' over a 5-minute interval.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the query.
// - labelMathers: A variadic parameter of Prometheus label matchers to filter the query.
//
// Returns:
// - panelgroup.Option: An option to add the configured panel to a panel group.
func PrometheusQueryRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Query Rate",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("rate(prometheus_engine_query_duration_seconds_count{job=~'$job',instance=~'$instance',slice='inner_eval'}[5m])", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{instance}}"),
			),
		),
	)
}

// PrometheusQueryStateDuration creates a panel option for displaying the stage duration
// of Prometheus queries. It uses the metric `prometheus_engine_query_duration_seconds`
// with a quantile of 0.9, filtered by job and instance labels. The panel displays the
// data in a time series chart with the y-axis formatted in seconds and the legend positioned
// at the bottom in table mode.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the query.
// - labelMathers: Optional PromQL label matchers to filter the query.
//
// Returns:
// - panelgroup.Option: A panel option configured with the specified settings.
func PrometheusQueryStateDuration(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Stage Duration",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: "seconds",
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("max by (slice) (prometheus_engine_query_duration_seconds{quantile='0.9', job=~'$job',instance=~'$instance'})", labelMathers),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{slice}}"),
			),
		),
	)
}

// PrometheusRemoteStorageTimestampLag creates a panel option for visualizing the timestamp lag
// between the highest timestamp in Prometheus remote storage and the highest sent
// timestamp in the remote storage queue. This can help in identifying delays or
// issues in data ingestion.
//
// Parameters:
//   - datasourceName: The name of the Prometheus datasource to be used for the query.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: An option that adds a panel to a panel group with the specified
//     configuration for visualizing timestamp lag.
func PrometheusRemoteStorageTimestampLag(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Timestamp Lag",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(prometheus_remote_storage_highest_timestamp_in_seconds{instance=~'$instance'} -  ignoring(remote_name, url) group_right(instance) (prometheus_remote_storage_queue_highest_sent_timestamp_seconds{instance=~'$instance', url='$url'} != 0))",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageRateLag creates a panel option for monitoring the rate lag of Prometheus remote storage.
// It generates a time series panel with a specific query to measure the lag between the highest timestamp in Prometheus
// remote storage and the highest sent timestamp in the remote storage queue.
//
// Parameters:
// - datasourceName: The name of the data source to be used in the query.
// - labelMathers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
// - panelgroup.Option: An option to add the configured panel to a panel group.
func PrometheusRemoteStorageRateLag(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate[5m]",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"clamp_min(rate(prometheus_remote_storage_highest_timestamp_in_seconds{instance=~'$instance'}[5m])  - ignoring (remote_name, url) group_right(instance) rate(prometheus_remote_storage_queue_highest_sent_timestamp_seconds{instance=~'$instance', url='$url'}[5m]), 0)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageSampleRate creates a panel option for visualizing the rate of Prometheus remote storage samples
// over a 5-minute interval. It displays the rate of incoming samples versus succeeded or dropped samples.
//
// Parameters:
//   - datasourceName: The name of the Prometheus datasource.
//   - labelMathers: A variadic parameter for Prometheus label matchers.
//
// Returns:
//   - panelgroup.Option: An option that adds the configured panel to a panel group.
func PrometheusRemoteStorageSampleRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate, in vs. succeeded or dropped [5m]",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(prometheus_remote_storage_samples_in_total{instance=~'$instance'}[5m]) - ignoring(remote_name, url) group_right(instance) (rate(prometheus_remote_storage_succeeded_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_total{instance=~'$instance', url='$url'}[5m])) - (rate(prometheus_remote_storage_dropped_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_dropped_total{instance=~'$instance', url='$url'}[5m]))",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageCurrentShards creates a panel option for displaying the current number of shards
// in Prometheus remote storage. It configures a time series panel with a legend at the bottom in table mode,
// and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageCurrentShards(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Shards",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_remote_storage_shards{instance=~'$instance', url='$url'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageDesiredShards creates a panel option for displaying the desired shards
// of Prometheus remote storage. It configures a time series panel with a legend positioned at the
// bottom in table mode and adds a PromQL query to fetch the desired shards metric.
//
// Parameters:
//   - datasourceName: The name of the Prometheus datasource to be used in the query.
//   - labelMathers: A variadic parameter of PromQL label matchers to filter the query.
//
// Returns:
//
//	A panelgroup.Option configured with the desired shards panel and query.
func PrometheusRemoteStorageDesiredShards(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Desired Shards",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_remote_storage_shards_desired{instance=~'$instance', url='$url'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageMaxShards creates a panel option for displaying the maximum number of shards
// in Prometheus remote storage. It configures a time series panel with a legend at the bottom in table mode,
// and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageMaxShards(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Max Shards",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_remote_storage_shards_max{instance=~'$instance', url='$url'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageMinShards creates a panel option for displaying the minimum number of shards
// in Prometheus remote storage. It configures a time series panel with a legend at the bottom in table mode,
// and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageMinShards(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Min Shards",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_remote_storage_shards_min{instance=~'$instance', url='$url'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageShardCapacity creates a panel option for displaying the shard capacity
// in Prometheus remote storage. It configures a time series panel with a legend at the bottom in table mode,
// and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageShardCapacity(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Shard Capacity",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_remote_storage_shard_capacity{instance=~'$instance', url='$url'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStoragePendingSamples creates a panel option for displaying the pending samples
// in Prometheus remote storage. It configures a time series panel with a legend at the bottom in table mode,
// and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStoragePendingSamples(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Pending Samples",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_remote_storage_pending_samples{instance=~'$instance', url='$url'} or prometheus_remote_storage_samples_pending{instance=~'$instance', url='$url'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusTSDBCurrentSegment creates a panel option for displaying the current segment
// of the Prometheus TSDB WAL (Write-Ahead Log). It configures a time series panel with a legend
// positioned at the bottom in table mode and adds a PromQL query to fetch the relevant metric.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusTSDBCurrentSegment(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TSDB Current Segment",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_tsdb_wal_segment_current{instance=~'$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// PrometheusRemoteWriteCurrentSegment creates a panel option for displaying the current segment
// of the Prometheus remote write WAL (Write-Ahead Log). It configures a time series panel with a legend
// positioned at the bottom in table mode and adds a PromQL query to fetch the relevant metric.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteWriteCurrentSegment(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Current Segment",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"prometheus_wal_watcher_current_segment{instance=~'$instance'}",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// PrometheusRemoteStorageDroppedSamplesRate creates a panel option for displaying the rate of dropped samples
// in Prometheus remote storage over a 5-minute interval. It configures a time series panel with a legend
// positioned at the bottom in table mode and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageDroppedSamplesRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Dropped Samples Rate",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(prometheus_remote_storage_dropped_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_dropped_total{instance=~'$instance', url='$url'}[5m])",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageFailedSamplesRate creates a panel option for displaying the rate of failed samples
// in Prometheus remote storage over a 5-minute interval. It configures a time series panel with a legend
// positioned at the bottom in table mode and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageFailedSamplesRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Failed Samples",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(prometheus_remote_storage_failed_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_failed_total{instance=~'$instance', url='$url'}[5m])",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageRetriedSamplesRate creates a panel option for displaying the rate of retried samples
// in Prometheus remote storage over a 5-minute interval. It configures a time series panel with a legend
// positioned at the bottom in table mode and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageRetriedSamplesRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Retried Samples",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(prometheus_remote_storage_retried_samples_total{instance=~'$instance', url=~'$url'}[5m]) or rate(prometheus_remote_storage_samples_retried_total{instance=~'$instance', url=~'$url'}[5m])",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}

// PrometheusRemoteStorageEnqueueRetriesRate creates a panel option for displaying the rate of enqueue retries
// in Prometheus remote storage over a 5-minute interval. It configures a time series panel with a legend
// positioned at the bottom in table mode and adds a PromQL query to fetch the relevant metrics.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: A variadic list of PromQL label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageEnqueueRetriesRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Enqueue Retries",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(prometheus_remote_storage_enqueue_retries_total{instance=~'$instance', url=~'$url'}[5m])",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
			),
		),
	)
}
