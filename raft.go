package raft

import (
    "github.com/BluePecker/raft/types"
    "github.com/BluePecker/snowflake"
)

type raft struct {
    Term, UniqueID uint64
    leaderID       uint64
    
    clock          *types.Clock
    
    member         *types.Nodes
    
    identity       *types.Identity
}

func NewRafter(NodeId uint64) (*raft, error) {
    IdWorker, err := snowflake.NewIdWorker(NodeId)
    if err != nil {
        return nil, err
    }
    
    UniqueId, err := IdWorker.NextId()
    if err != nil {
        return nil, err
    }
    
    return &raft{
        Term: 0x0,
        leaderID: 0x0,
        UniqueID: UniqueId,
    }, nil
}