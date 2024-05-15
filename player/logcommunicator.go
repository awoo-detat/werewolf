package player

import (
	"fmt"
)

func NewLogCommunicator() *LogCommunicator {
	return &LogCommunicator{}
}

type LogCommunicator struct {
}

func (mc *LogCommunicator) ReadMessage() (int, []byte, error) {
	return 0, []byte{}, nil
}

func (mc *LogCommunicator) WriteMessage(messageType int, data []byte) error {
	fmt.Println(string(data[:]))
	return nil
}

func (mc *LogCommunicator) Close() error {
	return nil
}
