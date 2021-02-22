/********************************************

@Author :yinjinlin<yinjinlin_uplook@163.com>
@Time : 2021/2/22 11:13
@Description:

*********************************************/
package bus

import "upelastic/rpcinterface/interface/elastic/bus"

type BusLogic struct {

}

// SearchAminBus ES商户后台搜索
func (b *BusLogic)SearchAminBus(args *bus.ArgsBusElastic, reply *bus.ReplyBusElastic)(){
	// 1： 初始化ElasticClient （NewElasiticClient）

	// 2：字段刷选： 分普通字段和关键字刷选
	// 字段精选刷选 SetFilter

	// 关键字搜索 Search
	// "*"+where.BrandName+"*"

	// 时间范围搜索 SetFilterGte 提交开始和结束时间

	// 3:对于多个搜索结果进行分页和排序
	// SetLimit SetSortMode

	// 4：请求

}
