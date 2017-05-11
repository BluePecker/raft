package types

type Watcher struct {
    Heartbeat func(term, leaderId uint64, member []uint64)
}
