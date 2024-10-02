package prometheus

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	panels "github.com/nicolastakashi/community-perses-dashboards/pkg/panels/prometheus"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withPrometheusRwTimestamps(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Timestamps",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusRemoteStorageTimestampLag(datasource, labelMatcher),
		panels.PrometheusRemoteStorageRateLag(datasource, labelMatcher),
	)
}

func withPrometheusRwSamples(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Samples",
		panelgroup.PanelsPerLine(1),
		panels.PrometheusRemoteStorageSampleRate(datasource, labelMatcher),
	)
}

func withPrometheusRwShard(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Shards",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusRemoteStorageCurrentShards(datasource, labelMatcher),
		panels.PrometheusRemoteStorageDesiredShards(datasource, labelMatcher),
		panels.PrometheusRemoteStorageMaxShards(datasource, labelMatcher),
		panels.PrometheusRemoteStorageMinShards(datasource, labelMatcher),
	)
}

func withPrometheusRwShardDetails(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Shard Details",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusRemoteStorageShardCapacity(datasource, labelMatcher),
		panels.PrometheusRemoteStoragePendingSamples(datasource, labelMatcher),
	)
}

func withPrometheusRwSegments(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Segments",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusTSDBCurrentSegment(datasource, labelMatcher),
		panels.PrometheusRemoteWriteCurrentSegment(datasource, labelMatcher),
	)
}

func withPrometheusRwMiscRates(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Misc. Rates",
		panelgroup.PanelsPerLine(4),
		panels.PrometheusRemoteStorageDroppedSamplesRate(datasource, labelMatcher),
		panels.PrometheusRemoteStorageFailedSamplesRate(datasource, labelMatcher),
		panels.PrometheusRemoteStorageRetriedSamplesRate(datasource, labelMatcher),
		panels.PrometheusRemoteStorageEnqueueRetriesRate(datasource, labelMatcher),
	)
}

func BuildPrometheusRemoteWrite(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("prometheus-remote-write",
		dashboard.Name("Prometheus / Remote Write"),
		dashboard.ProjectName(project),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "prometheus_remote_storage_shards"),
		dashboard.AddVariable("instance",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("instance",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"prometheus_remote_storage_shards",
							[]promql.LabelMatcher{clusterLabelMatcher},
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
						promql.SetLabelMatchers(
							"prometheus_remote_storage_shards{instance='$instance'}",
							[]promql.LabelMatcher{clusterLabelMatcher},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("url"),
			),
		),
		withPrometheusRwTimestamps(datasource, clusterLabelMatcher),
		withPrometheusRwSamples(datasource, clusterLabelMatcher),
		withPrometheusRwShard(datasource, clusterLabelMatcher),
		withPrometheusRwShardDetails(datasource, clusterLabelMatcher),
		withPrometheusRwSegments(datasource, clusterLabelMatcher),
		withPrometheusRwMiscRates(datasource, clusterLabelMatcher),
	)
}
