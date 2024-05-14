package player

import (
	"github.com/stretchr/testify/mock"
)

func NewMockCommunicator() *MockCommunicator {
	return &MockCommunicator{}
}

type MockCommunicator struct {
	mock.Mock
}

func (mc *MockCommunicator) ReadMessage() (int, []byte, error) {
	return 0, []byte{}, nil
}

func (mc *MockCommunicator) WriteMessage(messageType int, data []byte) error {
	return nil
}

func (mc *MockCommunicator) Close() error {
	return nil
}
