package main

import (
	"flag"
	"fmt"

	dashboards "github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards/prometheus"
	"github.com/perses/perses/go-sdk/dashboard"
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

	writer := dashboards.NewExec()
	builders := []dashboard.Builder{}

	po, err := prometheus.BuildPrometheusOverview(project, datasource, clusterLabelName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	builders = append(builders, po)

	prw, err := prometheus.BuildPrometheusRemoteWrite(project, datasource, clusterLabelName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	builders = append(builders, prw)

	for _, builder := range builders {
		writer.BuildDashboard(builder, nil)
	}
}
