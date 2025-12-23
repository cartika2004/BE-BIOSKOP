package service

import (
	"errors"
	"project-bioskop/models"
	"project-bioskop/repository"
)

type TransactionService interface {
	CreateBooking(userID uint, scheduleID uint, seatIDs []uint) (models.Transaction, error)
	GetUserHistory(userID uint) ([]models.Transaction, error)
	GetSeatAvailability(scheduleID uint) ([]SeatStatus, error)
}

type SeatStatus struct {
	SeatID      uint   `json:"seat_id"`
	SeatNumber  string `json:"seat_number"`
	IsAvailable bool   `json:"is_available"`
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

func (s *transactionService) CreateBooking(userID uint, scheduleID uint, seatIDs []uint) (models.Transaction, error) {
	price, err := s.repo.GetSchedulePrice(scheduleID)
	if err != nil {
		return models.Transaction{}, errors.New("schedule not found")
	}

	trx := models.Transaction{
		UserID:      userID,
		Status:      "PAID",
		TotalAmount: 0,
		Tickets:     []models.Ticket{},
	}

	for _, seatID := range seatIDs {
		if !s.repo.IsSeatAvailable(scheduleID, seatID) {
			return models.Transaction{}, errors.New("seat already booked")
		}

		ticket := models.Ticket{
			ScheduleID: scheduleID,
			SeatID:     seatID,
		}
		trx.Tickets = append(trx.Tickets, ticket)
		trx.TotalAmount += price
	}

	if err := s.repo.CreateTransaction(&trx); err != nil {
		return models.Transaction{}, err
	}

	return trx, nil
}

func (s *transactionService) GetUserHistory(userID uint) ([]models.Transaction, error) {
	return s.repo.FindByUserID(userID)
}

// --- (Logic Cek Kursi) ---
func (s *transactionService) GetSeatAvailability(scheduleID uint) ([]SeatStatus, error) {
	// Cek Studio ID dari Jadwal
	studioID, err := s.repo.GetStudioIDBySchedule(scheduleID)
	if err != nil {
		return nil, err
	}

	// Ambil Semua Kursi Fisik
	allSeats, err := s.repo.GetSeatsByStudioID(studioID)
	if err != nil {
		return nil, err
	}

	// Ambil Kursi yang sudah laku
	bookedSeatIDs, err := s.repo.GetBookedSeatIDs(scheduleID)
	if err != nil {
		return nil, err
	}

	// Optimization: Pakai Map biar ngeceknya cepat
	bookedMap := make(map[uint]bool)
	for _, id := range bookedSeatIDs {
		bookedMap[id] = true
	}

	// Gabungkan Status
	var result []SeatStatus
	for _, seat := range allSeats {
		isBooked := bookedMap[seat.ID]

		result = append(result, SeatStatus{
			SeatID:      seat.ID,
			SeatNumber:  seat.SeatNumber,
			IsAvailable: !isBooked, // True jika belum booked
		})
	}

	return result, nil
}