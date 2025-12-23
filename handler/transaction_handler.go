package handler

import (
	"encoding/json"
	"net/http"
	"project-bioskop/middleware"
	"project-bioskop/service"
	"strconv"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.KeyUserID).(float64)

	var input struct {
		ScheduleID uint   `json:"schedule_id"`
		SeatIDs    []uint `json:"seat_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	trx, err := h.service.CreateBooking(uint(userID), input.ScheduleID, input.SeatIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking Success!",
		"data":    trx,
	})
}

func (h *TransactionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.KeyUserID).(float64)

	history, err := h.service.GetUserHistory(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "My Booking History",
		"data":    history,
	})
}

func (h *TransactionHandler) GetSeatsBySchedule(w http.ResponseWriter, r *http.Request) {
	// /schedules/{id}/seats
	scheduleIDStr := r.PathValue("id")
	scheduleID, _ := strconv.Atoi(scheduleIDStr)

	seats, err := h.service.GetSeatAvailability(uint(scheduleID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"schedule_id": scheduleID,
		"seats":       seats,
	})
}