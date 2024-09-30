package main

import (
	"flag"

	"github.com/nicolastakashi/community-perses-dashboards/dashboards/prometheus"
	"github.com/perses/perses/go-sdk"
)

var (
	project          string
	datasource       string
	clusterLabelName string
)

func main() {

	flag.StringVar(&project, "project", "default", "The project name")
	flag.StringVar(&datasource, "datasource", "", "The datasource name")
	flag.StringVar(&clusterLabelName, "cluster-label-name", "cluster", "The cluster label name")
	flag.Parse()
	exec := sdk.NewExec()

	// prometheus.BuildPrometheusOverview(exec, project, datasource, clusterLabelName)
	prometheus.BuildPrometheusRemoteWrite(exec, project, datasource, clusterLabelName)
}
