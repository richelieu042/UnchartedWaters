package conf

import "time"

var (
	defSleepInterval = 1_300 * time.Millisecond

	battleSleepInterval = 5_000 * time.Millisecond
)

func SetDefSleepInterval(interval time.Duration) {
	defSleepInterval = interval
}

func GetDefSleepInterval() time.Duration {
	return defSleepInterval
}

func GetBattleSleepInterval() time.Duration {
	return battleSleepInterval
}
