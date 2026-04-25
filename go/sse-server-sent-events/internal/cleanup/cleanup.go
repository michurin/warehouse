package cleanup

import (
	"log/slog"
	"time"

	"github.com/michurin/minchat/internal/xdto"
	"github.com/michurin/minchat/internal/xhouse"
)

func RevisionLoop(ch *xhouse.House, inactiveTime time.Duration) {
	for {
		ms := time.Now().Add(-inactiveTime).UnixMilli()
		walls, users := ch.Audit(ms)
		for i, w := range walls {
			slog.Info("Run: notify")
			w.Pub(xdto.BuildResponse(xdto.BuildRobotMessage(ms, "Someone got out"), users[i], false))
		}
		time.Sleep(2 * time.Second)
	}
}
