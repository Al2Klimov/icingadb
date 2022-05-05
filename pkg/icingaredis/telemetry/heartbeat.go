package telemetry

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icinga/icingadb/internal"
	"github.com/icinga/icingadb/pkg/icingaredis"
	"github.com/icinga/icingadb/pkg/logging"
	"github.com/icinga/icingadb/pkg/periodic"
	"github.com/icinga/icingadb/pkg/utils"
	"go.uber.org/zap"
	"time"
)

// StartHeartbeat periodically writes heartbeats to Redis for being monitored by Icinga 2.
func StartHeartbeat(ctx context.Context, client *icingaredis.Client, logger *logging.Logger) {
	xaa := &redis.XAddArgs{
		Stream: "icingadb:telemetry:heartbeat",
		MaxLen: 1,
		Values: []string{"version", internal.Version},
	}

	periodic.Start(ctx, time.Second, func(_ periodic.Tick) {
		cmd := client.XAdd(ctx, xaa)
		if err := cmd.Err(); err != nil && !utils.IsContextCanceled(err) {
			logger.Warnw("Can't update own heartbeat", zap.Error(icingaredis.WrapCmdErr(cmd)))
		}
	})
}
