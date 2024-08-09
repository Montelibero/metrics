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
	isCouncilReady, isCouncilDelegation, isAssemblyDelegation, isFundDelegation, isBSNBasicFilled bool
	mtlapCount                                                                                    string
}

func (m *Metrics) MTLAPGaugeReset() {
	m.mtlapGauge.Reset()
	m.l.Debug("[metrics] mtlap_total reset")
}

func (m *Metrics) MTLAPGaugeInc(p MTLAPGaugeParams) {
	m.mtlapGauge.WithLabelValues(
		strconv.FormatBool(p.isCouncilReady),
		strconv.FormatBool(p.isCouncilDelegation),
		strconv.FormatBool(p.isAssemblyDelegation),
		strconv.FormatBool(p.isFundDelegation),
		strconv.FormatBool(p.isBSNBasicFilled),
		p.mtlapCount,
	).Inc()
	m.l.Debug("[metrics] mtlap_total inc", slog.Any("params", p))
}
