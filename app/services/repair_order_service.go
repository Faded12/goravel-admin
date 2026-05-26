package services

import (
	"github.com/goravel/framework/facades"

	"goravel/app/models"
)

type RepairOrderService interface {
	GetList(filters RepairOrderFilters, page, pageSize int) ([]models.RepairOrder, int64, error)
	GetByID(id uint) (*models.RepairOrder, error)
}

type RepairOrderFilters struct {
	WorksOrderNumber   string
	NotificationNumber string
	ContractNumber     string
	WorksOrderStatus   string
	StartDate          string
	EndDate            string
}

type RepairOrderServiceImpl struct{}

func NewRepairOrderService() RepairOrderService {
	return &RepairOrderServiceImpl{}
}

func (s *RepairOrderServiceImpl) GetList(filters RepairOrderFilters, page, pageSize int) ([]models.RepairOrder, int64, error) {
	var orders []models.RepairOrder
	var total int64

	query := facades.Orm().Query().Model(&models.RepairOrder{})

	if filters.WorksOrderNumber != "" {
		query = query.Where("works_order_number LIKE ?", "%"+filters.WorksOrderNumber+"%")
	}
	if filters.NotificationNumber != "" {
		query = query.Where("notification_number LIKE ?", "%"+filters.NotificationNumber+"%")
	}
	if filters.ContractNumber != "" {
		query = query.Where("contract_number", filters.ContractNumber)
	}
	if filters.WorksOrderStatus != "" {
		query = query.Where("works_order_status", filters.WorksOrderStatus)
	}
	if filters.StartDate != "" {
		query = query.Where("start_date >= ?", filters.StartDate)
	}
	if filters.EndDate != "" {
		query = query.Where("start_date <= ?", filters.EndDate)
	}

	err := query.Order("id desc").Paginate(page, pageSize, &orders, &total)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *RepairOrderServiceImpl) GetByID(id uint) (*models.RepairOrder, error) {
	var order models.RepairOrder
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&order); err != nil {
		return nil, err
	}

	return &order, nil
}
