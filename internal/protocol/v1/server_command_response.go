package v1

import "errors"

type ServerCommandResponse struct{}

func (r ServerCommandResponse) Error(response byte) error {

	success := byte(0xFE)
	success_ethernet := byte(0x45)
	success_gprs := byte(0x47)
	different_checksum := byte(0xE6)
	central_not_connected := byte(0xE4)
	connected_to_other_device := byte(0xE8)
	unknown_error := byte(0x00)

	switch response {
	case success, success_ethernet, success_gprs:
		return nil
	case different_checksum:
		return nil
	case central_not_connected:
		return errors.New("central not connected")
	case connected_to_other_device:
		return errors.New("connected to other device")
	case unknown_error:
		return errors.New("unknown error")
	default:
		return errors.New("unknown error")
	}
}
