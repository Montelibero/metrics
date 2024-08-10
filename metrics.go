package metrics

import (
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	l          Logger
	mtlapGauge *prometheus.GaugeVec
}

func NewMetrics(l Logger) *Metrics {
	mtlacGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mtlap_total",
			Help: "Total MTLAP accounts",
		},
		[]string{
			"is_council_ready",
			"is_council_delegation",
			"is_assembly_delegation",
			"is_fund_delegation",
			"is_bsn_basic_filled",
			"is_bsn_partialy_filled",
			"mtlap_count",
		},
	)

	prometheus.MustRegister(mtlacGauge)

	return &Metrics{
		l:          l,
		mtlapGauge: mtlacGauge,
	}
}

type MTLAPGaugeParams struct {
	MTLAPCount           string
	IsCouncilReady       bool
	IsCouncilDelegation  bool
	IsAssemblyDelegation bool
	IsFundDelegation     bool
	IsBSNBasicFilled     bool
	IsBSNPartialyFilled  bool
}

func (m *Metrics) MTLAPGaugeReset() {
	m.mtlapGauge.Reset()
	m.l.Debug("[metrics] mtlap_total reset")
}

func (m *Metrics) MTLAPGaugeInc(p MTLAPGaugeParams) {
	m.mtlapGauge.WithLabelValues(
		strconv.FormatBool(p.IsCouncilReady),
		strconv.FormatBool(p.IsCouncilDelegation),
		strconv.FormatBool(p.IsAssemblyDelegation),
		strconv.FormatBool(p.IsFundDelegation),
		strconv.FormatBool(p.IsBSNBasicFilled),
		strconv.FormatBool(p.IsBSNPartialyFilled),
		p.MTLAPCount,
	).Inc()
	m.l.Debug("[metrics] mtlap_total inc", slog.Any("params", p))
}
