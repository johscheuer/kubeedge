package messagelayer

import (
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/config"
)

// ContextMessageLayer build on context
type ContextMessageLayer struct {
	SendModuleName     string
	ReceiveModuleName  string
	ResponseModuleName string
}

// Send message
func (cml *ContextMessageLayer) Send(message model.Message) error {
	beehiveContext.Send(cml.SendModuleName, message)
	return nil
}

// Receive message
func (cml *ContextMessageLayer) Receive() (model.Message, error) {
	return beehiveContext.Receive(cml.ReceiveModuleName)
}

// Response message
func (cml *ContextMessageLayer) Response(message model.Message) error {
	beehiveContext.Send(cml.ResponseModuleName, message)
	return nil
}

// NewContextMessageLayer create a ContextMessageLayer
func NewContextMessageLayer() (*ContextMessageLayer, error) {
	return &ContextMessageLayer{
		SendModuleName:     config.ContextSendModule,
		ReceiveModuleName:  config.ContextReceiveModule,
		ResponseModuleName: config.ContextResponseModule,
	}, nil
}
