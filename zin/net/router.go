package net

import "zin/zinface"

type Baserouter struct {


}



//处理业务之前的方法
func (r *Baserouter)PreHandle(request zinface.IRequest) {
	//将interface的方法全部实现， 目的是 让用户重写这个方法
}
//真正处理业务的方法
func (r *Baserouter)Handle(request zinface.IRequest){

}
//处理业务之后的方法
func (r *Baserouter)PostHandle(request zinface.IRequest){

}