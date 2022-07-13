package model

type AlarmSector struct {
	Index                 int  `json:"index"`
	IsViolated            bool `json:"isViolated"`
	IsOpen                bool `json:"isOpen"`
	IsBypassed            bool `json:"isBypassed"`
	StayEnabled           bool `json:"stayEnabled"`
	HasCommunicationError bool `json:"hasCommunicationError"`
	HasLowBattery         bool `json:"hasLowBattery"`
}
