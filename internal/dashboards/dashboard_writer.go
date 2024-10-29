package dashboards

import "github.com/perses/perses/go-sdk/dashboard"

type DashboardWriter struct {
	dashboardResults []DashboardResult
	executor         Exec
}

type DashboardResult struct {
	builder dashboard.Builder
	err     error
}

func NewDashboardWriter() *DashboardWriter {
	return &DashboardWriter{
		executor: NewExec(),
	}
}

func (w *DashboardWriter) Add(builder dashboard.Builder, err error) {
	w.dashboardResults = append(w.dashboardResults, DashboardResult{
		builder: builder,
		err:     err,
	})
}

func (w *DashboardWriter) Write() {
	for _, result := range w.dashboardResults {
		w.executor.BuildDashboard(result.builder, result.err)
	}
}
