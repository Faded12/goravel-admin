package admin

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"

	"goravel/app/http/helpers"
	"goravel/app/http/response"
	"goravel/app/models"
)

type JiuOrderController struct{}

func NewJiuOrderController() *JiuOrderController {
	return &JiuOrderController{}
}

// Index lists JiuOrder records.
func (c *JiuOrderController) Index(ctx http.Context) http.Response {
	page := helpers.GetIntQuery(ctx, "page", 1)
	pageSize := helpers.GetIntQuery(ctx, "page_size", 10)

	query := facades.Orm().Query().Model(&models.JiuOrder{}).Where("deleted_at", nil)

	// 筛选条件
	if jiuNo := ctx.Request().Query("jiu_no", ""); jiuNo != "" {
		query = query.Where("jiu_no", "LIKE", "%"+jiuNo+"%")
	}
	if contractNo := ctx.Request().Query("contract_no", ""); contractNo != "" {
		query = query.Where("contract_no", "LIKE", "%"+contractNo+"%")
	}
	if dmo := ctx.Request().Query("dmo", ""); dmo != "" {
		query = query.Where("dmo", dmo)
	}
	if estate := ctx.Request().Query("estate", ""); estate != "" {
		query = query.Where("estate", estate)
	}
	if worker := ctx.Request().Query("worker", ""); worker != "" {
		query = query.Where("worker", "LIKE", "%"+worker+"%")
	}
	if orderStatus := ctx.Request().Query("order_status", ""); orderStatus != "" {
		query = query.Where("order_status", orderStatus)
	}
	if completeStatus := ctx.Request().Query("complete_status", ""); completeStatus != "" {
		query = query.Where("complete_status", completeStatus)
	}
	if projectType := ctx.Request().Query("project_type", ""); projectType != "" {
		query = query.Where("project_type", projectType)
	}

	// 日期范围筛选
	if dateStart := ctx.Request().Query("date_start", ""); dateStart != "" {
		query = query.Where("date_start", ">=", dateStart)
	}
	if dateEnd := ctx.Request().Query("date_end", ""); dateEnd != "" {
		query = query.Where("date_end", "<=", dateEnd)
	}

	total, err := query.Count()
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "query_failed")
	}

	var jiuOrders []models.JiuOrder
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Get(&jiuOrders); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "query_failed")
	}

	return response.Success(ctx, http.Json{
		"list":      jiuOrders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Show returns JiuOrder details.
func (c *JiuOrderController) Show(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")

	var jiuOrder models.JiuOrder
	if err := facades.Orm().Query().Where("id", id).Where("deleted_at", nil).First(&jiuOrder); err != nil {
		return response.Error(ctx, http.StatusNotFound, "record_not_found")
	}

	return response.Success(ctx, http.Json{
		"jiu_order": jiuOrder,
	})
}

// Store creates a new JiuOrder.
func (c *JiuOrderController) Store(ctx http.Context) http.Response {
	jiuOrder := models.JiuOrder{
		JiuNo:          ctx.Request().Input("jiu_no", ""),
		ContractNo:     ctx.Request().Input("contract_no", ""),
		Dmo:            ctx.Request().Input("dmo", ""),
		Estate:         ctx.Request().Input("estate", ""),
		Worker:         ctx.Request().Input("worker", ""),
		Type:           ctx.Request().Input("type", ""),
		Fee:            cast.ToFloat64(ctx.Request().Input("fee", "0")),
		OrderStatus:    ctx.Request().Input("order_status", ""),
		CompleteStatus: ctx.Request().Input("complete_status", ""),
		ProjectType:    ctx.Request().Input("project_type", ""),
		DateStart:      ctx.Request().Input("date_start", ""),
		DateEnd:        ctx.Request().Input("date_end", ""),
	}

	if err := facades.Orm().Query().Create(&jiuOrder); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "create_failed")
	}

	return response.Success(ctx, http.Json{
		"jiu_order": jiuOrder,
	})
}

// Update modifies an existing JiuOrder.
func (c *JiuOrderController) Update(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")

	var jiuOrder models.JiuOrder
	if err := facades.Orm().Query().Where("id", id).Where("deleted_at", nil).First(&jiuOrder); err != nil {
		return response.Error(ctx, http.StatusNotFound, "record_not_found")
	}

	jiuOrder.JiuNo = ctx.Request().Input("jiu_no", jiuOrder.JiuNo)
	jiuOrder.ContractNo = ctx.Request().Input("contract_no", jiuOrder.ContractNo)
	jiuOrder.Dmo = ctx.Request().Input("dmo", jiuOrder.Dmo)
	jiuOrder.Estate = ctx.Request().Input("estate", jiuOrder.Estate)
	jiuOrder.Worker = ctx.Request().Input("worker", jiuOrder.Worker)
	jiuOrder.Type = ctx.Request().Input("type", jiuOrder.Type)
	if fee := cast.ToFloat64(ctx.Request().Input("fee", "0")); fee > 0 {
		jiuOrder.Fee = fee
	}
	jiuOrder.OrderStatus = ctx.Request().Input("order_status", jiuOrder.OrderStatus)
	jiuOrder.CompleteStatus = ctx.Request().Input("complete_status", jiuOrder.CompleteStatus)
	jiuOrder.ProjectType = ctx.Request().Input("project_type", jiuOrder.ProjectType)
	jiuOrder.DateStart = ctx.Request().Input("date_start", jiuOrder.DateStart)
	jiuOrder.DateEnd = ctx.Request().Input("date_end", jiuOrder.DateEnd)

	if err := facades.Orm().Query().Save(&jiuOrder); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "update_failed")
	}

	return response.Success(ctx, http.Json{
		"jiu_order": jiuOrder,
	})
}

// Destroy deletes a JiuOrder (soft delete).
func (c *JiuOrderController) Destroy(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")

	if _, err := facades.Orm().Query().Model(&models.JiuOrder{}).Where("id", id).Delete(); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "delete_failed")
	}

	return response.Success(ctx, "delete_success", http.Json{})
}

// Export exports JiuOrder records.
func (c *JiuOrderController) Export(ctx http.Context) http.Response {
	query := facades.Orm().Query().Model(&models.JiuOrder{}).Where("deleted_at", nil)

	// 筛选条件
	if jiuNo := ctx.Request().Input("jiu_no", ctx.Request().Query("jiu_no", "")); jiuNo != "" {
		query = query.Where("jiu_no", "LIKE", "%"+jiuNo+"%")
	}
	if contractNo := ctx.Request().Input("contract_no", ctx.Request().Query("contract_no", "")); contractNo != "" {
		query = query.Where("contract_no", "LIKE", "%"+contractNo+"%")
	}
	if dmo := ctx.Request().Input("dmo", ctx.Request().Query("dmo", "")); dmo != "" {
		query = query.Where("dmo", dmo)
	}
	if estate := ctx.Request().Input("estate", ctx.Request().Query("estate", "")); estate != "" {
		query = query.Where("estate", estate)
	}

	var jiuOrders []models.JiuOrder
	if err := query.Order("created_at DESC").Get(&jiuOrders); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "export_failed")
	}

	var data [][]string
	header := []string{"ID", "九字头订单编号", "合同编号", "DMO编号", "屋邨名称", "工人", "类型", "费用", "订单状态", "完成状态", "项目类型", "开始日期", "结束日期"}
	data = append(data, header)

	for _, order := range jiuOrders {
		row := []string{
			fmt.Sprintf("%d", order.ID),
			order.JiuNo,
			order.ContractNo,
			order.Dmo,
			order.Estate,
			order.Worker,
			order.Type,
			fmt.Sprintf("%.2f", order.Fee),
			order.OrderStatus,
			order.CompleteStatus,
			order.ProjectType,
			order.DateStart,
			order.DateEnd,
		}
		data = append(data, row)
	}

	return response.Success(ctx, http.Json{
		"data":    data,
		"message": "export_success",
	})
}