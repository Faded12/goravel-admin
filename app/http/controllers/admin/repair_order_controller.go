package admin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/http/helpers"
	"goravel/app/http/response"
	"goravel/app/models"
	"goravel/app/services"
)

type RepairOrderController struct {
	repairOrderService services.RepairOrderService
}

func NewRepairOrderController() *RepairOrderController {
	return &RepairOrderController{
		repairOrderService: services.NewRepairOrderService(),
	}
}

func (r *RepairOrderController) Index(ctx http.Context) http.Response {
	page := helpers.GetIntQuery(ctx, "page", 1)
	pageSize := helpers.GetIntQuery(ctx, "page_size", 10)

	filters := services.RepairOrderFilters{
		WorksOrderNumber:   ctx.Request().Query("works_order_number", ""),
		NotificationNumber: ctx.Request().Query("notification_number", ""),
		ContractNumber:     ctx.Request().Query("contract_number", ""),
		WorksOrderStatus:   ctx.Request().Query("works_order_status", ""),
		StartDate:          ctx.Request().Query("start_date", ""),
		EndDate:            ctx.Request().Query("end_date", ""),
	}

	orders, total, err := r.repairOrderService.GetList(filters, page, pageSize)
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, http.Json{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (r *RepairOrderController) Show(ctx http.Context) http.Response {
	id := helpers.GetUintRoute(ctx, "id")
	order, err := r.repairOrderService.GetByID(id)
	if err != nil {
		return response.Error(ctx, http.StatusNotFound, err.Error())
	}

	return response.Success(ctx, http.Json{
		"order": order,
	})
}

func (r *RepairOrderController) Store(ctx http.Context) http.Response {
	return response.Success(ctx, "create_success", http.Json{})
}

func (r *RepairOrderController) Update(ctx http.Context) http.Response {
	return response.Success(ctx, "update_success", http.Json{})
}

func (r *RepairOrderController) Destroy(ctx http.Context) http.Response {
	return response.Success(ctx, "delete_success", http.Json{})
}

func (r *RepairOrderController) Export(ctx http.Context) http.Response {
	adminID, err := helpers.GetAdminIDFromContext(ctx)
	if err != nil {
		return response.Error(ctx, http.StatusUnauthorized, "unauthorized")
	}

	filters := services.RepairOrderFilters{
		WorksOrderNumber:   ctx.Request().Query("works_order_number", ""),
		NotificationNumber: ctx.Request().Query("notification_number", ""),
		ContractNumber:     ctx.Request().Query("contract_number", ""),
		WorksOrderStatus:   ctx.Request().Query("works_order_status", ""),
		StartDate:          ctx.Request().Query("start_date", ""),
		EndDate:            ctx.Request().Query("end_date", ""),
	}

	orders, _, err := r.repairOrderService.GetList(filters, 1, 10000)
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	exportID := fmt.Sprintf("RO%d%d", time.Now().Unix(), adminID)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "工单编号", "通知编号", "工单描述", "工单类型", "合同编号", "工单状态", "功能位置", "功能位置描述", "开始日期", "完成日期", "预估总成本", "项目官员岗位描述"}
	writer.Write(header)

	for _, order := range orders {
		row := []string{
			fmt.Sprintf("%d", order.ID),
			order.WorksOrderNumber,
			order.NotificationNumber,
			order.WorksOrderDescription,
			order.WorksOrderType,
			order.ContractNumber,
			order.WorksOrderStatus,
			order.FunctionalLocation,
			order.FunctionalLocationDescription,
			order.StartDate,
			order.FinishDate,
			order.EstimatedTotalCosts,
			order.ProjectOfficerPostDescription,
		}
		writer.Write(row)
	}
	writer.Flush()

	storage := facades.Storage().Disk("local")
	filename := fmt.Sprintf("repair_orders_%d.csv", time.Now().Unix())
	if err := storage.Put("exports/"+filename, buf.String()); err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "save_file_failed")
	}

	return response.Success(ctx, http.Json{
		"export_id": exportID,
		"message":   "导出成功",
	})
}

func (r *RepairOrderController) Import(ctx http.Context) http.Response {
	_, err := helpers.GetAdminIDFromContext(ctx)
	if err != nil {
		return response.Error(ctx, http.StatusUnauthorized, "unauthorized")
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "file_required")
	}

	filename := file.GetClientOriginalName()
	if !strings.HasSuffix(strings.ToLower(filename), ".csv") {
		return response.Error(ctx, http.StatusBadRequest, "invalid_file_type")
	}

	storage := facades.Storage().Disk("local")
	savedPath, err := storage.PutFile("", file)
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "save_file_failed")
	}

	csvContent, err := storage.Get(savedPath)
	if err != nil {
		_ = storage.Delete(savedPath)
		return response.Error(ctx, http.StatusInternalServerError, "read_file_failed")
	}

	reader := csv.NewReader(strings.NewReader(string(csvContent)))
	records, err := reader.ReadAll()
	if err != nil {
		_ = storage.Delete(savedPath)
		return response.Error(ctx, http.StatusInternalServerError, "parse_csv_failed")
	}

	successCount := 0
	failedCount := 0
	var errors []string

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 12 {
			failedCount++
			errors = append(errors, fmt.Sprintf("row %d: 数据不完整", i))
			continue
		}

		order := models.RepairOrder{
			WorksOrderNumber:                record[1],
			NotificationNumber:              record[2],
			WorksOrderDescription:           record[3],
			WorksOrderType:                  record[4],
			ContractNumber:                  record[5],
			WorksOrderStatus:                record[6],
			FunctionalLocation:              record[7],
			FunctionalLocationDescription:   record[8],
			StartDate:                       record[9],
			FinishDate:                      record[10],
			EstimatedTotalCosts:             record[11],
			ProjectOfficerPostDescription:   record[12],
		}

		if err := facades.Orm().Query().Create(&order); err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("row %d: %s", i, err.Error()))
		} else {
			successCount++
		}
	}

	_ = storage.Delete(savedPath)

	return response.Success(ctx, http.Json{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	})
}
