package models

import (
	"github.com/goravel/framework/database/orm"
)

type RepairOrder struct {
	orm.Model
	WorksOrderNumber                string  `gorm:"not null;size:50;uniqueIndex;comment:工单编号" json:"works_order_number"`
	NotificationNumber              string  `gorm:"size:50;comment:通知编号" json:"notification_number"`
	ExternalWorksOrderNumber         string  `gorm:"size:50;comment:外部工单编号" json:"external_works_order_number"`
	WorksOrderDescription            string  `gorm:"size:500;comment:工单描述" json:"works_order_description"`
	WorksOrderType                  string  `gorm:"size:50;comment:工单类型" json:"works_order_type"`
	ContractNumber                  string  `gorm:"size:50;comment:合同编号" json:"contract_number"`
	WorksOrderStatus                string  `gorm:"size:100;comment:工单状态" json:"works_order_status"`
	WoCancelled                     bool    `gorm:"default:false;comment:是否取消" json:"wo_cancelled"`
	FunctionalLocation              string  `gorm:"size:50;comment:功能位置" json:"functional_location"`
	FunctionalLocationDescription   string  `gorm:"size:200;comment:功能位置描述" json:"functional_location_description"`
	PriorityText                   string  `gorm:"size:50;comment:优先级" json:"priority_text"`
	IssueDate                       string  `gorm:"size:10;comment:发布日期" json:"issue_date"`
	StartDate                       string  `gorm:"size:10;comment:开始日期" json:"start_date"`
	FinishDate                      string  `gorm:"size:10;comment:完成日期" json:"finish_date"`
	ExtendedCompletionDate          string  `gorm:"size:10;comment:延展完成日期" json:"extended_completion_date"`
	EstimatedTotalCosts             string  `gorm:"size:20;comment:预估总成本" json:"estimated_total_costs"`
	ProjectOfficerPostID            string  `gorm:"size:50;comment:项目官员岗位ID" json:"project_officer_post_id"`
	ProjectOfficerPostDescription   string  `gorm:"size:200;comment:项目官员岗位描述" json:"project_officer_post_description"`
	ReportCompletionDate            string  `gorm:"size:10;comment:报告完成日期" json:"report_completion_date"`
	ActualEndDate                  string  `gorm:"size:10;comment:实际结束日期" json:"actual_end_date"`
	TmcCallDate                    string  `gorm:"size:10;comment:TMC呼叫日期" json:"tmc_call_date"`
	TmcCallTime                    string  `gorm:"size:8;comment:TMC呼叫时间" json:"tmc_call_time"`
	TmcArrivalDate                 string  `gorm:"size:10;comment:TMC到达日期" json:"tmc_arrival_date"`
	TmcArrivalTime                 string  `gorm:"size:8;comment:TMC到达时间" json:"tmc_arrival_time"`
	BreakdownStartDate             string  `gorm:"size:10;comment:故障开始日期" json:"breakdown_start_date"`
	BreakdownStartTime             string  `gorm:"size:8;comment:故障开始时间" json:"breakdown_start_time"`
	ResumeDate                     string  `gorm:"size:10;comment:恢复日期" json:"resume_date"`
	ResumeTime                     string  `gorm:"size:8;comment:恢复时间" json:"resume_time"`
	orm.SoftDeletes
}

func (RepairOrder) TableName() string {
	return "repair_orders"
}
