package raft

import (
    "testing"
    "github.com/BluePecker/raft/types"
    "time"
    "fmt"
)

func TestNewRafter(t *testing.T) {
    rafter, err := NewRafter(1, types.Second{
        Tenure: 6,
        Timeout: 2,
        Heartbeat: 1,
    })
    
    if err != nil {
        t.Fatal(err)
    }
    rafter.Start()
    for {
        time.Sleep(time.Duration(1) * time.Second)
        fmt.Println(rafter.identity.Show)
    }
}
