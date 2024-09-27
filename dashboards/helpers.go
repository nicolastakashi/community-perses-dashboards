package dashboards

import (
	"github.com/perses/perses/go-sdk/dashboard"
	"github.com/perses/perses/go-sdk/prometheus/query"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
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
			labelValuesVar.PrometheusLabelValues("cluster",
				labelValuesVar.Matchers(matcher),
				AddVariableDatasource(datasource),
			),
			listVar.DisplayName("cluster"),
		),
	)
}

func LabelsSetPromQL(query, labelMatchType, name, value string) string {
	if name == "" || value == "" {
		return query
	}

	expr, err := parser.ParseExpr(query)
	if err != nil {
		return ""
	}
	var matchType labels.MatchType
	switch labelMatchType {
	case parser.ItemType(parser.EQL).String():
		matchType = labels.MatchEqual
	case parser.ItemType(parser.NEQ).String():
		matchType = labels.MatchNotEqual
	case parser.ItemType(parser.EQL_REGEX).String():
		matchType = labels.MatchRegexp
	case parser.ItemType(parser.NEQ_REGEX).String():
		matchType = labels.MatchNotRegexp
	default:
		return ""
	}

	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		if n, ok := node.(*parser.VectorSelector); ok {
			var found bool
			for i, l := range n.LabelMatchers {
				if l.Name == name {
					n.LabelMatchers[i].Type = matchType
					n.LabelMatchers[i].Value = value
					found = true
				}
			}
			if !found {
				n.LabelMatchers = append(n.LabelMatchers, &labels.Matcher{
					Type:  matchType,
					Name:  name,
					Value: value,
				})
			}
		}
		return nil
	})
	return expr.Pretty(0)
}
