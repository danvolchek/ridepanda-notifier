package worker

import (
	"bytes"
	"time"
)

type Config struct {
	CheckFrequency duration `json:"checkFrequency"`
	Jitter         duration `json:"jitter"`
	StartHour      int      `json:"startHour"`
	StartMinute    int      `json:"startMinute"`
	NotifyNothing  int      `json:"notifyNothing"`
}

type duration struct {
	underlying time.Duration
}

func (d *duration) UnmarshalJSON(data []byte) error {
	dur, err := time.ParseDuration(string(bytes.Trim(data, "\"")))
	if err != nil {
		return err
	}

	d.underlying = dur
	return nil
}
