package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()

		logFields := logrus.Fields{
			"method":    c.Request().Method,
			"path":      c.Request().URL.Path,
			"status":    c.Response().Status,
			"latency":   stop.Sub(start),
			"client_id": c.RealIP(),
		}

		switch {
		case c.Response().Status >= 500:
			logrus.WithFields(logFields).Error("handled request with an error")

		case c.Response().Status >= 400:
			logrus.WithFields(logFields).Warn("handled request with a warn")

		case c.Response().Status >= 300:
			logrus.WithFields(logFields).Info("handled request with an info")

		default:
			logrus.WithFields(logFields).Debug("handled request")

		}
		return err
	}
}
