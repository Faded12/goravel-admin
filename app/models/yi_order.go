package models

import (
	"github.com/goravel/framework/database/orm"
)

type YiOrder struct {
	orm.Model
	YiOrderNumber string `gorm:"column:yi_no;type:varchar(50);comment:一字头订单编号" json:"yi_no"`
	JiuOrderNumber string `gorm:"column:jiu_order_number;type:varchar(50);comment:九字头订单编号" json:"jiu_order_number"`
	ContractNo    string `gorm:"column:contract_no;type:varchar(50);comment:合同编号" json:"contract_no"`
	Dmo           string `gorm:"column:dmo;type:varchar(50);comment:DMO编号" json:"dmo"`
	Estate        string `gorm:"column:estate;type:varchar(100);comment:屋邨名称" json:"estate"`
	Worker        string `gorm:"column:worker;type:varchar(100);comment:工人" json:"worker"`
	ProjectType   string `gorm:"column:project_type;type:varchar(50);comment:项目类型" json:"project_type"`
	IsComplete    bool   `gorm:"column:is_complete;type:tinyint(1);default:0;comment:是否完成" json:"is_complete"`
}

func (YiOrder) TableName() string {
	return "yi_orders"
}