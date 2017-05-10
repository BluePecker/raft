package types

import (
    "errors"
)

const (
    // 选民
    FOLLOWER int8 = 1 << iota
    // 候选人
    CANDIDATE
    // 领导
    LEADER
)

type Identity struct {
    role chan int8
}

func (i *Identity) init() {
    if i.role == nil {
        i.role = make(chan int8, 1)
    }
}

func (i *Identity) NightWatch(standby, campaign, order func()) {
    i.init()
    
    if len(i.role) < 1 {
        i.role <- FOLLOWER
    }
    
    for {
        switch <-i.role {
        case FOLLOWER:
            go standby()
        case CANDIDATE:
            go campaign()
        case LEADER:
            go order()
        default:
            panic(errors.New("undefined identity"))
        }
    }
}

func (i *Identity) BecomeFollower() {
    i.init()
    i.role <- FOLLOWER
}

func (i *Identity) BecomeCandidate() {
    i.init()
    i.role <- CANDIDATE
}

func (i *Identity) BecomeLeader() {
    i.init()
    i.role <- LEADER
}

func (i *Identity) New() *Identity {
    return &Identity{
        role: make(chan int8, 1),
    }
}