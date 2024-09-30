package main

import (
	"flag"
	"fmt"

	"github.com/nicolastakashi/community-perses-dashboards/dashboards/prometheus"
	dashboards "github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
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

	rw, err := prometheus.BuildPrometheusRemoteWrite(project, datasource, clusterLabelName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	builders = append(builders, rw)

	for _, builder := range builders {
		writer.BuildDashboard(builder, nil)
	}
}
