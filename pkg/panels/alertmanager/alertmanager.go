package alertmanager

import (
	"github.com/nicolastakashi/community-perses-dashboards/internal/dashboards"
	"github.com/nicolastakashi/community-perses-dashboards/internal/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
	"github.com/perses/perses/go-sdk/prometheus/query"
)

// Alerts creates a panel option for displaying the count of alerts
// from Alertmanager using Prometheus as the data source. It generates a time series panel with
// specific configurations for the Y-axis format and legend position.
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func Alerts(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alerts",
		panel.Description("Current set of alerts stored in the Alertmanager"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(alertmanager_alerts{job=~'$job'}) by (instance)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

// AlertsReceiveRate creates a panel option for displaying the rate of alerts received by the Alertmanager.
// It includes a description of the panel, a time series chart with a legend, and a PromQL query to fetch the data.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMathers: Optional Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func AlertsReceiveRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alerts receive rate",
		panel.Description("Rate of successful and invalid alerts received by the Alertmanager"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_alerts_received_total{job=~'$job'}[5m])) by (job,instance)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} Received"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_alerts_invalid_total{job=~'$job'}[5m])) by (job,instance)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} Invalid"),
			),
		),
	)
}

// NotificationsSendRate creates a panel option for displaying the rate of successful
// and invalid notifications sent by the Alertmanager. It generates a time series
// panel with a legend positioned at the bottom in table mode, showing the last
// calculation value. The panel includes two PromQL queries: one for the total
// notifications sent and another for the failed notifications, both grouped by
// integration and instance.
//
// Parameters:
// - datasourceName: The name of the data source to be used for the queries.
// - labelMathers: A variadic parameter for Prometheus label matchers to filter the queries.
//
// Returns:
// - panelgroup.Option: The configured panel option.
func NotificationsSendRate(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Notifications Send Rate",
		panel.Description("Rate of successful and invalid notifications sent by the Alertmanager"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_notifications_total{job=~'$job', integration=~'$integration'}[5m])) by (integration, instance)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ integration }} - {{instance}} Total"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_notifications_failed_total{job=~'$job', integration=~'$integration'}[5m])) by (integration, instance)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ integration }} - {{instance}} Total"),
			),
		),
	)
}

// NotificationDuration creates a panel option for displaying the notification duration metrics
// from Alertmanager. It generates a time series panel with queries for the 99th percentile,
// median, and average notification latency.
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the queries.
//   - labelMathers: A variadic parameter for Prometheus label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: An option that adds the configured panel to a panel group.
func NotificationDuration(datasourceName string, labelMathers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Notification Duration",
		panel.Description("Latency of notifications sent by the Alertmanager"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.SecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(alertmanager_notification_latency_seconds_bucket{job=~'$job', integration=~'$integration'}[5m])) by (le,integration,instance))",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ integration }} - {{instance}} 99th "),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(rate(alertmanager_notification_latency_seconds_bucket{job=~'$job', integration=~'$integration'}[5m])) by (le,integration,instance))",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ integration }} - {{instance}} Median"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_notification_latency_seconds_sum{job=~'$job', integration=~'$integration'}[5m])) by (integration,instance) / sum(rate(alertmanager_notification_latency_seconds_count{job=~'$job', integration=~'$integration'}[5m])) by (integration,instance)",
					labelMathers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ integration }} - {{instance}} Average"),
			),
		),
	)
}
