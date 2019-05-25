package net

import (
	"zin/zinface"
	"bytes"
	"encoding/binary"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取二进制包头部长度，固定返回8字节
func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}
//封包 将Message打包
func (dp *DataPack) Pack(msg zinface.IMessage) ([]byte,error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuffer,binary.LittleEndian,msg.GetMsgId());err!=nil {
		return nil,err
	}
	if err:=binary.Write(dataBuffer,binary.LittleEndian,msg.GetMsgData());err!=nil {
		return nil,err
	}
	//返回这个缓冲
	return dataBuffer.Bytes(),nil

}
//拆包方法
func (dp *DataPack) UnPack(binaryData []byte) (zinface.IMessage,error){
	//拆包分两次解压，第一次读取的固定的长度8字节，第二次根据len再次进行read
	msgHead:=&Message{}  //msgHead.Datalen, msgHead.dataID
	//创建一个读取二进制数据流的io.Reader
	dataBuff:=bytes.NewReader(binaryData)

	//将二进制流的 先读datalen 放在msg的DataLen属性中
	if err:=binary.Read(dataBuff,binary.LittleEndian,&msgHead.Datalen);err!=nil{
		return nil,err
	}
	if err:=binary.Read(dataBuff,binary.LittleEndian,&msgHead.Id);err!=nil{
		return  nil,err
	}

	return msgHead,nil

}