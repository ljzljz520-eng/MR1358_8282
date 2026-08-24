package operations

import "fmt"

type RouteMetric struct {
	RouteID       string
	Views         int
	Submissions   int
	Confirmations int
	Cancellations int
}

func (m RouteMetric) ConversionRate() float64 {
	if m.Views == 0 {
		return 0
	}
	return float64(m.Confirmations) / float64(m.Views)
}

func (m RouteMetric) ConfirmationRate() float64 {
	if m.Submissions == 0 {
		return 0
	}
	return float64(m.Confirmations) / float64(m.Submissions)
}

func (m RouteMetric) NetBookings() int {
	net := m.Confirmations - m.Cancellations
	if net < 0 {
		return 0
	}
	return net
}

type Dashboard struct{ Routes []RouteMetric }

func NewDashboard(metrics []RouteMetric) Dashboard {
	return Dashboard{Routes: append([]RouteMetric(nil), metrics...)}
}

func (d Dashboard) TotalViews() int {
	total := 0
	for _, metric := range d.Routes {
		total += metric.Views
	}
	return total
}

func (d Dashboard) TotalConfirmed() int {
	total := 0
	for _, metric := range d.Routes {
		total += metric.Confirmations
	}
	return total
}

func (d Dashboard) BestRoute() (RouteMetric, bool) {
	if len(d.Routes) == 0 {
		return RouteMetric{}, false
	}
	best := d.Routes[0]
	for _, metric := range d.Routes[1:] {
		if metric.ConfirmationRate() > best.ConfirmationRate() {
			best = metric
		}
	}
	return best, true
}

func (d Dashboard) Summary() string {
	best, ok := d.BestRoute()
	if !ok {
		return "no route activity"
	}
	return fmt.Sprintf("views=%d confirmed=%d best=%s", d.TotalViews(), d.TotalConfirmed(), best.RouteID)
}
