package zinface

type IMessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	SetMsgId(uint32)
	SetDatalen(uint32)
	SetData([]byte)


}