package service

import (
	"avitoTeastProj/internal/models"
	"gorm.io/gorm"
	"time"
)

type MonthlyRevenueReport struct {
	UserID            uint    `json:"userId"`
	TotalRevenue      float64 `json:"totalRevenue"`
	TransactionsCount int     `json:"transactionsCount"`
}

func GenerateMonthlyRevenueReport(db *gorm.DB, year int, month uint) ([]MonthlyRevenueReport, error) {
	var reports []MonthlyRevenueReport
	// Формируем начало и конец месяца для фильтрации данных
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// Выполняем запрос к базе данных
	err := db.Model(&models.ReportEntry{}).
		Select("user_id, SUM(revenue) AS total_revenue, COUNT(*) AS transactions_count").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Group("user_id").
		Scan(&reports).Error

	if err != nil {
		return nil, err
	}

	return reports, nil
}
