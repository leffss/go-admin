package zap

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func Ginzap(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return Logger(logger, WithTimeFormat(timeFormat), WithUTC(utc))
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return Recovery(logger, stack)
}

// Option logger/recover option
type Option func(c *Config)

// WithTimeFormat optional a time package format string (e.g. time.RFC3339).(default time.RFC3339Nano).
func WithTimeFormat(layout string) Option {
	return func(c *Config) {
		c.timeFormat = layout
	}
}

// WithUTC a boolean stating whether to use UTC time zone or local.(default local).
func WithUTC(b bool) Option {
	return func(c *Config) {
		c.utc = b
	}
}

func WithBody(b bool) Option {
	return func(c *Config) {
		c.withBody = b
	}
}

// WithCustomFields optional custom field
func WithCustomFields(fields ...func(c *gin.Context) zap.Field) Option {
	return func(c *Config) {
		c.customFields = fields
	}
}

// Config logger/recover config
type Config struct {
	timeFormat   string
	utc          bool
	withBody     bool
	customFields []func(c *gin.Context) zap.Field
}

// Logger returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// Default option:
//   1. A time package format string (e.g. time.RFC3339).(default time.RFC3339Nano)
//   2. Use time zone.(e.g. utc time zone).(default local).
//   3. Custom fields.(default nil)
func Logger(logger *zap.Logger, opts ...Option) gin.HandlerFunc {
	cfg := Config{
		time.RFC3339Nano,
		false,
		false,
		nil,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		var body []byte
		if cfg.withBody {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if cfg.utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			fields := make([]zap.Field, 0, 8 + len(cfg.customFields))
			fields = append(fields,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(cfg.timeFormat)),
				zap.Duration("latency", latency),
			)
			if cfg.withBody && string(body) != "" {
				fields = append(fields, zap.String("body", string(body)))
			}
			for _, field := range cfg.customFields {
				fields = append(fields, field(c))
			}
			logger.Info(path, fields...)
		}
	}
}

// Recovery returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func Recovery(logger *zap.Logger, stack bool, opts ...Option) gin.HandlerFunc {
	cfg := Config{
		time.RFC3339Nano,
		false,
		false,
		nil,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if stack {
		cfg.customFields = append(cfg.customFields, func(c *gin.Context) zap.Field {
			return zap.ByteString("stack", debug.Stack())
		})
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.ByteString("request", httpRequest),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				now := time.Now()
				if cfg.utc {
					now = now.UTC()
				}
				fields := make([]zap.Field, 0, 3+len(cfg.customFields))
				fields = append(fields,
					zap.String("time", now.Format(cfg.timeFormat)),
					zap.Any("error", err),
					zap.ByteString("request", httpRequest),
				)
				for _, field := range cfg.customFields {
					fields = append(fields, field(c))
				}
				logger.Error("[Recovery from panic]", fields...)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
