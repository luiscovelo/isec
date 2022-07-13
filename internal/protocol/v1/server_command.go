package v1

import "github.com/luiscovelo/isec/internal/utils"

type ServerCommand struct{}

type commandStruct struct {
	GET_BYTE int
	CONNECT  int
}

type typeDeviceStruct struct {
	ANDROID  int
	GUARDIAN int
}

type connectionTypeStruct struct {
	ETHERNET int
	SIM01    int
}

var commands commandStruct
var typeDevice typeDeviceStruct
var connectionType connectionTypeStruct

func init() {
	commands = commandStruct{
		GET_BYTE: 0xFB,
		CONNECT:  0xE5,
	}
	typeDevice = typeDeviceStruct{
		ANDROID:  0x02,
		GUARDIAN: 0x05,
	}
	connectionType = connectionTypeStruct{
		ETHERNET: 0x45,
		SIM01:    0x47,
	}
}

func (c ServerCommand) AssembleGetByteCommand() []byte {
	var command []byte
	command = append(command, 0x01)
	command = append(command, byte(commands.GET_BYTE))
	command = append(command, utils.Checksum(command))

	return command
}

func (c ServerCommand) AssembleConnectionCommand(macAddress string, identifier string, byteEncrypt byte) []byte {
	var command []byte
	macAddressHex := utils.HexToByte(macAddress)
	identifierHex := utils.HexToByte(identifier)

	command = append(command, byte(len(macAddressHex)+len(identifierHex)+4))
	command = append(command, byte(commands.CONNECT))
	command = append(command, byte(typeDevice.GUARDIAN))
	command = append(command, identifierHex...)
	command = append(command, macAddressHex...)
	command = append(command, 0x00)
	command = append(command, byte(connectionType.ETHERNET))
	command = append(command, utils.Checksum(command))

	newCommand := utils.EncryptCommand(command, byteEncrypt)
	return newCommand
}
