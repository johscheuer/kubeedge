package context

import (
	"sync"
	"time"

	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/core/model"
)

//define channel type
const (
	MsgCtxTypeChannel = "channel"
)

var (
	// singleton
	context *Context
	once    sync.Once
)

// InitContext gets global context instance
func InitContext(contextType string) {
	once.Do(func() {
		context = &Context{}
		switch contextType {
		case MsgCtxTypeChannel:
			channelContext := NewChannelContext()
			context.messageContext = channelContext
			context.moduleContext = channelContext
		default:
			klog.Fatalf("Do not support context type:%s", contextType)
		}
	})
}

// AddModule adds module into module context
func AddModule(module string) {
	context.moduleContext.AddModule(module)
}

// AddModuleGroup adds module into module context group
func AddModuleGroup(module, group string) {
	context.moduleContext.AddModuleGroup(module, group)
}

// Cleanup cleans up module
func Cleanup(module string) {
	context.moduleContext.Cleanup(module)
}

// Send the message
func Send(module string, message model.Message) {
	context.messageContext.Send(module, message)
}

// Receive the message
// module : local module name
func Receive(module string) (model.Message, error) {
	message, err := context.messageContext.Receive(module)
	if err == nil {
		return message, nil
	}
	klog.Warningf("Receive: failed to receive message, error:%v", err)
	return message, err
}

// SendSync sends message in sync mode
// module: the destination of the message
// timeout: if <= 0 using default value(30s)
func SendSync(module string,
	message model.Message, timeout time.Duration) (model.Message, error) {
	resp, err := context.messageContext.SendSync(module, message, timeout)
	if err == nil {
		return resp, nil
	}
	return model.Message{}, err
}

// SendResp sends response
// please get resp message using model.NewRespByMessage
func SendResp(resp model.Message) {
	context.messageContext.SendResp(resp)
}

// SendToGroup broadcasts the message to all of group members
func SendToGroup(moduleType string, message model.Message) {
	context.messageContext.SendToGroup(moduleType, message)
}

// sendToGroupSync broadcasts the message to all of group members in sync mode
func sendToGroupSync(moduleType string, message model.Message, timeout time.Duration) error {
	return context.messageContext.SendToGroupSync(moduleType, message, timeout)
}
