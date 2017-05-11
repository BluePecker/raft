package raft

import (
    "time"
    "github.com/BluePecker/raft/types"
    "github.com/BluePecker/snowflake"
)

type raft struct {
    Term, UniqueID uint64
    leaderID       uint64
    
    refresh        chan struct{}
    
    clock          *types.Clock
    member         *types.Nodes
    identity       *types.Identity
    watcher        *types.Watcher
}

func iToSec(i int) time.Duration {
    return time.Duration(i) * time.Second
}

func (r *raft) Start() {
    // todo
    go r.identity.NightWatch(r)
}

func (r *raft) Follower() {
    Sec := iToSec(r.clock.Second.Tenure)
    r.clock.Timer.Tenure = time.NewTimer(Sec)
    for {
        select {
        case <-r.clock.Timer.Tenure.C:
            go r.identity.BecomeCandidate()
            return
        case <-r.refresh:
            r.clock.Timer.Tenure.Reset(Sec)
        }
    }
}

func (r *raft) Candidate() {
    
}

func (r *raft) Leader() {
    r.leaderID = r.UniqueID
    Sec := iToSec(r.clock.Second.Heartbeat)
    r.clock.Ticker.Heartbeat = time.NewTicker(Sec)
    for {
        select {
        case <-r.refresh:
            r.clock.Ticker.Heartbeat.Stop()
            r.identity.BecomeFollower()
            return
        case <-r.clock.Ticker.Heartbeat.C:
            var Nodes []uint64;
            for Next := r.member.Front(r.Term); Next != nil; Next = Next.Next() {
                if NodeId, ok := Next.Value.(uint64); ok {
                    Nodes = append(Nodes, NodeId)
                }
            }
            
            go func(term, leaderId uint64, member []uint64) {
                if r.watcher.Heartbeat == nil {
                    return
                }
                r.watcher.Heartbeat(term, leaderId, member)
            }(r.Term, r.leaderID, Nodes)
        }
    }
}

func NewRafter(NodeId int64, Times types.Second) (*raft, error) {
    if err := Times.Validator(); err != nil {
        return nil, err
    }
    
    IdWorker, err := snowflake.NewIdWorker(NodeId)
    if err != nil {
        return nil, err
    }
    
    UniqueId, err := IdWorker.NextId()
    if err != nil {
        return nil, err
    }
    
    rafter := &raft{
        UniqueID: uint64(UniqueId),
        Term: 0x0,
        leaderID: 0x0,
        
        refresh: make(chan struct{}),
        
        clock: &types.Clock{
            Second: Times,
        },
    }
    
    rafter.identity = (&types.Identity{}).New()
    
    return rafter, nil
}