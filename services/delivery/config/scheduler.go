package config

import "github.com/kanthorlabs/common/validator"

var ServiceNameScheduler = "scheduler"

type Scheduler struct {
	Request SchedulerRequest `json:"request" yaml:"request" mapstructure:"request"`
}

func (c *Scheduler) Validate() error {
	if err := c.Request.Validate(); err != nil {
		return err
	}

	return nil
}

type SchedulerRequest struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (conf *SchedulerRequest) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThanOrEqual("SCHEDULER.CONFIG.REQUEST.SCHEDULE.TIMEOUT", conf.Timeout, 1000),
	)
}
