package main

import "github.com/bolt"

//定义区块链结构(使用数组模拟区块链)
type BlockChain struct {
	Blocks []*Block //区块链
}

//创世语
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const blockchainDBFile ="blockchain.db"
const bucketBlock ="bucketBlock"
const lastBlockHashKey = "lastBlockHashKey"

//创建区块，从无到有：这个函数仅仅执行一次
func CreateBlockChain()error{
	// 1. 区块链不存在，创建
	db,err:=bolt.Open(blockchainDBFile,0600,nil)
	if err!=nil{
		return err

	}
	//不要db.Close，后续要使用这个句柄
	defer db.Close()
	// 2. 开始创建
	err=db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(bucketBlock))
		//如果bucket为空，说明不存在
		if bucket==nil {
			//创建bucket
			bucket,err:=tx.CreateBucket([]byte(bucketBlock))
			if err!=nil{
					return err
			}
			//写入创世块
			//创建BlockChain，同时添加一个创世块
		genesisBlock:=NewBlock(genesisInfo,nil)
			//key是区块的哈希值，value是block的字节流
		bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
			//更新最后区块哈希值到数据库
		bucket.Put([]byte(lastBlockHashKey),genesisBlock.Hash)
		}
		return nil

	})
		return err
}

//获取区块链实例，用于后续的操作，每一次有业务的时候都会调用
func GetBlockChainInstance()(*BlockChain,error) {
	var lastHash []byte	//内存中最后一个区块的哈希值
	//如果两个区块连不存在，则创建，同时返回blockchain的实例
	bolt.Open(blockchainDBFile,0400,nil)
}





//提供一个创建区块链的方法
func NewBlockChain() *BlockChain {
	//创建BlockChain，同时添加一个创世块
	genesisBlock := NewBlock(genesisInfo, nil)

	bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

	return &bc
}

//提供一个向区块链中添加区块的方法
//参数：数据，不需要提供前区块的哈希值，因为bc可以通过自己的下标拿到
func (bc *BlockChain) AddBlock(data string) {
	//通过下标，得到最后一个区块
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	//最后一个区块哈希值是新区块的前哈希
	prevHash := lastBlock.Hash

	//创建block
	newBlcok := NewBlock(data, prevHash)

	//添加bc中
	bc.Blocks = append(bc.Blocks, newBlcok)
}
