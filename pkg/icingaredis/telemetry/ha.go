package telemetry

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icinga/icingadb/pkg/icingaredis"
	"github.com/icinga/icingadb/pkg/logging"
	"github.com/icinga/icingadb/pkg/utils"
	"go.uber.org/zap"
)

var boolToIntStr = map[bool]string{false: "0", true: "1"}

// ReportResponsibility updates whether and since when we're responsible (HA) in Redis for being monitored by Icinga 2.
func ReportResponsibility(ctx context.Context, client *icingaredis.Client, logger *logging.Logger, isResponsible bool) {
	cmd := client.XAdd(ctx, &redis.XAddArgs{
		Stream: "icingadb:telemetry:ha",
		MaxLen: 1,
		Values: []string{"is-responsible", boolToIntStr[isResponsible]},
	})
	if err := cmd.Err(); err != nil && !utils.IsContextCanceled(err) {
		logger.Warnw("Can't update own responsibility", zap.Error(icingaredis.WrapCmdErr(cmd)))
	}
}
