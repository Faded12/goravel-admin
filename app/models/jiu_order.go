package models

import (
	"github.com/goravel/framework/database/orm"
)

type JiuOrder struct {
	orm.Model
	JiuNo          string  `gorm:"column:jiu_no;type:varchar(50);comment:九字头订单编号" json:"jiu_no"`
	ContractNo     string  `gorm:"column:contract_no;type:varchar(50);comment:合同编号" json:"contract_no"`
	Dmo            string  `gorm:"column:dmo;type:varchar(50);comment:DMO编号" json:"dmo"`
	Estate         string  `gorm:"column:estate;type:varchar(100);comment:屋邨名称" json:"estate"`
	Worker         string  `gorm:"column:worker;type:varchar(100);comment:工人" json:"worker"`
	Type           string  `gorm:"column:type;type:varchar(50);comment:类型" json:"type"`
	Fee            float64 `gorm:"column:fee;type:decimal(10,2);comment:费用" json:"fee"`
	OrderStatus    string  `gorm:"column:order_status;type:varchar(50);comment:订单状态" json:"order_status"`
	CompleteStatus string  `gorm:"column:complete_status;type:varchar(50);comment:完成状态" json:"complete_status"`
	ProjectType    string  `gorm:"column:project_type;type:varchar(50);comment:项目类型" json:"project_type"`
	DateStart      string  `gorm:"column:date_start;type:date;comment:开始日期" json:"date_start"`
	DateEnd        string  `gorm:"column:date_end;type:date;comment:结束日期" json:"date_end"`
}

func (JiuOrder) TableName() string {
	return "jiu_orders"
}