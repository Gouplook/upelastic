/********************************************

@Author :yinjinlin<yinjinlin_uplook@163.com>
@Time : 2021/2/22 11:05
@Description:

*********************************************/
package bus

import (
	"context"
	"upelastic/rpcinterface/interface/elastic/bus"
)

type Bus struct {
	//client.Baseclient
}

func (bus *Bus)Init() *Bus{
	//bus.ServiceName = "rpc_upelastic"
	//bus.ServicePaht = "Bus/Bus"
	return bus
}

func (bus *Bus)SearchAdminBus(ctx context.Context, args *bus.ArgsBusElastic, reply *bus.ReplyBusElastic) error {
	//return bus.Call(ctx, "SearchAminBus", args,reply)
	return nil
}
