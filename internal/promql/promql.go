package promql

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

type LabelMatcher struct {
	Name  string
	Value string
	Type  string
}

func SetLabelMatchers(query string, labelMathers []LabelMatcher) string {
	for _, l := range labelMathers {
		query = LabelsSetPromQL(query, l.Type, l.Name, l.Value)
	}
	return query
}

func LabelsSetPromQL(query, labelMatchType, name, value string) string {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return ""
	}

	if name == "" || value == "" {
		return expr.Pretty(0)
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
