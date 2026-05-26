package admin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/xuri/excelize/v2"

	"goravel/app/http/helpers"
	"goravel/app/http/response"
	"goravel/app/services"
)

func parseDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}
	if strings.Contains(dateStr, "T") {
		if t, err := time.Parse("2006-01-02T00:00:00Z", dateStr); err == nil {
			return t.Format("2006-01-02")
		}
		if t, err := time.Parse("2006-01-02T00:00:00+00:00", dateStr); err == nil {
			return t.Format("2006-01-02")
		}
		if t, err := time.Parse("2006-01-02T00:00:00", dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	layouts := []string{
		"2006/1/2", "2006-1-2", "2006.1.2",
		"02/01/2006", "01/02/2006",
		"01-02-06", "02-01-06",
		"2/1/2006", "1/2/2006",
		"2-1-2006", "1-2-2006",
		"2/1/06", "1/2/06",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return ""
}

func excelDateToString(excelDate float64) string {
	excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	date := excelEpoch.AddDate(0, 0, int(excelDate))
	return date.Format("2006-01-02")
}

func valuesToPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return "(" + strings.Join(placeholders, ",") + ")"
}

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
	isXlsx := strings.HasSuffix(strings.ToLower(filename), ".xlsx")
	isCsv := strings.HasSuffix(strings.ToLower(filename), ".csv")
	if !isXlsx && !isCsv {
		return response.Error(ctx, http.StatusBadRequest, "invalid_file_type")
	}

	storage := facades.Storage().Disk("local")
	savedPath, err := storage.PutFile("", file)
	if err != nil {
		return response.Error(ctx, http.StatusInternalServerError, "save_file_failed")
	}

	content, err := storage.Get(savedPath)
	if err != nil {
		_ = storage.Delete(savedPath)
		return response.Error(ctx, http.StatusInternalServerError, "read_file_failed")
	}

	var records [][]string

	if isXlsx {
		f, err := excelize.OpenReader(bytes.NewReader([]byte(content)), excelize.Options{})
		if err != nil {
			_ = storage.Delete(savedPath)
			return response.Error(ctx, http.StatusInternalServerError, "parse_xlsx_failed")
		}
		rows, err := f.GetRows("Sheet1")
		f.Close()
		if err != nil {
			_ = storage.Delete(savedPath)
			return response.Error(ctx, http.StatusInternalServerError, "parse_xlsx_failed")
		}
		for _, row := range rows {
			record := make([]string, len(row))
			for j, cell := range row {
				if cell == "" {
					continue
				}
				if strings.Contains(cell, "/") || strings.Contains(cell, "-") || strings.Contains(cell, ".") {
					record[j] = strings.TrimSpace(cell)
				} else if num, err := strconv.ParseFloat(cell, 64); err == nil && num > 0 {
					record[j] = excelDateToString(num)
				} else {
					record[j] = strings.TrimSpace(cell)
				}
			}
			records = append(records, record)
		}
	} else {
		reader := csv.NewReader(strings.NewReader(string(content)))
		records, err = reader.ReadAll()
		if err != nil {
			_ = storage.Delete(savedPath)
			return response.Error(ctx, http.StatusInternalServerError, "parse_csv_failed")
		}
	}

	_ = storage.Delete(savedPath)

	successCount := 0
	failedCount := 0
	var errors []string

	for i, record := range records {
		if i == 0 {
			facades.Log().Info(fmt.Sprintf("Excel header row (%d columns): %v", len(record), record))
			continue
		}

		if len(record) < 13 {
			failedCount++
			errors = append(errors, fmt.Sprintf("row %d: 数据不完整 (columns: %d)", i, len(record)))
			continue
		}

		if i == 1 {
			facades.Log().Info(fmt.Sprintf("First data row (%d columns):", len(record)))
			for idx, val := range record {
				facades.Log().Info(fmt.Sprintf("  record[%d] = %q", idx, val))
			}
		}

		fields := []string{}
		values := []any{}

		addField := func(name string, value string) {
			if value != "" {
				fields = append(fields, name)
				values = append(values, value)
			}
		}

		addDateField := func(name string, dateStr string) {
			parsed := parseDate(dateStr)
			if parsed != "" {
				fields = append(fields, name)
				values = append(values, parsed)
			}
		}

		addTimeField := func(name string, timeStr string) {
			timeStr = strings.TrimSpace(timeStr)
			if timeStr != "" {
				fields = append(fields, name)
				values = append(values, timeStr)
			}
		}

		addBoolField := func(name string, value string) bool {
			value = strings.TrimSpace(value)
			if value != "" {
				boolVal := value == "1" || strings.ToLower(value) == "true" || strings.ToLower(value) == "x" || value == "X"
				fields = append(fields, name)
				values = append(values, boolVal)
				return true
			}
			return false
		}

		addField("works_order_number", record[0])
		addField("notification_number", record[1])
		addField("external_works_order_number", record[2])
		addField("works_order_description", record[3])
		addField("works_order_type", record[4])
		addField("contract_number", record[5])
		addField("works_order_status", record[6])

		if len(record) > 7 { addBoolField("wo_cancelled", record[7]) }

		if len(record) > 8 { addField("functional_location", record[8]) }
		if len(record) > 9 { addField("functional_location_description", record[9]) }
		if len(record) > 10 { addField("priority_text", record[10]) }

		if len(record) > 11 { addDateField("issue_date", record[11]) }
		if len(record) > 12 { addDateField("start_date", record[12]) }
		if len(record) > 13 { addDateField("finish_date", record[13]) }
		if len(record) > 14 { addDateField("extended_completion_date", record[14]) }

		if len(record) > 15 { addField("estimated_total_costs", record[15]) }
		if len(record) > 16 { addField("project_officer_post_id", record[16]) }
		if len(record) > 17 { addField("project_officer_post_description", record[17]) }

		if len(record) > 18 { addDateField("report_completion_date", record[18]) }
		if len(record) > 19 { addDateField("actual_end_date", record[19]) }
		if len(record) > 20 { addDateField("tmc_call_date", record[20]) }
		if len(record) > 21 { addTimeField("tmc_call_time", record[21]) }
		if len(record) > 22 { addDateField("tmc_arrival_date", record[22]) }
		if len(record) > 23 { addTimeField("tmc_arrival_time", record[23]) }
		if len(record) > 24 { addDateField("breakdown_start_date", record[24]) }
		if len(record) > 25 { addTimeField("breakdown_start_time", record[25]) }
		if len(record) > 26 { addDateField("resume_date", record[26]) }
		if len(record) > 27 { addTimeField("resume_time", record[27]) }

		if len(fields) == 0 {
			failedCount++
			errors = append(errors, fmt.Sprintf("row %d: 没有有效数据", i))
			continue
		}

		sql := "INSERT INTO repair_orders (" + strings.Join(fields, ",") + ") VALUES " + valuesToPlaceholders(len(values))
		_, err := facades.Orm().Query().Exec(sql, values...)

		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("row %d: %s", i, err.Error()))
		} else {
			successCount++
		}
	}

	return response.Success(ctx, http.Json{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	})
}
