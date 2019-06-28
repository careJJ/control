package main

import (
	"math/big"
	"fmt"
	"bytes"
	"crypto/sha256"
)

type ProofOfWork struct{
	// ​区块：block
	block *Block
	// ​目标值：target，这个目标值要与生成哈希值比较
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork{
	pow:=ProofOfWork{
		block:block,
	}
	//难度值先写死，不去推导，后面补充推导方式
	//64位的16进制数: 64 * 4 = 256
	targetStr := "0000010000000000000000000000000000000000000000000000000000000000"
	tmpBigInt := new(big.Int)
	//将我们的难度值赋值给bigint
	tmpBigInt.SetString(targetStr,16)
	pow.target=tmpBigInt
	return &pow
}

//挖矿函数，不断变化nonce，使得sha256(数据+nonce)<难度值
//返回：区块哈希，nonce
func (pow *ProofOfWork)Run()([]byte,uint64)  {
	//定义随机数
	var nonce uint64
	var hash  [32]byte
	fmt.Println("开始挖矿......")
	for{
		fmt.Printf("%x\r",hash[:])
		// 1. 拼接字符串 + nonce
		data := pow.PrepareData(nonce)
		//2.哈希值=sha256(data)
		hash=sha256.Sum256(data)
		//将hash转换成bigInt类型
		tmpInt:=new(big.Int)
		tmpInt.SetBytes(hash[:])
		//当前计算的哈希.Cmp(难度值)
		if tmpInt.Cmp(pow.target)==-1{
			fmt.Printf("挖矿成功,hash :%x, nonce :%d\n", hash[:], nonce)
			break
		}else{
			//如果不小于难度值
			nonce++
		}

	}
	return hash[:],nonce


}


func (pow *ProofOfWork)PrepareData(nonce uint64)[]byte{
	b:=pow.block
	tmp:=[][]byte{
		uintToByte(b.Version), //将uint64转换为[]byte
		b.PrevHash,
		b.MerkleRoot,
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(nonce),
		// b.Hash, //它不应该参与哈希运算
		b.Data,

	}
	//使用join方法，将二维切片转为1维切片
	data:=bytes.Join(tmp,[]byte{})
	return data

}

func(pow *ProofOfWork)IsValid()bool{
	// 	1. 获取区块
	// 2. 拼装数据（block + nonce）
	data:=pow.PrepareData(pow.block.Nonce)
	//3.计算sha256
	hash:=sha256.Sum256(data)
	//与难度值作对比
	tmpInt:=new(big.Int)
	tmpInt.SetBytes(hash[:])
	return tmpInt.Cmp(pow.target)==-1



}