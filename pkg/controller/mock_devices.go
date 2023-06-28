package controller

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

type MockDHT11 struct {
	mock.Mock
}

type MockDevice struct {
	mock.Mock
}

func (f *MockDevice) ReadMoisture(tries int) (float32, error) {
	args := f.Called(tries)
	fn, ok := args.Get(0).(func() float32)
	if !ok {
		return 0, fmt.Errorf("couldn't unpack mock read moisture function call")
	}
	return fn(), args.Error(1)
}

func (m *MockDHT11) ReadTempHumidity(tries int) (float32, float32, error) {
	args := m.Called(tries)
	fn, ok := args.Get(0).(func() (float32, float32, error))
	if !ok {
		return 0, 0, fmt.Errorf("couldn't unpack mock temp humidity function call")
	}
	return fn()
}

func (f *MockDevice) Name() string {
	args := f.Called()
	return args.String(0)
}

func (f *MockDevice) On() {
	_ = f.Called()
	return
}

func (f *MockDevice) Off() {
	_ = f.Called()
	return
}

func (f *MockDevice) Number() uint8 {
	return 1
}

func (f *MockDevice) State() (bool, error) {
	args := f.Called()
	fn, ok := args.Get(0).(func() (bool, error))
	if !ok {
		return false, fmt.Errorf("couldn't unpack mock state function call")
	}
	return fn()
}
