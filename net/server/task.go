package server

import "time"

type Task func() (time.Duration, Action)
