package model

type AlarmPartition struct {
	Index     int           `json:"index"`
	State     string        `json:"state"`
	Sectors   []AlarmSector `json:"sectors"`
	IsInAlarm bool          `json:"isInAlarm"`
}
