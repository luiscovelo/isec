package model

type AlarmPartitions struct {
	IsEnabled       bool             `json:"isEnabled"`
	SupportStayMode bool             `json:"supportStayMode"`
	PartitionsList  []AlarmPartition `json:"partitionsList"`
}
