package v1

import (
	"encoding/hex"
	"errors"
)

type CentralCommandResponse struct{}

func (r CentralCommandResponse) Error(response []byte) error {

	var errorCode byte
	if len(response) == 1 {
		errorCode = response[0]
	} else if len(response) == 4 {
		errorCode = response[2]
	} else if len(response) > 4 {
		return nil
	}

	success := byte(0xFE)
	invalid_package := byte(0xE0)
	incorrect_password := byte(0xE1)
	invalid_command := byte(0xE2)
	central_does_not_have_partitions := byte(0xE3)
	open_zones := byte(0xE4)
	command_deprecated := byte(0xE5)
	invalid_model := byte(0xFF)
	bypass_denied := byte(0xE6)
	deactivation_denied := byte(0xE7)
	bypass_central_activated := byte(0xE8)
	unknown_error := byte(0x00)

	switch errorCode {
	case success:
		return errors.New(hex.EncodeToString([]byte{success}))
	case invalid_package:
		return errors.New("invalid package")
	case incorrect_password:
		return errors.New("incorrect password")
	case invalid_command:
		return errors.New("invalid command")
	case central_does_not_have_partitions:
		return errors.New("central does not have partitions")
	case open_zones:
		return errors.New("open zones")
	case command_deprecated:
		return errors.New("command deprecated")
	case invalid_model:
		return errors.New("invalid model")
	case bypass_denied:
		return errors.New("bypass denied")
	case deactivation_denied:
		return errors.New("deactivation denied")
	case bypass_central_activated:
		return errors.New("bypass central activated")
	case unknown_error:
		return errors.New("unknown error")
	default:
		return nil
	}
}
