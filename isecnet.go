package isecnet

import (
	"errors"

	v1 "github.com/luiscovelo/isec/internal/protocol/v1"

	"github.com/luiscovelo/isec/internal/protocol/model"
)

func NewInstance(protocolVersion int) (model.IsecNetProtocol, error) {
	if protocolVersion == 1 {
		return model.IsecNetProtocol{
			CentralCommand:         v1.CentralCommand{},
			CentralCommandResponse: v1.CentralCommandResponse{},
			ServerCommand:          v1.ServerCommand{},
			ServerCommandResponse:  v1.ServerCommandResponse{},
			Parser:                 v1.Parser{},
		}, nil
	}
	return model.IsecNetProtocol{}, errors.New("version not supported")
}
