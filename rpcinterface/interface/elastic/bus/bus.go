/********************************************

@Author :yinjinlin<yinjinlin_uplook@163.com>
@Time : 2021/2/22 10:52
@Description:

*********************************************/
package bus

import "context"

type ArgsBusElastic struct {
	CompanyName string //　企业名称
	BrandName   string //　店铺名称
	Pid         int    // 直属省份/城市
	Cid         int    // 所属城市ID
	Status      string // 商户审核状态 0=待审核 1=审核失败 2=已通过审核 3=下架
	CtimeStart  int64  //　提交开始时间戳
	CtimeEnd    int64  //　提交结束时间戳
}

type ReplyBusElastic struct {

}

// ES 搜索
type Bus interface {
	// SetAdminBus 设置商户信息到ES
	SetAdminBus(ctx context.Context,busId *int, reply *bool)
	// SearchAdminBus ES商户后台搜索
	SearchAdminBus(ctx context.Context, args *ArgsBusElastic, reply *ReplyBusElastic) error
}