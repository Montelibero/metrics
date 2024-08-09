package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Montelibero/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/stellar/go/clients/horizonclient"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	cl := horizonclient.DefaultPublicNetClient

	m := metrics.NewMetrics(l)
	mtlapGauge := metrics.NewMTLAPGauge(l, cl, m)

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(metrics.NewCronLogger(l)),
	)

	metrics.WrapDebug("mtlapGauge", mtlapGauge.Update)()

	_, err := c.AddFunc("@hourly",
		metrics.WrapDebug("mtlapGauge", mtlapGauge.Update))
	if err != nil {
		l.Error("failed to add func", slog.Any("error", err))
		os.Exit(1)
	}

	c.Start()
	defer c.Stop()

	http.Handle("/metrics", promhttp.Handler())

	l.Info("starting server on :80")

	if err := http.ListenAndServe(":80", nil); err != nil {
		l.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
