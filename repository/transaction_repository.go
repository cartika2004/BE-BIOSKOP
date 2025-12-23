package repository

import (
	"project-bioskop/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	IsSeatAvailable(scheduleID uint, seatID uint) bool
	GetSchedulePrice(scheduleID uint) (float64, error)
	FindByUserID(userID uint) ([]models.Transaction, error)
	GetStudioIDBySchedule(scheduleID uint) (uint, error)
	GetSeatsByStudioID(studioID uint) ([]models.Seat, error)
	GetBookedSeatIDs(scheduleID uint) ([]uint, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) IsSeatAvailable(scheduleID uint, seatID uint) bool {
	var count int64
	r.db.Model(&models.Ticket{}).
		Where("schedule_id = ? AND seat_id = ?", scheduleID, seatID).
		Count(&count)

	return count == 0
}

func (r *transactionRepository) GetSchedulePrice(scheduleID uint) (float64, error) {
	var schedule models.Schedule
	if err := r.db.First(&schedule, scheduleID).Error; err != nil {
		return 0, err
	}
	return schedule.Price, nil
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	tx := r.db.Begin()

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *transactionRepository) FindByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.
		Preload("Tickets.Seat").
		Preload("Tickets.Schedule.Movie").
		Preload("Tickets.Schedule.Studio").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&transactions).Error

	return transactions, err
}

// --- (Seat Availability) ---

func (r *transactionRepository) GetStudioIDBySchedule(scheduleID uint) (uint, error) {
	var schedule models.Schedule
	err := r.db.Select("studio_id").First(&schedule, scheduleID).Error
	return schedule.StudioID, err
}

func (r *transactionRepository) GetSeatsByStudioID(studioID uint) ([]models.Seat, error) {
	var seats []models.Seat
	err := r.db.Where("studio_id = ?", studioID).Find(&seats).Error
	return seats, err
}

func (r *transactionRepository) GetBookedSeatIDs(scheduleID uint) ([]uint, error) {
	var bookedIDs []uint
	err := r.db.Model(&models.Ticket{}).
		Where("schedule_id = ?", scheduleID).
		Pluck("seat_id", &bookedIDs).Error
	return bookedIDs, err
}