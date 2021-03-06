package net

type Message struct {
	// 消息ID
	Id uint32

	// 消息长度
	DataLen uint32

	// 消息内容
	Data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// 设置消息长度
func (m *Message) SetMsgLen(length uint32) {
	m.DataLen = length
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
