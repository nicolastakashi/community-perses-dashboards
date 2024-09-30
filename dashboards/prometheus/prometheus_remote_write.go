package prometheus

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
	"github.com/perses/perses/go-sdk/prometheus/query"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func BuildPrometheusRemoteWrite(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	return dashboard.New("prometheus-remote-write",
		dashboard.Name("Prometheus / Remote Write"),
		dashboard.ProjectName(project),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "prometheus_remote_storage_shards"),
		dashboard.AddVariable("instance",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("instance",
					labelValuesVar.Matchers(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shards",
							"=",
							clusterLabelName,
							"$cluster",
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("instance"),
			),
		),
		dashboard.AddVariable("url",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("url",
					labelValuesVar.Matchers(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shards{instance='$instance'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("url"),
			),
		),

		dashboard.AddPanelGroup("Timestamps",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("Highest Timestamp In vs. Highest Timestamp Sent",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"(prometheus_remote_storage_highest_timestamp_in_seconds{instance=~'$instance'} -  ignoring(remote_name, url) group_right(instance) (prometheus_remote_storage_queue_highest_sent_timestamp_seconds{instance=~'$instance', url='$url'} != 0))",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Rate[5m]",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"clamp_min(rate(prometheus_remote_storage_highest_timestamp_in_seconds{instance=~'$instance'}[5m])  - ignoring (remote_name, url) group_right(instance) rate(prometheus_remote_storage_queue_highest_sent_timestamp_seconds{instance=~'$instance', url='$url'}[5m]), 0)",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Samples",
			panelgroup.PanelsPerLine(1),
			panelgroup.AddPanel("Rate, in vs. succeeded or dropped [5m]",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"rate(prometheus_remote_storage_samples_in_total{instance=~'$instance'}[5m]) - ignoring(remote_name, url) group_right(instance) (rate(prometheus_remote_storage_succeeded_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_total{instance=~'$instance', url='$url'}[5m])) - (rate(prometheus_remote_storage_dropped_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_dropped_total{instance=~'$instance', url='$url'}[5m]))",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Shards",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("Current Shards",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shards{instance=~'$instance', url='$url'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Desired Shards",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shards_desired{instance=~'$instance', url='$url'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Max Shards",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shards_max{instance=~'$instance', url='$url'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Min Shards",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shards_min{instance=~'$instance', url='$url'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Shard Details",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("Shard Capacity",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_shard_capacity{instance=~'$instance', url='$url'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Pending Samples",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_remote_storage_pending_samples{instance=~'$instance', url='$url'} or prometheus_remote_storage_samples_pending{instance=~'$instance', url='$url'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Segments",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("TSDB Current Segment",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_tsdb_wal_segment_current{instance=~'$instance'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}"),
					),
				),
			),
			panelgroup.AddPanel("Remote Write Current Segment",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"prometheus_wal_watcher_current_segment{instance=~'$instance'}",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Misc. Rates",
			panelgroup.PanelsPerLine(4),
			panelgroup.AddPanel("Dropped Samples",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"rate(prometheus_remote_storage_dropped_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_dropped_total{instance=~'$instance', url='$url'}[5m])",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Failed Samples",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"rate(prometheus_remote_storage_failed_samples_total{instance=~'$instance', url='$url'}[5m]) or rate(prometheus_remote_storage_samples_failed_total{instance=~'$instance', url='$url'}[5m])",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Retried Samples",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"rate(prometheus_remote_storage_retried_samples_total{instance=~'$instance', url=~'$url'}[5m]) or rate(prometheus_remote_storage_samples_retried_total{instance=~'$instance', url=~'$url'}[5m])",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
			panelgroup.AddPanel("Retried Samples",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL(
						promql.LabelsSetPromQL(
							"rate(prometheus_remote_storage_enqueue_retries_total{instance=~'$instance', url=~'$url'}[5m])",
							"=",
							clusterLabelName,
							"$cluster",
						),
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{instance}}:{{remote_name}}:{{url}}"),
					),
				),
			),
		),
	)
}
