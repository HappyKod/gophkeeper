package constatns

import "time"

// TimeOutSync sync server timeout
const TimeOutSync = time.Duration(10 * time.Second)

// TimeSleepSync next sync timeout
const TimeSleepSync = time.Duration(5 * time.Second)
