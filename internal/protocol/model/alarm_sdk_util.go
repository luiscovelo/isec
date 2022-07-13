package model

type AlarmSdkUtil struct {
	SectorPermissions    []bool
	PartitionPermissions []bool
	PartitionSectors     [][]int
	IsDefaultPermission  bool
	SectorsWithStay      []bool
	Pgms                 []AlarmPgm
}
