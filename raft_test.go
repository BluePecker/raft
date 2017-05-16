package raft

import (
    "testing"
    "github.com/BluePecker/raft/types"
    "time"
    "fmt"
)

func TestNewRafter(t *testing.T) {
    Second := types.Second{
        Tenure: 6,
        Timeout: 2,
        Heartbeat: 1,
    }
    
    rafter1, err := NewRafter(1, Second, []uint64{1, 2})
    if err != nil {
        t.Fatal(1, err)
    }
    rafter2, err := NewRafter(2, Second, []uint64{1, 2})
    if err != nil {
        t.Fatal(2, err)
    }
    
    Container := map[uint64]*raft{
        1: rafter1,
        2: rafter2,
    }
    
    Watcher := &types.Watcher{
        Heartbeat: func(Term, LeaderId uint64, Members []uint64) {
            for _, UniqueId := range Members {
                Container[UniqueId].Sync(LeaderId, Term, Members)
            }
        },
        Canvassing: func(FollowerId uint64, Bill types.Bill) bool {
            return Container[FollowerId].Vote(Bill)
        },
    }
    
    rafter1.NightWatch(Watcher)
    rafter2.NightWatch(Watcher)
    
    rafter1.Start()
    rafter2.Start()
    
    for {
        time.Sleep(time.Duration(1) * time.Second)
        fmt.Println(rafter1.Term, "->", rafter1.identity.Show, " | ", rafter2.Term, "->", rafter2.identity.Show)
    }
}
