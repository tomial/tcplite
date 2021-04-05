package net

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/tomial/tcplite/iface"
	"github.com/tomial/tcplite/utils"
)

// 对Messagge进行封包、拆包
type DataPacker struct{}

func NewDataPacker() *DataPacker {
	return &DataPacker{}
}

// 获取头部长度
func (d *DataPacker) GetHeadLen() uint32 {
	// Len(uint32) 4bytes + ID(uint32) 4bytes = 8 bytes
	return 8
}

// 封包
func (d *DataPacker) Pack(msg iface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})

	// 写入Msg长度
	err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	// 写入MsgID
	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 写入Msg内容
	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

// 拆包
func (d *DataPacker) Unpack(binaryData []byte) (iface.IMessage, error) {
	dataBuffer := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if msg.DataLen < 0 || msg.DataLen > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Invalid packet size")
	}

	return msg, nil
}
