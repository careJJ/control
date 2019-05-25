package net

type Message struct {
	Id uint32
	Datalen uint32
	Data []byte
}


func (m *Message)GetMsgId() uint32{
	return m.Id
}


func (m *Message)GetMsgLen() uint32{
	return m.Datalen
}
func (m *Message)GetMsgData() []byte{
	return m.Data
}

//setter
func (m *Message)SetMsgId(id uint32){
	m.Id=id

}
func (m *Message)SetData(data []byte){
	m.Data=data
}
func (m *Message)SetDatalen(len uint32){
	m.Datalen=len
}