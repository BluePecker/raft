package raft

import (
    "time"
    "math/rand"
    "github.com/BluePecker/raft/types"
    "github.com/BluePecker/snowflake"
)

type raft struct {
    Term, UniqueId uint64
    leaderId       uint64
    
    refresh        chan struct{}
    
    bill           types.Bill
    
    clock          *types.Clock
    member         *types.Nodes
    identity       *types.Identity
    watcher        *types.Watcher
    ballotBox      *types.Nodes
}

func iToSec(I int) time.Duration {
    return time.Duration(I) * time.Second
}

func randWait(Millisecond int) {
    WaitSec := rand.New(rand.NewSource(time.Now().UnixNano()))
    time.Sleep(time.Duration(WaitSec.Intn(Millisecond)) * time.Millisecond)
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

func (r *raft) prepare() {
    var BackupTerm uint64 = r.Term
    for ; r.Term <= r.bill.Term; r.Term++ {}
    
    r.bill = types.Bill{
        Term: r.Term,
        UniqueId: r.UniqueId,
    }
    
    Sec := iToSec(r.clock.Second.Timeout)
    r.clock.Timer.Timeout = time.NewTimer(Sec)
    
    if r.ballotBox == nil {
        r.ballotBox = &types.Nodes{}
    }
    
    var Canvassing = func(UniqueId uint64) {
        if r.watcher.Canvassing == nil {
            return
        }
        var Bill types.Bill = types.Bill{
            Term: r.Term,
            UniqueId: r.UniqueId,
        }
        if r.watcher.Canvassing(UniqueId, Bill) {
            return
        }
        r.ballotBox.PushBack(r.Term, UniqueId)
    }
    
    if r.member == nil {
        r.member = &types.Nodes{}
    }
    
    go func() {
        Next := r.member.Front(BackupTerm)
        for ; Next != nil; Next = Next.Next() {
            UniqueId, ok := Next.Value.(uint64)
            if ok && UniqueId != r.UniqueId {
                go Canvassing(UniqueId)
            }
        }
    }()
}

func (r *raft) aggregate() bool {
    if r.ballotBox != nil && r.member != nil {
        Win := r.ballotBox.Len(r.Term) >= r.member.Len(r.Term) / 2 + 1
        return Win
    }
    return false
}

func (r *raft) Candidate() {
    again:r.prepare()
    for {
        select {
        case <-r.clock.Timer.Timeout.C:
            if r.aggregate() {
                randWait(1000)
                goto again
            }
            r.clock.Timer.Timeout.Stop()
            go r.identity.BecomeLeader()
            return
        case <-r.refresh:
            r.clock.Timer.Timeout.Stop()
            go r.identity.BecomeFollower()
            return
        }
    }
}

func (r *raft) Leader() {
    r.leaderId = r.UniqueId
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
            Next := r.member.Front(r.Term)
            for ; Next != nil; Next = Next.Next() {
                if NodeId, ok := Next.Value.(uint64); ok {
                    Nodes = append(Nodes, NodeId)
                }
            }
            go func(term, leaderId uint64, member []uint64) {
                if r.watcher == nil || r.watcher.Heartbeat == nil {
                    return
                }
                r.watcher.Heartbeat(term, leaderId, member)
            }(r.Term, r.leaderId, Nodes)
        }
    }
}

func (r *raft) Vote(Bill types.Bill) bool {
    if Bill.Term <= r.bill.Term {
        return false
    }
    r.refresh <- struct{}{}
    r.bill = Bill
    return true
}

func (r *raft) Sync(LeaderId, Term uint64, Members []uint64) {
    if Term >= r.Term {
        r.refresh <- struct{}{}
        
        r.Term = Term
        r.leaderId = LeaderId
        
        for _, UniqueId := range Members {
            r.member.PushBack(Term, UniqueId)
        }
    }
}

// todo
func (r *raft) NightWatch() {
    
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
        UniqueId: uint64(UniqueId),
        Term: 0x0,
        leaderId: 0x0,
        refresh: make(chan struct{}),
        clock: &types.Clock{
            Second: Times,
        },
        member: &types.Nodes{},
        ballotBox: &types.Nodes{},
    }
    
    rafter.identity = (&types.Identity{}).New()
    
    return rafter, nil
}