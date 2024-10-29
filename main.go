package main

import (
	"flag"

	dashboards "github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	nodeexporter "github.com/nicolastakashi/community-perses-dashboards/internal/dashboards/node_exporter"
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards/prometheus"
)

var (
	project          string
	datasource       string
	clusterLabelName string
)

func main() {

	flag.StringVar(&project, "project", "default", "The project name")
	flag.StringVar(&datasource, "datasource", "", "The datasource name")
	flag.StringVar(&clusterLabelName, "cluster-label-name", "", "The cluster label name")
	flag.Parse()

	dashboardWriter := dashboards.NewDashboardWriter()

	dashboardWriter.Add(prometheus.BuildPrometheusOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(prometheus.BuildPrometheusRemoteWrite(project, datasource, clusterLabelName))
	dashboardWriter.Add(nodeexporter.BuildNodeExporterNodes(project, datasource, clusterLabelName))
	dashboardWriter.Add(nodeexporter.BuildNodeExporterClusterUseMethod(project, datasource, clusterLabelName))

	dashboardWriter.Write()
}
