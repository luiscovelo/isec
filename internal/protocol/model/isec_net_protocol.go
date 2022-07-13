package model

type CentralCommand interface {
	AssembleActivateCentral(userPassword string, partition int) []byte
	AssembleDeactivateCentral(userPassword string, partition int) []byte
	AssembleGetStatus(userPassword string, model byte) []byte
	AssembleGetCompleteStatus(userPassword string) []byte
	AssembleGetPartialStatus(userPassword string) []byte
	AssembleGetEEPROMPartitions(address int, numberOfBytes int, userPassword string) []byte
	AssembleGetEEPROMSectors(address int, numberOfBytes int, userPassword string) []byte
	AssembleGetEEPROMBytes(address int, numberOfBytes int, userPassword string) []byte
}

type CentralCommandResponse interface {
	Error(response []byte) error
}

type ServerCommand interface {
	AssembleGetByteCommand() []byte
	AssembleConnectionCommand(macAddress string, identifier string, byteEncrypt byte) []byte
}

type ServerCommandResponse interface {
	Error(response byte) error
}

type Parser interface {
	ParseCompleteStatus(bytes []byte, sdkUtil AlarmSdkUtil) AlarmDeviceInfo
}

type IsecNetProtocol struct {
	CentralCommand         CentralCommand
	CentralCommandResponse CentralCommandResponse
	ServerCommand          ServerCommand
	ServerCommandResponse  ServerCommandResponse
	Parser                 Parser
}
