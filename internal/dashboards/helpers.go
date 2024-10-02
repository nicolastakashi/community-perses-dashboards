package dashboards

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	"github.com/perses/perses/go-sdk/prometheus/query"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func AddVariableDatasource(datasourceName string) labelValuesVar.Option {
	if datasourceName == "" {
		return func(plugin *labelValuesVar.Builder) error {
			return nil
		}
	}
	return labelValuesVar.Datasource(datasourceName)
}

func AddQueryDataSource(datasourceName string) query.Option {
	if datasourceName == "" {
		return func(plugin *query.Builder) error {
			return nil
		}
	}
	return query.Datasource(datasourceName)
}

func AddClusterVariable(datasource, clusterLabelName, matcher string) dashboard.Option {
	if clusterLabelName == "" {
		return func(builder *dashboard.Builder) error {
			return nil
		}
	}
	return dashboard.AddVariable("cluster",
		listVar.List(
			labelValuesVar.PrometheusLabelValues(clusterLabelName,
				labelValuesVar.Matchers(matcher),
				AddVariableDatasource(datasource),
			),
			listVar.DisplayName(clusterLabelName),
		),
	)
}

func GetClusterLabelMatcher(clusterLabelName string) promql.LabelMatcher {
	return promql.LabelMatcher{
		Name:  clusterLabelName,
		Value: "$cluster",
		Type:  "=",
	}
}
