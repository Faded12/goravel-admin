package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260624000002CreateYiOrdersTable struct{}

func (m *M20260624000002CreateYiOrdersTable) Signature() string {
	return "20260624000002_create_yi_orders_table"
}

func (m *M20260624000002CreateYiOrdersTable) Up() error {
	if facades.Schema().HasTable("yi_orders") {
		return nil
	}
	return facades.Schema().Create("yi_orders", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("yi_no", 50).Comment("一字头订单编号")
		table.String("jiu_order_number", 50).Nullable().Comment("九字头订单编号（关联收款单，非必填）")
		table.String("contract_no", 50).Comment("合同编号")
		table.String("dmo", 50).Comment("DMO编号")
		table.String("estate", 100).Comment("屋邨名称")
		table.String("worker", 100).Comment("工人")
		table.String("project_type", 50).Comment("项目类型")
		table.Boolean("is_complete").Default(false).Comment("是否完成")
		table.Timestamps()
		table.SoftDeletes()

		table.Index("yi_no")
		table.Index("jiu_order_number")
		table.Index("contract_no")
	})
}

func (m *M20260624000002CreateYiOrdersTable) Down() error {
	return facades.Schema().DropIfExists("yi_orders")
}