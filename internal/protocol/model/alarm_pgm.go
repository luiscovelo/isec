package model

type AlarmPgm struct {
	Position            int  `json:"position"`
	IsActivated         bool `json:"isActivated"`
	HasLowBattery       bool `json:"hasLowBattery"`
	HasCommucationError bool `json:"hasCommunicationError"`
}
