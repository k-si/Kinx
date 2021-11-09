package knet

import "kinx/kiface"

// TLV协议：MessageId | MessageLen | MessageData
type Message struct {
	id     uint32
	length uint32
	data []byte
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) GetMsgLen() uint32 {
	return m.length
}

func (m *Message) GetMsgData() []byte {
	return m.data
}

func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

func (m *Message) SetMsgLen(length uint32) {
	m.length = length
}

func (m *Message) SetMsgData(data []byte) {
	m.data = data
}

func NewMessage(id uint32, data []byte) kiface.IMessage {
	message := &Message{
		id:     id,
		length: uint32(len(data)),
		data:   data,
	}
	return message
}
