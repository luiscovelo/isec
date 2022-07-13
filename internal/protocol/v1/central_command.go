package v1

import (
	"github.com/luiscovelo/isec/internal/protocol/model"
	"github.com/luiscovelo/isec/internal/utils"
)

type CentralCommand struct{}

const (
	_ISEC_PROGRAM       = 0xE9
	_EEPROM_BYTES       = 0x5C
	_ACTIVATE_CENTRAL   = 0x41
	_DEACTIVATE_CENTRAL = 0x44
	_COMPLETE_STATUS    = 0x5B
	_PARTIAL_STATUS     = 0x5A
	_DELIMITER          = 0x21
	_PARTITION_A        = 0x41
	_PARTITION_B        = 0x42
	_PARTITION_C        = 0x43
	_PARTITION_D        = 0x44
)

func (c CentralCommand) AssembleActivateCentral(userPassword string, partition int) []byte {
	var command []byte
	command = append(command, _ACTIVATE_CENTRAL)

	partitionIndex := c.getPartitionIndex(partition)
	if partitionIndex != 0 {
		command = append(command, partitionIndex)
	}

	return c.assembleIsec(command, userPassword)
}

func (c CentralCommand) AssembleDeactivateCentral(userPassword string, partition int) []byte {
	var command []byte
	command = append(command, _DEACTIVATE_CENTRAL)

	partitionIndex := c.getPartitionIndex(partition)
	if partitionIndex != 0 {
		command = append(command, partitionIndex)
	}

	return c.assembleIsec(command, userPassword)
}

func (c CentralCommand) AssembleGetStatus(userPassword string, modelByte byte) []byte {
	if modelByte == model.AMT_4010_SMART {
		return c.AssembleGetCompleteStatus(userPassword)
	}
	return c.AssembleGetPartialStatus(userPassword)
}

func (c CentralCommand) AssembleGetPermissions(userPassword string) []byte {
	return c.AssembleGetEEPROMBytes(0x00C0, 33, userPassword)
}

func (c CentralCommand) AssembleGetCompleteStatus(userPassword string) []byte {
	var command []byte
	command = append(command, _COMPLETE_STATUS)
	return c.assembleIsec(command, userPassword)
}

func (c CentralCommand) AssembleGetPartialStatus(userPassword string) []byte {
	var command []byte
	command = append(command, _PARTIAL_STATUS)
	return c.assembleIsec(command, userPassword)
}

func (c CentralCommand) AssembleGetEEPROMPartitions(address int, numberOfBytes int, userPassword string) []byte {
	return c.AssembleGetEEPROMBytes(address, numberOfBytes, userPassword)
}

func (c CentralCommand) AssembleGetEEPROMSectors(address int, numberOfBytes int, userPassword string) []byte {
	return c.AssembleGetEEPROMBytes(address, numberOfBytes, userPassword)
}

func (c CentralCommand) AssembleGetEEPROMBytes(address int, numberOfBytes int, userPassword string) []byte {
	var command []byte
	highAddress := utils.ExtractCrcHigh(address)
	lowAddress := utils.ExtractCrcLow(address)

	command = append(command, _EEPROM_BYTES)
	command = append(command, byte(highAddress))
	command = append(command, byte(lowAddress))
	command = append(command, byte(numberOfBytes))
	return c.assembleIsec(command, userPassword)
}

func (c CentralCommand) assembleIsec(command []byte, userPassword string) []byte {
	var response []byte

	commandSize := len(command)
	userPasswordSize := len(userPassword)

	response = append(response, byte(commandSize+userPasswordSize+3))
	response = append(response, _ISEC_PROGRAM)
	response = append(response, _DELIMITER)
	response = append(response, utils.StringToByte(userPassword)...)
	response = append(response, command...)
	response = append(response, _DELIMITER)
	response = append(response, utils.Checksum(response))

	return response
}

func (c CentralCommand) getPartitionIndex(partition int) byte {
	switch partition {
	case 1:
		return _PARTITION_A
	case 2:
		return _PARTITION_B
	case 3:
		return _PARTITION_C
	case 4:
		return _PARTITION_D
	default:
		return 0x00
	}
}
