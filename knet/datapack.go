package knet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/k-si/Kinx/kiface"
)

const (
	MessageHeadLength = 8
)

// 统一根据TLV协议，处理数据封装和拆解，以此处理tcp黏包问题
type DataPack struct{}

// 将message转为二进制切片
func (d *DataPack) Pack(msg kiface.IMessage) ([]byte, error) {

	// 创建二进制buf，存放message
	buf := bytes.NewBuffer([]byte{})

	// 将msg中的id、len、data放入buf
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 将二进制数据中的head抽离出来
func (d *DataPack) UnPack(data []byte) (kiface.IMessage, error) {

	// 将二进制数据存入buf
	buf := bytes.NewBuffer(data)
	head := &Message{}

	// 将二进制数据读取出来，存入结构体中
	if err := binary.Read(buf, binary.LittleEndian, &head.id); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &head.length); err != nil {
		return nil, err
	}

	// 如果该message的data过长，返回err
	if head.GetMsgLen() > config.MaxPackageSize {
		return nil, errors.New("message length larger than max package length")
	}

	return head, nil
}

func NewDataPack() kiface.IDataPack {
	return &DataPack{}
}
