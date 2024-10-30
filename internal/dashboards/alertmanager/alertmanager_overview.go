package alertmanager

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	panels "github.com/nicolastakashi/community-perses-dashboards/pkg/panels/alertmanager"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withAlertsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alerts",
		panelgroup.PanelsPerLine(2),
		panels.Alerts(datasource, labelMatcher),
		panels.AlertsReceiveRate(datasource, labelMatcher),
	)
}

func withNotificationsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Notifications",
		panelgroup.PanelsPerLine(2),
		panels.NotificationsSendRate(datasource, labelMatcher),
		panels.NotificationDuration(datasource, labelMatcher),
	)
}

func BuildAlertManagerOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("alertmanager-overview",
		dashboard.ProjectName(project),
		dashboard.Name("Alertmanager / Overview"),
		dashboard.AddVariable("job",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("job",
					labelValuesVar.Matchers("alertmanager_alerts"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("job"),
			),
		),
		dashboards.AddClusterVariable(datasource, clusterLabelName, "alertmanager_alerts"),
		dashboard.AddVariable("integration",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("integration",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"alertmanager_notifications_total",
							[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.AllowAllValue(true),
				listVar.AllowMultiple(true),
				listVar.DisplayName("integration"),
			),
		),
		withAlertsGroup(datasource, clusterLabelMatcher),
		withNotificationsGroup(datasource, clusterLabelMatcher),
	)
}
