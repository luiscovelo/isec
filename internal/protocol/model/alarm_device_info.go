package model

type AlarmDeviceInfo struct {
	Model             AlarmModel      `json:"model"`
	FirmwareVersion   string          `json:"firmwareVersion"`
	CentralHasProblem bool            `json:"centralHasProblem"`
	Time              string          `json:"time"`
	Sectors           []AlarmSector   `json:"sectors"`
	Partitions        AlarmPartitions `json:"partitions"`
	Pgms              []AlarmPgm      `json:"pgms"`
	GeneralState      string          `json:"generalState"`
	IsInAlarm         bool            `json:"isInAlarm"`
	PartialMode       string          `json:"partialMode"`
}
