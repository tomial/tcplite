package iface

// 封装请求的数据内容
type IMessage interface {
	// 获取消息ID
	GetMsgId() uint32

	// 获取消息的长度
	GetMsgLen() uint32

	// 获取消息内容
	GetData() []byte

	// 设置消息ID
	SetMsgId(uint32)

	// 设置消息长度
	SetMsgLen(uint32)

	// 设置消息内容
	SetData([]byte)
}
