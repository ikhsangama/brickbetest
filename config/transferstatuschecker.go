package config

import "strconv"

type TransferStatusCheckerConfig struct {
	CronSpec     string
	IntervalDays int
	FetchLimit   int
}

func GetTransferStatusCheckerConfig() TransferStatusCheckerConfig {
	const (
		defaultCron         = "0 * * * * *"
		defaultIntervalDays = "3" // for example if looking for transfer request stuck more than 3 days, put 3. For test purpose just put 0.
		defaultFetchLimit   = "3"
	)

	d := getEnvWithDefault("INTERVAL_DAYS", defaultIntervalDays)
	intervalDays, err := strconv.Atoi(d)
	if err != nil {
		panic(err)
	}

	l := getEnvWithDefault("FETCH_LIMIT", defaultFetchLimit)
	fetchLimit, err := strconv.Atoi(l)
	if err != nil {
		panic(err)
	}

	return TransferStatusCheckerConfig{
		CronSpec:     getEnvWithDefault("CRON_SPEC", defaultCron),
		IntervalDays: intervalDays,
		FetchLimit:   fetchLimit,
	}
}
