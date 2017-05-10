package types

import "time"

type Clock struct {
    Second struct {
               tenure, timeout, heartbeat int
           }
    Timer  struct {
               timeout *time.Timer
               tenure  *time.Timer
           }
    Ticker struct {
               heartbeat *time.Ticker
           }
}