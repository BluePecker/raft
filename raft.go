package raft

import (
    "github.com/BluePecker/raft/types"
)

type raft struct {
    Term, UniqueID uint64
    leaderID       uint64
    
    clock          *types.Clock
    
    member         *types.Nodes

    identity       *types.Identity
}