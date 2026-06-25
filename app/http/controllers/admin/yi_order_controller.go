package admin

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/http/helpers"
	"goravel/app/http/response"
	"goravel/app/models"
)

type YiOrderController struct{}

func NewYiOrderController() *YiOrderController {
	return &YiOrderController{}
}

// Index lists YiOrder records.
func (c *YiOrderController) Index(ctx http.Context) http.Response {
	page := helpers.GetIntQuery(ctx, "page", 1)
	pageSize := helpers.GetIntQuery(ctx, "page_size", 10)

	query := facades.Orm().Query().Model(&models.YiOrder{}).Where("deleted_at", nil)

	// 筛选条件
	if yiNo := ctx.Request().Query("yi_no", ""); yiNo != "" {
		query = query.Where("yi_no", "LIKE", "%"+yiNo+"%")
	}
	if jiuOrderNumber := ctx.Request().Query("jiu_order_number", ""); jiuOrderNumber != "" {
		query = query.Where("jiu_order_number", "LIKE", "%"+jiuOrderNumber+"%")
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
	if projectType := ctx.Request().Query("project_type", ""); projectType != "" {
		query = query.Where("project_type", projectType)
	}
	if isComplete := ctx.Request().Query("is_complete", ""); isComplete != "" {
		query = query.Where("is_complete", isComplete == "true")
	}

	total, err := query.Count()
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "query_failed")
	}

	var yiOrders []models.YiOrder
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Get(&yiOrders); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "query_failed")
	}

	return response.Success(ctx, http.Json{
		"list":      yiOrders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Show returns YiOrder details.
func (c *YiOrderController) Show(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")

	var yiOrder models.YiOrder
	if err := facades.Orm().Query().Where("id", id).Where("deleted_at", nil).First(&yiOrder); err != nil {
		return response.Error(ctx, http.StatusNotFound, "record_not_found")
	}

	return response.Success(ctx, http.Json{
		"yi_order": yiOrder,
	})
}

// Store creates a new YiOrder.
func (c *YiOrderController) Store(ctx http.Context) http.Response {
	yiOrder := models.YiOrder{
		YiOrderNumber: ctx.Request().Input("yi_no", ""),
		JiuOrderNumber:  ctx.Request().Input("jiu_order_number", ""),
		ContractNo:      ctx.Request().Input("contract_no", ""),
		Dmo:             ctx.Request().Input("dmo", ""),
		Estate:          ctx.Request().Input("estate", ""),
		Worker:          ctx.Request().Input("worker", ""),
		ProjectType:     ctx.Request().Input("project_type", ""),
		IsComplete:      ctx.Request().Input("is_complete", "false") == "true",
	}

	if err := facades.Orm().Query().Create(&yiOrder); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "create_failed")
	}

	return response.Success(ctx, http.Json{
		"yi_order": yiOrder,
	})
}

// Update modifies an existing YiOrder.
func (c *YiOrderController) Update(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")

	var yiOrder models.YiOrder
	if err := facades.Orm().Query().Where("id", id).Where("deleted_at", nil).First(&yiOrder); err != nil {
		return response.Error(ctx, http.StatusNotFound, "record_not_found")
	}

	yiOrder.YiOrderNumber = ctx.Request().Input("yi_no", yiOrder.YiOrderNumber)
	yiOrder.JiuOrderNumber = ctx.Request().Input("jiu_order_number", yiOrder.JiuOrderNumber)
	yiOrder.ContractNo = ctx.Request().Input("contract_no", yiOrder.ContractNo)
	yiOrder.Dmo = ctx.Request().Input("dmo", yiOrder.Dmo)
	yiOrder.Estate = ctx.Request().Input("estate", yiOrder.Estate)
	yiOrder.Worker = ctx.Request().Input("worker", yiOrder.Worker)
	yiOrder.ProjectType = ctx.Request().Input("project_type", yiOrder.ProjectType)
	if isComplete := ctx.Request().Input("is_complete", ""); isComplete != "" {
		yiOrder.IsComplete = isComplete == "true"
	}

	if err := facades.Orm().Query().Save(&yiOrder); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "update_failed")
	}

	return response.Success(ctx, http.Json{
		"yi_order": yiOrder,
	})
}

// Destroy deletes a YiOrder (soft delete).
func (c *YiOrderController) Destroy(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")

	if _, err := facades.Orm().Query().Model(&models.YiOrder{}).Where("id", id).Delete(); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "delete_failed")
	}

	return response.Success(ctx, "delete_success", http.Json{})
}

// Export exports YiOrder records.
func (c *YiOrderController) Export(ctx http.Context) http.Response {
	query := facades.Orm().Query().Model(&models.YiOrder{}).Where("deleted_at", nil)

	// 筛选条件
	if yiNo := ctx.Request().Input("yi_no", ctx.Request().Query("yi_no", "")); yiNo != "" {
		query = query.Where("yi_no", "LIKE", "%"+yiNo+"%")
	}
	if jiuOrderNumber := ctx.Request().Input("jiu_order_number", ctx.Request().Query("jiu_order_number", "")); jiuOrderNumber != "" {
		query = query.Where("jiu_order_number", "LIKE", "%"+jiuOrderNumber+"%")
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

	var yiOrders []models.YiOrder
	if err := query.Order("created_at DESC").Get(&yiOrders); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "export_failed")
	}

	var data [][]string
	header := []string{"ID", "一字头订单编号", "九字头订单编号", "合同编号", "DMO编号", "屋邨名称", "工人", "项目类型", "是否完成"}
	data = append(data, header)

	for _, order := range yiOrders {
		isComplete := "否"
		if order.IsComplete {
			isComplete = "是"
		}
		row := []string{
			fmt.Sprintf("%d", order.ID),
			order.YiOrderNumber,
			order.JiuOrderNumber,
			order.ContractNo,
			order.Dmo,
			order.Estate,
			order.Worker,
			order.ProjectType,
			isComplete,
		}
		data = append(data, row)
	}

	return response.Success(ctx, http.Json{
		"data":    data,
		"message": "export_success",
	})
}