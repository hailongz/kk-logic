package lib

import (
	"github.com/hailongz/kk-logic/logic"
)

type OutletLogic struct {
	logic.Logic
}

func (L *OutletLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	ctx.Set(logic.OutletKeys, L)

	return nil
}

func init() {
	logic.Openlib("kk.Logic.Outlet", func(object interface{}) logic.ILogic {
		v := OutletLogic{}
		v.Init(object)
		return &v
	})
}
