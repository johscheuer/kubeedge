package devicetwin

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/edge/mocks/beego"
	"github.com/kubeedge/kubeedge/edge/mocks/beehive"
	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"github.com/kubeedge/kubeedge/edge/pkg/devicetwin/dtcommon"
)

const (
	//TestModule is name of test.
	TestModule = "test"
	//DeviceTwinModuleName is name of twin
	DeviceTwinModuleName = "twin"
)

// TestName is function to test Name().
func TestName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "DeviceTwinNametest",
			want: "twin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := &DeviceTwin{}
			if got := dt.Name(); got != tt.want {
				t.Errorf("DeviceTwin.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGroup is function to test Group().
func TestGroup(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "DeviceTwinGroupTest",
			want: "twin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := &DeviceTwin{}
			if got := dt.Group(); got != tt.want {
				t.Errorf("DeviceTwin.Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStart is function to test Start().
func TestStart(t *testing.T) {
	beehiveContext.InitContext(beehiveContext.MsgCtxTypeChannel)
	//test is for sending test messages from devicetwin module.
	var test model.Message
	// ormerMock is mocked Ormer implementation.
	var ormerMock *beego.MockOrmer
	// querySeterMock is mocked QuerySeter implementation.
	var querySeterMock *beego.MockQuerySeter
	// fakeModule is mocked implementation of TestModule.
	var fakeModule *beehive.MockModule

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ormerMock = beego.NewMockOrmer(mockCtrl)
	querySeterMock = beego.NewMockQuerySeter(mockCtrl)
	fakeModule = beehive.NewMockModule(mockCtrl)
	dbm.DBAccess = ormerMock

	fakeModule.EXPECT().Name().Return(TestModule).Times(3)
	core.Register(fakeModule)
	beehiveContext.AddModule(TestModule)

	core.Register(&DeviceTwin{})
	dt := DeviceTwin{}
	beehiveContext.AddModule(dt.Name())
	beehiveContext.AddModuleGroup(dt.Name(), dt.Group())
	ormerMock.EXPECT().QueryTable(gomock.Any()).Return(querySeterMock).Times(1)
	querySeterMock.EXPECT().All(gomock.Any()).Return(int64(1), nil).Times(1)
	go dt.Start()
	time.Sleep(1 * time.Millisecond)
	// Sending a message from devicetwin module to the created fake module(TestModule) to check context is initialized properly.
	beehiveContext.Send(TestModule, test)
	_, err := beehiveContext.Receive(TestModule)
	t.Run("MessagePingTest", func(t *testing.T) {
		if err != nil {
			t.Errorf("Error while receiving message: %v", err)
			return
		}
	})
	//Checking whether Mem,Twin,Device and Comm modules are registered and started successfully.
	tests := []struct {
		name       string
		moduleName string
	}{
		{
			name:       "MemModuleHealthCheck",
			moduleName: dtcommon.MemModule,
		},
		{
			name:       "TwinModuleHealthCheck",
			moduleName: dtcommon.TwinModule,
		},
		{
			name:       "DeviceModuleHealthCheck",
			moduleName: dtcommon.DeviceModule,
		},
		{
			name:       "CommModuleHealthCheck",
			moduleName: dtcommon.CommModule,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			moduleCheck := false
			for _, module := range dt.DTModules {
				if test.moduleName == module.Name {
					moduleCheck = true
					err := dt.DTContexts.HeartBeat(test.moduleName, "ping")
					if err != nil {
						t.Errorf("Heartbeat of module %v is expired and dtcontroller will start it again", test.moduleName)
					}
					break
				}
			}
			if moduleCheck == false {
				t.Errorf("Registration of module %v failed", test.moduleName)
			}
		})
	}
}

// TestCleanup is function to test Cleanup().
func TestCleanup(t *testing.T) {
	beehiveContext.InitContext(beehiveContext.MsgCtxTypeChannel)

	var test model.Message
	beehiveContext.AddModule(TestModule)

	core.Register(&DeviceTwin{})
	dt := DeviceTwin{}
	beehiveContext.AddModule(dt.Name())
	beehiveContext.AddModuleGroup(dt.Name(), dt.Group())
	_, cancel := context.WithCancel(context.Background())
	deviceTwin := DeviceTwin{
		cancel: cancel,
	}
	deviceTwin.Cleanup()
	//Send message to avoid deadlock if channel deletion has failed after cleanup
	go beehiveContext.Send(DeviceTwinModuleName, test)
	_, err := beehiveContext.Receive(DeviceTwinModuleName)
	t.Run("CheckCleanUp", func(t *testing.T) {
		if err == nil {
			t.Errorf("DeviceTwin Module still has channel after cleanup")
		}
	})
}
