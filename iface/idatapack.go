package iface

type IDataPacker interface {
	GetHeadLen() uint32

	Pack(IMessage) ([]byte, error)

	Unpack([]byte) (IMessage, error)
}
