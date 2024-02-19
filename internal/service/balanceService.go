package service

import (
	"avitoTeastProj/internal/models"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

type UserService struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}

func NewUserService(db *gorm.DB, logger *zap.SugaredLogger) *UserService {
	return &UserService{
		DB:     db,
		Logger: logger,
	}
}

// UpdateBalance функция обновляет баланс пользователя. Если пользователь не существует, он будет создан.
func (s UserService) UpdateBalance(id uint, balance uint) error {
	// Логируем начало операции обновления баланса
	s.Logger.Infow("Attempting to update user balance", "userID", id, "newBalance", balance)

	// Инициализация пустого объекта User
	user := &models.User{}

	// Проверка наличия пользователя в базе данных или его создание, если отсутствует
	res := s.DB.FirstOrCreate(user, models.User{Model: gorm.Model{ID: id}})
	if res.Error != nil {
		// Вывод в лог ошибки для операции FirstOrCreate
		s.Logger.Errorw("Error in Save during UpdateBalance", "error", res.Error, "userID", id)
		return res.Error
	}

	// Обновление баланса пользователя
	user.Balance = balance

	// Сохранение обновленного пользователя в базу данных
	res = s.DB.Save(user)
	if res.Error != nil {
		// Вывод в лог ошибки для операции Save
		log.Fatal("UpdateBalance, error in Save", res.Error)
	}
	// Если ошибок не произошло, функция возвращает nil вместо ошибки
	return nil
}

// Метод резервирования средств с основного баланса на отдельном счете. Принимает id пользователя, ИД услуги, ИД заказа, стоимость.
func (s UserService) ReserveMoney(userId uint, serviceId uint, orderId uint, price uint) error {
	s.Logger.Infow("Attempting to accept money", "userId", userId, "serviceId", serviceId, "orderId", orderId, "price", price)
	reserve := models.Reserve{
		UserID:         userId,
		ServiceID:      serviceId,
		OrderID:        orderId,
		Price:          price,
		ReservedAmount: price,
	}
	tx := s.DB.Set("gorm:query_option", "FOR UPDATE").Begin()

	if err := tx.Create(&reserve).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ReserveMoney, error in Create: %w", err)
	}

	user := &models.User{}

	if err := tx.First(user, userId).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ReserveMoney, error in First: %w", err)
	}

	if user.Balance < price {
		tx.Rollback()
		return fmt.Errorf("ReserveMoney, insufficient balance")
	}

	user.Balance -= price
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ReserveMoney, error in Save: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("ReserveMoney, error in Commit: %w", err)
	}

	s.Logger.Infow("Successfully accepted money", "userId", userId, "price", price)
	return nil
}

// AcceptMoney Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии.
// Принимает id пользователя, ИД услуги, ИД заказа, сумму.
func (s UserService) AcceptMoney(userId uint, serviceId uint, orderId uint, price uint) error {
	s.Logger.Infow("Attempting to accept money", "userId", userId, "serviceId", serviceId, "orderId", orderId, "price", price)
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var reverse models.Reserve
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("UserID = ? AND ServiseID = ? AND OrderID = ?", userId, serviceId, orderId).First(reverse).Error; err != nil {
		tx.Rollback() // Откат транзакции в случае ошибки
		return fmt.Errorf("AcceptMoney, reserve not found: %w", err)
	}

	if reverse.ReservedAmount < price {
		tx.Rollback()
		return fmt.Errorf("AcceptMoney, requested amount exceeds reserved amount")
	}

	report := models.ReportEntry{
		UserID:  userId,
		OrderID: orderId,
		Revenue: price,
	}
	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("AcceptMoney, error creating report entry: %w", err)
	}

	//уменьшение резерва
	reverse.ReservedAmount -= price
	if err := tx.Save(&reverse).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("AcceptMoney, error updating reserve: %w", err)
	}

	// Проверка, можно ли удалить резерв (если зарезервированная сумма равна 0)
	if reverse.ReservedAmount == 0 {
		if err := tx.Delete(&reverse).Error; err != nil {
			tx.Rollback() // Откат транзакции в случае ошибки при удалении резерва
			return fmt.Errorf("AcceptMoney, error deleting empty reserve: %w", err)
		}
	}

	// Подтверждение транзакции
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("AcceptMoney, error committing transaction: %w", err)
	}

	return nil

}

// GetBalance Метод получения баланса пользователя. Принимает id пользователя.
func (s UserService) GetBalance(userID uint) (uint, error) {
	s.Logger.Infow("Fetching user balance", "userID", userID)
	var user models.User

	res := s.DB.First(&user, userID)
	if res.Error != nil {
		return 0, fmt.Errorf("GetBalance, error finding user with ID %d: %w", userID, res.Error)
	}
	s.Logger.Infow("Successfully fetched user balance", "userID", userID, "balance", user.Balance)
	return user.Balance, nil
}
