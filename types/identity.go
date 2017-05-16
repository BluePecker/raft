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
    Show int8
}

type politician interface {
    Follower()
    Candidate()
    Leader()
}

func (i *Identity) init() {
    if i.role == nil {
        i.role = make(chan int8, 1)
    }
    if i.Show < FOLLOWER {
        i.Show = FOLLOWER
    }
}

func (i *Identity) NightWatch(pol politician) {
    i.init()
    
    for {
        switch <-i.role {
        case FOLLOWER:
            go pol.Follower()
        case CANDIDATE:
            go pol.Candidate()
        case LEADER:
            go pol.Leader()
        default:
            panic(errors.New("undefined identity"))
        }
    }
}

func (i *Identity) BecomeFollower() {
    i.init()
    i.Show = FOLLOWER
    i.role <- FOLLOWER
}

func (i *Identity) BecomeCandidate() {
    i.init()
    i.Show = CANDIDATE
    i.role <- CANDIDATE
}

func (i *Identity) BecomeLeader() {
    i.init()
    i.Show = LEADER
    i.role <- LEADER
}

func (i *Identity) New() *Identity {
    Identity := &Identity{
        role: make(chan int8, 1),
    }
    Identity.role <- FOLLOWER
    return Identity
}