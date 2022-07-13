package v1

import (
	"encoding/hex"
	"time"

	"github.com/luiscovelo/isec/internal/protocol/model"
	"github.com/luiscovelo/isec/internal/utils"
)

type Parser struct{}

func (Parser) ParseCompleteStatus(bytes []byte, sdkUtil model.AlarmSdkUtil) model.AlarmDeviceInfo {
	centralModelByte := bytes[26]
	firmwareVersionByte := bytes[27]
	timeBytes := bytes[32:37]
	outputStateByte := bytes[47]
	extenderPgmsStateBytes := []byte{}

	if len(bytes) >= 56 {
		extenderPgmsStateBytes = bytes[54:56]
	}

	alarmModel := model.GetAlarmModel(int(centralModelByte))
	firmwareVersion := ParseFirmwareVersion(firmwareVersionByte)
	time := ParseTime(timeBytes)

	pgms := ParsePgms(outputStateByte, centralModelByte, extenderPgmsStateBytes)

	alarmSectors := ParseAlarmSectors(bytes, centralModelByte, sdkUtil)
	partitionBytes := GetPartitionBytes(bytes, centralModelByte)

	alarmPartitions := ParseAlarmPartitions(partitionBytes, outputStateByte, centralModelByte, alarmSectors, sdkUtil)

	var generalState string
	if alarmPartitions.IsEnabled {
		generalState = ParsePartitionGeneralState(alarmPartitions.PartitionsList)
	} else {
		if len(alarmPartitions.PartitionsList) > 0 {
			generalState = alarmPartitions.PartitionsList[0].State
		} else {
			generalState = "NOT_AUTHORIZED"
		}
	}

	isSirenOn := false
	if len(alarmPartitions.PartitionsList) > 0 {
		isSirenOn = alarmPartitions.PartitionsList[0].IsInAlarm
	}

	partialMode := GetPartialMode(centralModelByte)

	alarmDeviceInfo := model.AlarmDeviceInfo{
		Model:           alarmModel,
		FirmwareVersion: firmwareVersion,
		Sectors:         alarmSectors,
		Partitions:      alarmPartitions,
		IsInAlarm:       isSirenOn,
		Time:            time,
		Pgms:            pgms,
		GeneralState:    generalState,
		PartialMode:     partialMode,
	}
	return alarmDeviceInfo
}

func ParseCentralModel(value byte) model.AlarmModel {
	return model.AlarmModel{}
}

func ParseFirmwareVersion(value byte) string {
	version := hex.EncodeToString([]byte{value})
	major := version[0]
	minor := version[1]
	return string(major) + "." + string(minor)
}

func ParseTime(bytes []byte) string {
	hour := int(bytes[0])
	minute := int(bytes[1])
	day := int(bytes[2])
	month := int(bytes[3])
	year := int(bytes[4]) + 2000
	datetime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
	datetimeFormatted := datetime.Format("2006/01/02 15:04")
	return datetimeFormatted
}

func ParsePgms(outputStateByte byte, modelByte byte, pgmsExtenderBytes []byte) []model.AlarmPgm {
	var pgmList = make([]model.AlarmPgm, 0)
	switch modelByte {
	case model.AMT_4010_SMART:
		pgms := utils.ByteToBooleanSlice([]byte{outputStateByte})
		pgmList = append(pgmList, model.AlarmPgm{IsActivated: pgms[6]})
		pgmList = append(pgmList, model.AlarmPgm{IsActivated: pgms[5]})
		pgmList = append(pgmList, model.AlarmPgm{IsActivated: pgms[4]})
		if len(pgmsExtenderBytes) > 0 {
			for _, pgm := range pgmsExtenderBytes {
				pgms = utils.ByteToBooleanSlice([]byte{pgm})
				for _, bool := range pgms {
					pgmList = append(pgmList, model.AlarmPgm{IsActivated: bool})
				}
			}
		}
		return pgmList
	default:
		return pgmList
	}
}

func ParseAlarmSectors(bytes []byte, modelByte byte, sdkUtil model.AlarmSdkUtil) []model.AlarmSector {
	sectorStatusBytes := getSectorStatusBytes(bytes, modelByte)

	openSectors := utils.ByteToBooleanSlice(sectorStatusBytes[0])
	violatedSectors := utils.ByteToBooleanSlice(sectorStatusBytes[1])
	bypassedSectors := utils.ByteToBooleanSlice(sectorStatusBytes[2])
	lowBatterySectors := utils.ByteToBooleanSlice(sectorStatusBytes[3])

	// TODO: implements to v2
	// staySectors := parseStayMode(sdkUtil, model)

	var alarmSectors = make([]model.AlarmSector, 0)

	for index, permission := range sdkUtil.SectorPermissions {
		hasLowBattery := false
		if index < len(lowBatterySectors) {
			hasLowBattery = lowBatterySectors[index]
		}

		if permission {
			alarmSector := model.AlarmSector{
				Index:         index,
				IsOpen:        openSectors[index],
				IsViolated:    violatedSectors[index],
				IsBypassed:    bypassedSectors[index],
				HasLowBattery: hasLowBattery,
			}
			alarmSectors = append(alarmSectors, alarmSector)
		}

	}

	return alarmSectors
}

func getSectorStatusBytes(bytes []byte, modelByte byte) [][]byte {
	var sectorStatus = make([][]byte, 0)
	switch modelByte {
	case model.AMT_4010_SMART:
		sectorStatus = append(sectorStatus, bytes[2:10])
		sectorStatus = append(sectorStatus, bytes[10:18])
		sectorStatus = append(sectorStatus, bytes[18:26])
		sectorStatus = append(sectorStatus, bytes[48:54])
		return sectorStatus
	case model.ANM_24_NET:
		sectorStatus = append(sectorStatus, bytes[2:5])
		sectorStatus = append(sectorStatus, bytes[8:11])
		sectorStatus = append(sectorStatus, bytes[14:17])
		sectorStatus = append(sectorStatus, bytes[40:43])
		return sectorStatus
	default:
		sectorStatus = append(sectorStatus, bytes[2:8])
		sectorStatus = append(sectorStatus, bytes[8:14])
		sectorStatus = append(sectorStatus, bytes[14:20])
		sectorStatus = append(sectorStatus, bytes[40:45])
		return sectorStatus
	}
}

func GetDefaultUserPermissions(modelByte byte) model.AlarmSdkUtil {
	partitions := getDefaultPartitionsPermissions(modelByte)
	sectors := getDefaultSectorsPermissions(modelByte)
	partitionSectors := getDefaultPartitionsSectors(modelByte)

	sdk := model.AlarmSdkUtil{
		SectorPermissions:    sectors,
		PartitionPermissions: partitions,
		PartitionSectors:     partitionSectors,
		IsDefaultPermission:  true,
	}
	return sdk
}

func getDefaultSectorsPermissions(model byte) []bool {
	numberOfSectors := getNumberOfSectors(model)
	var sectors = make([]bool, numberOfSectors)
	for i := 0; i < len(sectors); i++ {
		sectors[i] = true
	}
	return sectors
}

func getDefaultPartitionsPermissions(model byte) []bool {
	numberOfPartitions := getNumberOfPartitions(model)
	var partitions = make([]bool, numberOfPartitions)
	for i := 0; i < len(partitions); i++ {
		partitions[i] = true
	}
	return partitions
}

func getDefaultPartitionsSectors(model byte) [][]int {
	numberOfPartitions := getNumberOfPartitions(model) + 1
	var partitions = make([][]int, numberOfPartitions)
	for i := 0; i < len(partitions); i++ {
		if i == 0 {
			var sectors = make([]int, getNumberOfSectors(model))
			for i := 0; i < len(sectors); i++ {
				sectors[i] = i
			}
			partitions[i] = sectors
		} else {
			partitions[i] = []int{}
		}
	}
	return partitions
}

func getNumberOfPartitions(modelByte byte) int {
	switch modelByte {
	case model.ANM_24_NET:
		return 0
	case model.AMT_4010_SMART:
		return 4
	default:
		return 0
	}
}

func getNumberOfSectors(modelByte byte) int {
	switch modelByte {
	case model.ANM_24_NET:
		return 24
	case model.AMT_4010_SMART:
		return 64
	default:
		return 48
	}
}

// func getNumberOfUsers(model byte) int {
// 	switch model {
// 	case model.ANM_24_NET:
// 		return 30
// 	default:
// 		return 64
// 	}
// }

func ParseAlarmPartitions(partitionBytes [][]byte, outputStateByte byte, modelByte byte, sectors []model.AlarmSector, permissions model.AlarmSdkUtil) model.AlarmPartitions {
	partitionEnabledByte := partitionBytes[0]
	partitionsActivatedBytes := partitionBytes[1]

	partitionEnabled := false
	if partitionEnabledByte[0] == 0x01 {
		partitionEnabled = true
	}

	partitionsActivated := parsePartitions(partitionsActivatedBytes)

	partitionPermissions := permissions.PartitionPermissions
	sectorsIndex := permissions.PartitionSectors

	generalState := parseGeneralState(partitionsActivated, modelByte)
	isSirenOn := getSirenStatus(outputStateByte, modelByte)

	var partitionsList = make([]model.AlarmPartition, 0)

	if partitionEnabled {

		generalPartition := parseAlarmPartition(0, generalState, sectorsIndex[0], sectors, isSirenOn)
		partitionsList = append(partitionsList, generalPartition)

		for index, permission := range partitionPermissions {
			if permission {

				state := "DEACTIVATED"
				if partitionsActivated[index] {
					state = "ACTIVATED"
				}

				partitionIndex := index + 1
				partition := parseAlarmPartition(partitionIndex, state, sectorsIndex[partitionIndex], sectors, isSirenOn)

				partitionsList = append(partitionsList, partition)
			}
		}

	} else {
		generalPartition := parseAlarmPartition(0, generalState, []int{}, sectors, isSirenOn)
		partitionsList = append(partitionsList, generalPartition)
	}

	return model.AlarmPartitions{
		IsEnabled:       partitionEnabled,
		SupportStayMode: false,
		PartitionsList:  partitionsList,
	}

}

func parseAlarmPartition(partitionIndex int, state string, sectorsIndex []int, sectorsList []model.AlarmSector, isSirenOn bool) model.AlarmPartition {

	var sectors = make([]model.AlarmSector, 0)

	if len(sectorsIndex) > 0 {
		for _, sectorIndex := range sectorsIndex {

			for i := 0; i < len(sectorsList); i++ {
				if sectorsList[i].Index == sectorIndex {
					alarm := sectorsList[i]
					sectors = append(sectors, alarm)
					break
				}
			}

		}
	}

	return model.AlarmPartition{
		Index:     partitionIndex,
		State:     state,
		Sectors:   sectors,
		IsInAlarm: isSirenOn,
	}
}

func parseGeneralState(partitions []bool, modelByte byte) string {
	switch modelByte {
	case model.ANM_24_NET:
		return ""
	default:
		return parseAlarmState(partitions)
	}
}

func parseAlarmState(partitions []bool) string {

	contains := false
	someEquals := true
	for i := 1; i < len(partitions); i++ {
		if partitions[i] != partitions[0] {
			someEquals = false
			break
		}
	}

	for _, item := range partitions {
		if item {
			contains = true
			break
		}
	}

	if someEquals {
		state := partitions[0]
		if state {
			return "ACTIVATED"
		} else {
			return "DEACTIVATED"
		}
	}

	if contains {
		return "PARTIAL"
	}

	return "NOT_AUTHORIZED"
}

func getSirenStatus(outputStateByte byte, modelByte byte) bool {
	outputStateList := utils.ByteToBooleanSlice([]byte{outputStateByte})
	switch modelByte {
	case model.AMT_4010_SMART:
		return outputStateList[3]
	default:
		return outputStateList[2]
	}
}

func parsePartitions(bytes []byte) []bool {
	booleans := utils.ByteToBooleanSlice(bytes)
	partitions := make([]bool, 0)

	if len(bytes) > 1 {
		partitions = append(partitions, booleans[0:2]...)
		partitions = append(partitions, booleans[8:10]...)
	} else {
		partitions = append(partitions, booleans[0:2]...)
	}

	return partitions
}

func GetPartitionBytes(bytes []byte, modelByte byte) [][]byte {
	var partition = make([][]byte, 0)
	switch modelByte {
	case model.AMT_4010_SMART:
		partition = append(partition, []byte{bytes[28]})
		partition = append(partition, bytes[29:31])
		return partition
	default:
		partition = append(partition, []byte{bytes[22]})
		partition = append(partition, []byte{bytes[23]})
		return partition
	}
}

func ParsePartitionGeneralState(partitions []model.AlarmPartition) string {
	partitionList := partitions[1:]
	var states = make([]string, 0)

	for _, partition := range partitionList {
		if partition.Index > 0 {
			states = append(states, partition.State)
		}
	}

	contains := false
	someEquals := true
	for i := 1; i < len(states); i++ {
		if states[i] != states[0] {
			someEquals = false
			break
		}
	}

	for _, item := range states {
		if item == "ACTIVATED" {
			contains = true
			break
		}
	}

	if someEquals && len(states) > 0 {
		state := states[0]
		if state == "ACTIVATED" {
			return "ACTIVATED"
		} else if state == "DEACTIVATED" {
			return "DEACTIVATED"
		}
	}

	if contains {
		return "PARTIAL"
	}

	return "NOT_AUTHORIZED"
}

func GetPartialMode(modelByte byte) string {
	switch modelByte {
	case model.AMT_4010_SMART:
		return "NONE"
	case model.ANM_24_NET:
		return "STAY"
	default:
		return "PARTIAL"
	}
}

func ParseUserPermissions(partitionEnabled bool, partitionBytes []byte, sectorsBytes []byte, modelByte byte) model.AlarmSdkUtil {

	userPartitions := ParseUserPartitions(partitionBytes, modelByte)
	userSectors, partitionSectors := ParseUserSectors(sectorsBytes, partitionEnabled, modelByte, userPartitions)

	return model.AlarmSdkUtil{
		SectorPermissions:    userSectors,
		PartitionPermissions: userPartitions,
		PartitionSectors:     partitionSectors,
		IsDefaultPermission:  false,
	}
}

func ParseUserPartitions(bytes []byte, modelByte byte) []bool {

	var partitions = make([]bool, 0)
	userIndex := bytes[2]

	if userIndex == 0 {

		if modelByte == model.AMT_4010_SMART {
			partitions = append(partitions, true)
			partitions = append(partitions, true)
			partitions = append(partitions, true)
			partitions = append(partitions, true)
			return partitions
		} else {
			partitions = append(partitions, true)
			partitions = append(partitions, true)
			return partitions
		}

	} else {

		index := (int(userIndex) - 1) / 8

		bytePartitionA := bytes[3+index]
		bytePartitionB := bytes[11+index]

		var bit int
		if index == 0 {
			bit = int(userIndex) - 1
		} else {
			bit = (int(userIndex) - 1) % 8
		}

		booleansPartA := utils.ByteToBooleanSlice([]byte{bytePartitionA})
		partA := booleansPartA[bit]

		booleansPartB := utils.ByteToBooleanSlice([]byte{bytePartitionB})
		partB := booleansPartB[bit]

		if modelByte == model.AMT_4010_SMART {

			bytePartitionC := bytes[19+index]
			bytePartitionD := bytes[27+index]

			booleansPartC := utils.ByteToBooleanSlice([]byte{bytePartitionC})
			partC := booleansPartC[bit]

			booleansPartD := utils.ByteToBooleanSlice([]byte{bytePartitionD})
			partD := booleansPartD[bit]

			partitions = append(partitions, partA)
			partitions = append(partitions, partB)
			partitions = append(partitions, partC)
			partitions = append(partitions, partD)
			return partitions

		} else {
			partitions = append(partitions, partA)
			partitions = append(partitions, partB)
			return partitions
		}

	}
}

func ParseUserSectors(bytes []byte, partitionEnabled bool, modelByte byte, userPartitions []bool) ([]bool, [][]int) {

	var sectors = make([]bool, 0)

	sectorsPerPartition := make([][]int, len(userPartitions))
	sectorsWithoutPartition := make([]int, 0)

	numberOfSectors := getNumberOfSectors(modelByte)

	if !partitionEnabled {

		for sectorIndex := 0; sectorIndex < numberOfSectors; sectorIndex++ {
			index := sectorIndex / 8
			byteActive := bytes[67+index]
			var bit int
			if index == 0 {
				bit = sectorIndex
			} else {
				bit = sectorIndex % 8
			}
			booleansActive := utils.ByteToBooleanSlice([]byte{byteActive})
			active := booleansActive[bit]

			if active {
				sectorsWithoutPartition = append(sectorsWithoutPartition, sectorIndex)
			}

			sectors = append(sectors, active)
		}

	} else {

		for sectorIndex := 0; sectorIndex < numberOfSectors; sectorIndex++ {
			index := sectorIndex / 8
			byteActive := bytes[67+index]
			var bit int
			if index == 0 {
				bit = sectorIndex
			} else {
				bit = sectorIndex % 8
			}
			booleansActive := utils.ByteToBooleanSlice([]byte{byteActive})
			active := booleansActive[bit]

			noPartitionFlag := true
			partitionList := make([]bool, 0)

			bytePartitionA := bytes[3+index]
			bytePartitionB := bytes[11+index]

			booleansPartA := utils.ByteToBooleanSlice([]byte{bytePartitionA})
			partA := booleansPartA[bit]
			partitionList = append(partitionList, partA)

			booleansPartB := utils.ByteToBooleanSlice([]byte{bytePartitionB})
			partB := booleansPartB[bit]
			partitionList = append(partitionList, partB)

			var isUserSector bool
			if modelByte == model.AMT_4010_SMART {

				bytePartitionC := bytes[19+index]
				bytePartitionD := bytes[27+index]

				booleansPartC := utils.ByteToBooleanSlice([]byte{bytePartitionC})
				partC := booleansPartC[bit]
				partitionList = append(partitionList, partC)

				booleansPartD := utils.ByteToBooleanSlice([]byte{bytePartitionD})
				partD := booleansPartD[bit]
				partitionList = append(partitionList, partD)

				userPartitionA := partA
				userPartitionB := partB
				userPartitionC := partC
				userPartitionD := partD

				if !userPartitions[0] {
					userPartitionA = false
				}

				if !userPartitions[1] {
					userPartitionB = false
				}

				if !userPartitions[2] {
					userPartitionC = false
				}

				if !userPartitions[3] {
					userPartitionD = false
				}

				if active && ((!partA && !partB && !partC && !partD) || userPartitionA || userPartitionB || userPartitionC || userPartitionD) {
					isUserSector = true
				} else {
					isUserSector = false
				}

			} else {
				userPartitionA := partA
				userPartitionB := partB

				if !userPartitions[0] {
					userPartitionA = false
				}

				if !userPartitions[1] {
					userPartitionB = false
				}

				if active && ((!partA && !partB) || userPartitionA || userPartitionB) {
					isUserSector = true
				} else {
					isUserSector = false
				}
			}

			for partitionIndex, isFromPartition := range partitionList {
				if active && isFromPartition {
					sectorsPerPartition[partitionIndex] = append(sectorsPerPartition[partitionIndex], sectorIndex)
					noPartitionFlag = false
				}
			}

			if noPartitionFlag && active {
				sectorsWithoutPartition = append(sectorsWithoutPartition, sectorIndex)
			}

			sectors = append(sectors, isUserSector)

		}

	}

	partitionsWithPermissions := make([][]int, 0)

	for index, list := range sectorsPerPartition {
		if userPartitions[index] {
			partitionsWithPermissions = append(partitionsWithPermissions, list)
		} else {
			partitionsWithPermissions = append(partitionsWithPermissions, nil)
		}
	}

	partitions := make([][]int, 0)

	partitions = append(partitions, sectorsWithoutPartition)
	partitions = append(partitions, partitionsWithPermissions...)

	return sectors, partitions

}
