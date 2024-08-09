package metrics

import (
	"encoding/base64"
	"log/slog"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stellar/go/protocols/horizon"
)

func WrapDebug(name string, fn cron.FuncJob) cron.FuncJob {
	return func() {
		now := time.Now()
		slog.Debug("start", slog.Any("name", name))
		defer func() {
			slog.Debug("finish",
				slog.Any("name", name),
				slog.Any("duration", time.Since(now)))
		}()
		fn()
	}
}

func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func base64Decode(s string) string {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func isDataEqual(acc horizon.Account, key, value string) bool {
	b, _ := acc.GetData(key)
	return strings.ToLower(string(b)) == value
}

func isDataExist(acc horizon.Account, key string) bool {
	b, err := acc.GetData(key)
	return err == nil && len(b) > 0
}
