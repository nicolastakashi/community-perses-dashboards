package nodeexporter

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	panels "github.com/nicolastakashi/community-perses-dashboards/pkg/panels/node_exporter"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withNodeExporterNodesCPU(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU",
		panelgroup.PanelsPerLine(2),
		panels.NodeCPUUsagePercentage(datasource, labelMatcher),
		panels.NodeAverage(datasource, labelMatcher),
	)
}

func withNodeExporterNodesMemory(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory",
		panelgroup.PanelsPerLine(2),
		panels.NodeMemoryUsageBytes(datasource, labelMatcher),
		panels.NodeMemoryUsagePercentage(datasource, labelMatcher),
	)
}

func withNodeExporterNodesDisk(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk",
		panelgroup.PanelsPerLine(2),
		panels.NodeDiskIOBytes(datasource, labelMatcher),
		panels.NodeDiskIOSeconds(datasource, labelMatcher),
	)
}

func withNodeExporterNodesNetwork(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network",
		panelgroup.PanelsPerLine(2),
		panels.NodeNetworkReceivedBytes(datasource, labelMatcher),
		panels.NodeNetworkTransmitedBytes(datasource, labelMatcher),
	)
}

func BuildNodeExporterNodes(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("node-exporter-nodes",
		dashboard.ProjectName(project),
		dashboard.Name("Node Exporter / Nodes"),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "node_uname_info{job='node', sysname!='Darwin'}"),
		dashboard.AddVariable("instance",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("instance",
					dashboards.AddVariableDatasource(datasource),
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"node_uname_info{job='node', sysname!='Darwin'}",
							[]promql.LabelMatcher{clusterLabelMatcher},
						)),
				),
				listVar.DisplayName("instance"),
				listVar.AllowAllValue(true),
			),
		),
		withNodeExporterNodesCPU(datasource, clusterLabelMatcher),
		withNodeExporterNodesMemory(datasource, clusterLabelMatcher),
		withNodeExporterNodesDisk(datasource, clusterLabelMatcher),
		withNodeExporterNodesNetwork(datasource, clusterLabelMatcher),
	)
}
