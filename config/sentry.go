package config

import (
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
)

func SentryConfig() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config().Sentry.Dsn, // уникальный идентификатор проекта в Sentry
		EnableTracing:    true,                // мониторинга производительности
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
		os.Exit(1)
	}
}
