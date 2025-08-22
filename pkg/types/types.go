package types

import "time"

type ConvertResult struct {
	OrigName string        `json:"orig_name"`
	OutPath  string        `json:"out_path"`
	Err      error         `json:"error"`
	Size     int64         `json:"size"`
	Duration time.Duration `json:"duration"`
}
