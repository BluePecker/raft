package types

import "context"

type Watcher struct {
    Heartbeat func(cancel context.CancelFunc, term, leaderId uint64, member []uint64)
}
