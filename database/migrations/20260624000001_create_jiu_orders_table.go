package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260624000001CreateJiuOrdersTable struct{}

func (m *M20260624000001CreateJiuOrdersTable) Signature() string {
	return "20260624000001_create_jiu_orders_table"
}

func (m *M20260624000001CreateJiuOrdersTable) Up() error {
	if facades.Schema().HasTable("jiu_orders") {
		return nil
	}
	return facades.Schema().Create("jiu_orders", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("jiu_no", 50).Comment("九字头订单编号")
		table.String("contract_no", 50).Comment("合同编号")
		table.String("dmo", 50).Comment("DMO编号")
		table.String("estate", 100).Comment("屋邨名称")
		table.String("worker", 100).Comment("工人")
		table.String("type", 50).Comment("类型")
		table.Decimal("fee").Total(10).Places(2).Comment("费用")
		table.String("order_status", 50).Comment("订单状态")
		table.String("complete_status", 50).Comment("完成状态")
		table.String("project_type", 50).Comment("项目类型")
		table.Date("date_start").Comment("开始日期")
		table.Date("date_end").Comment("结束日期")
		table.Timestamps()
		table.SoftDeletes()

		table.Index("jiu_no")
		table.Index("contract_no")
		table.Index("dmo")
	})
}

func (m *M20260624000001CreateJiuOrdersTable) Down() error {
	return facades.Schema().DropIfExists("jiu_orders")
}