package main

import (
	"github.com/bolt"
	"errors"
)

//定义区块链结构(使用数组模拟区块链)
type BlockChain struct {
	db   *bolt.DB //用于存储数据
	//Blocks []*Block //区块链
	tail []byte   //最后一个区块的哈希值
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


//获取区块链实例，用于后续操作, 每一次有业务时都会调用
func GetBlockChainInstance() (*BlockChain, error) {
	var lastHash []byte //内存中最后一个区块的哈希值

	//两个功能：
	// 1. 如果区块链不存在，则创建，同时返回blockchain的示例
	db, err := bolt.Open(blockchainDBFile, 0400, nil) //rwx  0100 => 4
	if err != nil {
		return nil, err
	}

	//不要db.Close，后续要使用这个句柄

	// 2. 如果区块链存在，则直接返回blockchain示例
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			return errors.New("bucket不应为nil")
		} else {
			//直接读取特定的key，得到最后一个区块的哈希值
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	//5. 拼成BlockChain然后返回
	bc := BlockChain{db, lastHash}
	return &bc, nil
}

//提供一个向区块链中添加区块的方法
func (bc *BlockChain) AddBlock(data string) error {
	lashBlockHash := bc.tail //区块链中最后一个区块的哈希值

	//1. 创建区块
	newBlock := NewBlock(data, lashBlockHash)

	//2. 写入数据库
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("AddBlock时Bucket不应为空")
		}

		//key是新区块的哈希值， value是这个区块的字节流
		bucket.Put(newBlock.Hash, newBlock.Serialize())
		bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)

		//更新bc的tail，这样后续的AddBlock才会基于我们newBlock追加
		bc.tail = newBlock.Hash
		return nil
	})

	return err
}




////提供一个创建区块链的方法
//func NewBlockChain() *BlockChain {
//	//创建BlockChain，同时添加一个创世块
//	genesisBlock := NewBlock(genesisInfo, nil)
//
//	bc := BlockChain{
//		Blocks: []*Block{genesisBlock},
//	}
//
//	return &bc
//}

//提供一个向区块链中添加区块的方法
//参数：数据，不需要提供前区块的哈希值，因为bc可以通过自己的下标拿到
//func (bc *BlockChain) AddBlock(data string) {
//	//通过下标，得到最后一个区块
//	lastBlock := bc.Blocks[len(bc.Blocks)-1]
//
//	//最后一个区块哈希值是新区块的前哈希
//	prevHash := lastBlock.Hash
//
//	//创建block
//	newBlcok := NewBlock(data, prevHash)
//
//	//添加bc中
//	bc.Blocks = append(bc.Blocks, newBlcok)
//}
