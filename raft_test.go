package raft

import (
    "testing"
    "github.com/BluePecker/raft/types"
)

func TestNewRafter(t *testing.T) {
    rafter, _ := NewRafter(1, types.Second{
        Tenure: 10,
        Timeout: 2,
        Heartbeat: 1,
    })
    
    rafter.Start()
}
