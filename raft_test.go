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
    
    rafter1, err := NewRafter(1, Second)
    if err != nil {
        t.Fatal(1, err)
    }
    rafter2, err := NewRafter(2, Second)
    if err != nil {
        t.Fatal(2, err)
    }
    
    Watcher1 := types.Watcher{
        Heartbeat: func() {},
        Canvassing: func() {},
    }
    Watcher2 := types.Watcher{
        Heartbeat: func() {},
        Canvassing: func() {},
    }
    
    rafter1.NightWatch(Watcher1)
    rafter2.NightWatch(Watcher2)
    
    rafter1.Start()
    rafter2.Start()
    
    for {
        time.Sleep(time.Duration(1) * time.Second)
        fmt.Println(rafter1.identity.Show, rafter2.identity.Show)
    }
}
