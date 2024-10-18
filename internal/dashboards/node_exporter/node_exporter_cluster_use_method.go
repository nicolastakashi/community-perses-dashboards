package nodeexporter

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	panels "github.com/nicolastakashi/community-perses-dashboards/pkg/panels/node_exporter"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
)

func withClusterCPU(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("CPU",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeCPUUsagePercentage(datasource, labelMatcher),
		panels.ClusterNodeCPUSaturationPercentage(datasource, labelMatcher),
	)
}

func withClusterMemory(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeMemoryUsagePercentage(datasource, labelMatcher),
		panels.ClusterNodeMemorySaturationPercentage(datasource, labelMatcher),
	)
}

func withClusterNetwork(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeNetworkUsageBytes(datasource, labelMatcher),
		panels.ClusterNodeNetworkSaturationBytes(datasource, labelMatcher),
	)
}

func withClusterDiskIO(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk IO",
		panelgroup.PanelsPerLine(2),
		panels.ClusterNodeDiskUsagePercentage(datasource, labelMatcher),
		panels.ClusterNodeDiskSaturationPercentage(datasource, labelMatcher),
	)
}

func withClusterDiskSpace(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Disk Space",
		panelgroup.PanelsPerLine(1),
		panels.ClusterNodeDiskSpacePercentage(datasource, labelMatcher),
	)
}

func BuildNodeExporterClusterUseMethod(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("node-exporter-cluster-use-method",
		dashboard.ProjectName(project),
		dashboard.Name("Node Exporter / USE Method / Cluster"),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "node_uname_info{job='node', sysname!='Darwin'}"),
		withClusterCPU(datasource, clusterLabelMatcher),
		withClusterMemory(datasource, clusterLabelMatcher),
		withClusterNetwork(datasource, clusterLabelMatcher),
		withClusterDiskIO(datasource, clusterLabelMatcher),
		withClusterDiskSpace(datasource, clusterLabelMatcher),
	)
}
