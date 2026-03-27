package httpx_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"go.uber.org/mock/gomock"

	"github.com/wwwangxc/gopkg/httpx"
	"github.com/wwwangxc/gopkg/httpx/mockhttpx"
)

func TestNewClientProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client proxy
	mockedCli := mockhttpx.NewMockClientProxy(ctrl)
	mockedCli.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Head(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Options(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockedCli.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Mock function httpx.NewClientProxy
	patches := gomonkey.ApplyFunc(httpx.NewClientProxy,
		func(string, ...httpx.ClientOption) httpx.ClientProxy {
			return mockedCli
		})
	defer patches.Reset()

	// dosomething...
}
