package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260525000001CreateRepairOrdersTable struct{}

func (m *M20260525000001CreateRepairOrdersTable) Signature() string {
	return "20260525000001_create_repair_orders_table"
}

func (m *M20260525000001CreateRepairOrdersTable) Up() error {
	return facades.Schema().Create("repair_orders", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("works_order_number", 50)
		table.String("notification_number", 50)
		table.String("external_works_order_number", 50)
		table.String("works_order_description", 500)
		table.String("works_order_type", 50)
		table.String("contract_number", 50)
		table.String("works_order_status", 100)
		table.Boolean("wo_cancelled").Default(false)
		table.String("functional_location", 50)
		table.String("functional_location_description", 200)
		table.String("priority_text", 50)
		table.Date("issue_date")
		table.Date("start_date")
		table.Date("finish_date")
		table.Date("extended_completion_date")
		table.String("estimated_total_costs", 20)
		table.String("project_officer_post_id", 50)
		table.String("project_officer_post_description", 200)
		table.Date("report_completion_date")
		table.Date("actual_end_date")
		table.Date("tmc_call_date")
		table.Time("tmc_call_time")
		table.Date("tmc_arrival_date")
		table.Time("tmc_arrival_time")
		table.Date("breakdown_start_date")
		table.Time("breakdown_start_time")
		table.Date("resume_date")
		table.Time("resume_time")
		table.Timestamps()
		table.SoftDeletes()

		table.Unique("works_order_number")
		table.Index("works_order_number")
		table.Index("notification_number")
		table.Index("contract_number")
	})
}

func (m *M20260525000001CreateRepairOrdersTable) Down() error {
	return facades.Schema().DropIfExists("repair_orders")
}
