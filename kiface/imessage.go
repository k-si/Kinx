package kiface

// 对请求数据的封装
type IMessage interface {

	// Getter
	GetMsgUuid() uint64
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	// Setter
	SetMsgUuid(uint64)
	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetMsgData([]byte)
}
