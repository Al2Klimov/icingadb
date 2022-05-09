package telemetry

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icinga/icingadb/pkg/com"
	"github.com/icinga/icingadb/pkg/icingaredis"
	"github.com/icinga/icingadb/pkg/logging"
	"github.com/icinga/icingadb/pkg/periodic"
	"github.com/icinga/icingadb/pkg/utils"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// Stats periodically forward reported metrics to Redis for being monitored by Icinga 2.
type Stats struct {
	// Config & co. are to be increased by the T sync once for every T object synced.
	Config, State, History, Overdue, HistoryCleanup com.Counter

	ctx    context.Context
	redis  *icingaredis.Client
	logger *logging.Logger
}

// NewStats creates new Stats.
func NewStats(ctx context.Context, redis *icingaredis.Client, logger *logging.Logger) *Stats {
	s := &Stats{
		ctx:    ctx,
		redis:  redis,
		logger: logger,
	}

	s.start()

	return s
}

// start kicks off the actual work.
func (s *Stats) start() {
	counters := map[string]*com.Counter{
		"sync_config":     &s.Config,
		"sync_state":      &s.State,
		"sync_history":    &s.History,
		"sync_overdue":    &s.Overdue,
		"cleanup_history": &s.HistoryCleanup,
	}

	periodic.Start(s.ctx, 3*time.Second, func(_ periodic.Tick) {
		var data []string
		for kind, counter := range counters {
			if cnt := counter.Reset(); cnt > 0 {
				data = append(data, kind, strconv.FormatUint(cnt, 10))
			}
		}

		if data != nil {
			cmd := s.redis.XAdd(s.ctx, &redis.XAddArgs{
				Stream: "icingadb:telemetry:stats",
				MaxLen: 5 * 60,
				Approx: true,
				Values: data,
			})
			if err := cmd.Err(); err != nil && !utils.IsContextCanceled(err) {
				s.logger.Warnw("Can't update own stats", zap.Error(icingaredis.WrapCmdErr(cmd)))
			}
		}
	})
}
