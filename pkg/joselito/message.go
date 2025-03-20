package joselito

type MessageType int

const (
	GROUP_JOIN    MessageType = 1
	GROUP_LEAVE   MessageType = 2
	GROUP_RESET   MessageType = 3
	CALL_START    MessageType = 11 // 0x0B
	CALL_DROP     MessageType = 12 // 0x0C
	CALL_END      MessageType = 13 // 0x0D
	CALL_AUDIO    MessageType = 20 // 0x14
	CALL_ALIAS    MessageType = 21 // 0x15
	CALL_METER    MessageType = 22 // 0x16
	SYSTEM_RESCUE MessageType = 80 // 0x50
)

type Message interface {
	MessageType() MessageType
	Marshall() ([]byte, error)
	Unmarshall() error
}
