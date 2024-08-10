package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Montelibero/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/urfave/cli/v2"
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}()

func main() {
	logLevelMap := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}

	app := &cli.App{
		Version: Commit,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Value:   ":9090",
				Usage:   "Address to listen on",
				EnvVars: []string{"MTL_METRICS_LISTEN"},
			},
			&cli.StringFlag{
				Name:    "mtlap-total-cron",
				Aliases: []string{"mtc"},
				Value:   "* */10 * * * *",
				Usage:   "Cron expression for updating mtlap_total",
				EnvVars: []string{"MTL_METRICS_MTLAP_TOTAL_CRON"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"log", "l"},
				Value:   "INFO",
				Usage:   "Log level",
				EnvVars: []string{"MTL_METRICS_LOG_LEVEL"},
				Action: func(_ *cli.Context, v string) error {
					_, ok := logLevelMap[strings.ToUpper(v)]
					if !ok {
						return fmt.Errorf("unknown log level %q", v)
					}
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevelMap[strings.ToUpper(c.String("log-level"))],
			}))
			cl := horizonclient.DefaultPublicNetClient

			l.Info("starting metrics", slog.String("commit", Commit))

			m := metrics.NewMetrics(l)
			mtlapGauge := metrics.NewMTLAPGauge(l, cl, m)

			cr := cron.New(
				cron.WithSeconds(),
				cron.WithLogger(metrics.NewCronLogger(l)),
			)

			l.Info("updating mtlap_total")

			metrics.WrapDebug("mtlapGauge", mtlapGauge.Update)()

			l.Info("adding cron job", slog.String("cron", c.String("mtlap-total-cron")))

			_, err := cr.AddFunc(c.String("mtlap-total-cron"),
				metrics.WrapDebug("mtlapGauge", mtlapGauge.Update))
			if err != nil {
				return err
			}

			cr.Start()
			defer cr.Stop()

			http.Handle("/metrics", promhttp.Handler())

			l.Info("starting server", slog.String("addr", c.String("addr")))

			if err := http.ListenAndServe(c.String("addr"), nil); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
