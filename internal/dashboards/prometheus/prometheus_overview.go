package prometheus

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	"github.com/nicolastakashi/community-perses-dashboards/pkg/panels"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withPrometheusOverviewStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Prometheus Stats",
		panelgroup.PanelsPerLine(1),
		panels.PrometheusStatsTable(datasource, labelMatcher),
	)
}

func withPrometheusOverviewDiscoveryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Discovery",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusTargetSync(datasource, labelMatcher),
		panels.PrometheusTargets(datasource, labelMatcher),
	)
}

func withPrometheusRetrievalGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Retrieval",
		panelgroup.PanelsPerLine(3),
		panels.PrometheusAverageScrapeIntervalDuration(datasource, labelMatcher),
		panels.PrometheusScrapeFailures(datasource, labelMatcher),
		panels.PrometheusAppendedSamples(datasource, labelMatcher),
	)
}

func withPrometheusStorageGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusHeadSeries(datasource, labelMatcher),
		panels.PrometheusHeadChunks(datasource, labelMatcher),
	)
}

func withPrometheusQueryGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusQueryRate(datasource, labelMatcher),
		panels.PrometheusQueryStateDuration(datasource, labelMatcher),
	)
}

func BuildPrometheusOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("prometheus-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Prometheus / Overview"),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("prometheus_build_info{}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("job"),
			),
		),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "prometheus_build_info"),
		dashboard.AddVariable("instance",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("instance",
					labelValuesVar.Matchers("prometheus_build_info{job='$job'}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("instance"),
				listVar.AllowAllValue(true),
			),
		),
		withPrometheusOverviewStatsGroup(datasource, clusterLabelMatcher),
		withPrometheusOverviewDiscoveryGroup(datasource, clusterLabelMatcher),
		withPrometheusRetrievalGroup(datasource, clusterLabelMatcher),
		withPrometheusStorageGroup(datasource, clusterLabelMatcher),
		withPrometheusQueryGroup(datasource, clusterLabelMatcher),
	)
}
