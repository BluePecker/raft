package types

import (
    "time"
    "errors"
)

const SAFE_TIME = 5

type Second struct {
    Tenure, Timeout, Heartbeat int
}

func (s Second) Validator() error {
    if s.Tenure < 0 || s.Timeout < 0 || s.Heartbeat < 0 {
        return errors.New("tenure can not less than 0")
    }
    
    if s.Tenure - s.Timeout - s.Heartbeat < SAFE_TIME {
        return errors.New("tenure - timeout - heartbeat there is not enough safe space")
    }
    
    return nil
}

type Clock struct {
    Second Second
    Timer  struct {
               Timeout *time.Timer
               Tenure  *time.Timer
           }
    Ticker struct {
               Heartbeat *time.Ticker
           }
}