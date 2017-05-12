package raft

import (
    "testing"
    "github.com/BluePecker/raft/types"
    "time"
    "fmt"
)

func TestNewRafter(t *testing.T) {
    rafter, _ := NewRafter(1, types.Second{
        Tenure: 8,
        Timeout: 2,
        Heartbeat: 1,
    })
    
    rafter.Start()
    
    for {
        time.Sleep(time.Duration(1) * time.Second)
        fmt.Println(rafter.identity.Show)
    }
}
