package telemetry

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icinga/icingadb/pkg/icingaredis"
	"github.com/icinga/icingadb/pkg/logging"
	"github.com/icinga/icingadb/pkg/utils"
	"go.uber.org/zap"
	"sync"
)

// Err reports the current error on $subject on demand to Redis for being monitored by Icinga 2.
type Err struct {
	ctx         context.Context
	redis       *icingaredis.Client
	subject     string
	key         string
	logger      *logging.Logger
	refurbisher *utils.Refurbisher

	mtx      sync.Mutex
	lastErr  string
	lastNil  bool
	everSent bool
}

// NewErr creates new Err.
func NewErr(ctx context.Context, redis *icingaredis.Client, subject string, logger *logging.Logger) *Err {
	return &Err{
		ctx:         ctx,
		redis:       redis,
		subject:     subject,
		key:         "icingadb:telemetry:err:" + subject,
		logger:      logger,
		refurbisher: utils.NewRefurbisher(),
	}
}

// Report does the actual reporting.
func (e *Err) Report(err error) {
	var currentErr string
	currentNil := true

	if err != nil {
		currentErr = err.Error()
		currentNil = false
	}

	e.mtx.Lock()

	if e.everSent && currentNil == e.lastNil && currentErr == e.lastErr {
		e.mtx.Unlock()
		return
	}

	e.lastErr = currentErr
	e.lastNil = currentNil
	e.everSent = true

	e.mtx.Unlock()

	_ = e.refurbisher.Do(e.ctx, func(ctx context.Context) error {
		var cmd redis.Cmder
		if currentNil {
			cmd = e.redis.Del(ctx, e.key)
		} else {
			cmd = e.redis.Set(ctx, e.key, currentErr, 0)
		}

		if err := cmd.Err(); err != nil && !utils.IsContextCanceled(err) {
			e.logger.Warnw(
				"Can't update last own error", zap.String("subject", e.subject), zap.Error(icingaredis.WrapCmdErr(cmd)),
			)
		}

		return nil
	})
}
