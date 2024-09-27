package prometheus

import (
	"github.com/nicolastakashi/community-perses-dashboards/dashboards"
	"github.com/perses/perses/go-sdk"
	"github.com/perses/perses/go-sdk/dashboard"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func BuildPrometheusOverview(exec sdk.Exec, project string, datasource string, clusterLabelName string) {
	builder, buildErr := dashboard.New("prometheus-overview",
		dashboard.Name("Prometheus / Overview"),
		dashboard.ProjectName(project),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("prometheus_build_info{}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("job"),
			),
		),
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("prometheus_build_info{job='$job'}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("cluster"),
			),
		),
		dashboard.AddVariable("instance",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("instance",
					labelValuesVar.Matchers("prometheus_build_info{job='$job'}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("instance"),
			),
		),

		dashboard.AddPanelGroup("Prometheus Stats",
			panelgroup.PanelsPerLine(1),
			panelgroup.AddPanel("Prometheus Stats",
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
					query.PromQL("count by (job, instance, version) (prometheus_build_info{job=~'$job', instance=~'$instance'})",
						dashboards.AddQueryDataSource(datasource),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Discovery",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("Target Sync",
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
					query.PromQL("sum(rate(prometheus_target_sync_length_seconds_sum{job=~'$job',instance=~'$instance'}[5m])) by (job, scrape_job, instance)",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}}:{{instance}}:{{scrape_job}}"),
					),
				),
			),
			panelgroup.AddPanel("Targets",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL("sum by (job, instance) (prometheus_sd_discovered_targets{job=~'$job',instance=~'$instance'})",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}}:{{instance}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Retrieval",
			panelgroup.PanelsPerLine(3),
			panelgroup.AddPanel("Average Scrape Interval Duration",
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
					query.PromQL("rate(prometheus_target_interval_length_seconds_sum{job=~'$job',instance=~'$instance'}[5m]) / rate(prometheus_target_interval_length_seconds_count{job=~'$job',instance=~'$instance'}[5m])",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}}:{{instance}} {{interval}} configured"),
					),
				),
			),
			panelgroup.AddPanel("Scrape failures",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL("sum by (job, instance) (rate(prometheus_target_scrapes_exceeded_body_size_limit_total{job=~'$job',instance=~'$instance'}[1m]))",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("exceeded body size limit: {{job}} {{instance}}"),
					),
				),
				panel.AddQuery(
					query.PromQL("sum by (job, instance) (rate(prometheus_target_scrapes_exceeded_sample_limit_total{job=~'$job',instance=~'$instance'}[1m]))",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("exceeded sample limit: {{job}} {{instance}}"),
					),
				),
				panel.AddQuery(
					query.PromQL("sum by (job, instance) (rate(prometheus_target_scrapes_sample_duplicate_timestamp_total{job=~'$job',instance=~'$instance'}[1m]))",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("duplicate timestamp: {{job}} {{instance}}"),
					),
				),
				panel.AddQuery(
					query.PromQL("sum by (job, instance) (rate(prometheus_target_scrapes_sample_out_of_bounds_total{job=~'$job',instance=~'$instance'}[1m]))",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("out of bounds: {{job}} {{instance}}"),
					),
				),
				panel.AddQuery(
					query.PromQL("sum by (job, instance) (rate(prometheus_target_scrapes_sample_out_of_order_total{job=~'$job',instance=~'$instance'}[1m]))",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("out of order: {{job}} {{instance}}"),
					),
				),
			),
			panelgroup.AddPanel("Appended Samples",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL("rate(prometheus_tsdb_head_samples_appended_total{job=~'$job',instance=~'$instance'}[5m])",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}} {{instance}}"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Storage",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("Head Series",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL("prometheus_tsdb_head_series{job=~'$job',instance=~'$instance'}",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}} {{instance}} head series"),
					),
				),
			),
			panelgroup.AddPanel("Head Chunks",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL("prometheus_tsdb_head_chunks{job=~'$job',instance=~'$instance'}",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}} {{instance}} head chunks"),
					),
				),
			),
		),

		dashboard.AddPanelGroup("Query",
			panelgroup.PanelsPerLine(2),
			panelgroup.AddPanel("Query Rate",
				timeSeriesPanel.Chart(
					timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
						Position: timeSeriesPanel.BottomPosition,
						Mode:     timeSeriesPanel.TableMode,
					}),
				),
				panel.AddQuery(
					query.PromQL("rate(prometheus_engine_query_duration_seconds_count{job=~'$job',instance=~'$instance',slice='inner_eval'}[5m])",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{job}} {{instance}}"),
					),
				),
			),
			panelgroup.AddPanel("Stage Duration",
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
					query.PromQL("max by (slice) (prometheus_engine_query_duration_seconds{quantile='0.9', job=~'$job',instance=~'$instance'})",
						dashboards.AddQueryDataSource(datasource),
						query.SeriesNameFormat("{{slice}}"),
					),
				),
			),
		),
	)
	exec.BuildDashboard(builder, buildErr)
}
