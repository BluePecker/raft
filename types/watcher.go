package types

type Watcher struct {
    Heartbeat  func(term, leaderId uint64, member []uint64)
    Canvassing func(followerId uint64, bill Bill) bool
}
