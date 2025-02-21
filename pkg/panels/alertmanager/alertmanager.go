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
// The panel uses the following Prometheus metrics:
// - alertmanager_alerts: Current number of alerts stored in Alertmanager
//
// Parameters:
//   - datasourceName: The name of the Prometheus data source.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: A panel option that can be added to a panel group.
func Alerts(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alerts",
		panel.Description("Shows current alerts in Alertmanager"),
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
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Alertmanager - Alerts"),
			),
		),
	)
}

// AlertsReceiveRate creates a panel option for displaying the rate of alerts received by the Alertmanager.
// The panel uses the following Prometheus metrics:
// - alertmanager_alerts_received_total: Total number of alerts received
// - alertmanager_alerts_invalid_total: Total number of invalid alerts received
//
// The panel shows:
// - Rate of received alerts per instance
// - Rate of invalid alerts per instance
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the query.
//   - labelMatchers: Optional Prometheus label matchers to filter the query.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func AlertsReceiveRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Alerts receive rate",
		panel.Description("Shows alert receive rate in Alertmanager"),
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
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Alertmanager - Received"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_alerts_invalid_total{job=~'$job'}[5m])) by (job,instance)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Alertmanager - Invalid"),
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
// The panel uses the following Prometheus metrics:
// - alertmanager_notifications_total: Total count of notifications sent
// - alertmanager_notifications_failed_total: Total count of failed notification attempts
//
// The panel shows:
// - Rate of total notifications sent per integration
// - Rate of failed notifications per integration
//
// Parameters:
// - datasourceName: The name of the data source to be used for the queries.
// - labelMatchers: A variadic parameter for Prometheus label matchers to filter the queries.
//
// Returns:
// - panelgroup.Option: The configured panel option.
func NotificationsSendRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Notifications Send Rate",
		panel.Description("Shows notification send rate for the Alertmanager"),
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
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{integration}} - Total"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_notifications_failed_total{job=~'$job', integration=~'$integration'}[5m])) by (integration, instance)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{integration}} - Failed"),
			),
		),
	)
}

// NotificationDuration creates a panel option for displaying the notification duration metrics
// from Alertmanager. It generates a time series panel with queries for the 99th percentile,
// median, and average notification latency.
//
// The panel uses the following Prometheus metrics:
// - alertmanager_notification_latency_seconds_bucket: Histogram of notification latency
// - alertmanager_notification_latency_seconds_sum: Total sum of notification latency
// - alertmanager_notification_latency_seconds_count: Total count of notifications
//
// The panel shows:
// - 99th percentile of notification latency
// - Median notification latency
// - Average notification latency
//
// Parameters:
//   - datasourceName: The name of the data source to be used for the queries.
//   - labelMatchers: A variadic parameter for Prometheus label matchers to filter the metrics.
//
// Returns:
//   - panelgroup.Option: An option that adds the configured panel to a panel group.
func NotificationDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Notification Duration",
		panel.Description("Shows notification latency for the Alertmanager"),
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
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{integration}} - 99th Percentile"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(rate(alertmanager_notification_latency_seconds_bucket{job=~'$job', integration=~'$integration'}[5m])) by (le,integration,instance))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{integration}} - Median"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(alertmanager_notification_latency_seconds_sum{job=~'$job', integration=~'$integration'}[5m])) by (integration,instance) / sum(rate(alertmanager_notification_latency_seconds_count{job=~'$job', integration=~'$integration'}[5m])) by (integration,instance)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{integration}} - Average"),
			),
		),
	)
}
