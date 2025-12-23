package models

import (
	"time"

	"gorm.io/gorm"
)

// 1. DATA USER
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role" gorm:"default:'user'"` // 'admin' atau 'user'
}

// 2. DATA FILM
type Movie struct {
	gorm.Model
	Title       string `json:"title"`
	Poster      string `json:"poster"`
	Duration    int    `json:"duration"` // menit
	Description string `json:"description"`
	Genre       string `json:"genre"`
}

// 3. DATA STUDIO
type Studio struct {
	gorm.Model
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Seats    []Seat `json:"seats"` // 1 Studio punya Banyak Kursi
}

// 4. DATA KURSI (Master Data)
type Seat struct {
	gorm.Model
	StudioID   uint   `json:"studio_id"`
	SeatNumber string `json:"seat_number"` // A1, A2, B1...
}

// 5. DATA JADWAL (Schedule)
type Schedule struct {
	gorm.Model
	MovieID   uint      `json:"movie_id"`
	Movie     Movie     `json:"movie"`
	StudioID  uint      `json:"studio_id"`
	Studio    Studio    `json:"studio"`
	StartTime time.Time `json:"start_time"`
	Price     float64   `json:"price"`
}

// 6. TRANSAKSI (Nota Header)
type Transaction struct {
	gorm.Model
	UserID      uint     `json:"user_id"`
	User        User     `json:"user"`
	TotalAmount float64  `json:"total_amount"`
	Status      string   `json:"status"` // PENDING, PAID, CANCEL
	Tickets     []Ticket `json:"tickets"`
}

// 7. TIKET (Detail Transaksi - Kursi yang dipesan)
type Ticket struct {
	gorm.Model
	TransactionID uint     `json:"transaction_id"`
	ScheduleID    uint     `json:"schedule_id"`
	Schedule      Schedule `json:"schedule"`
	SeatID        uint     `json:"seat_id"`
	Seat          Seat     `json:"seat"`
}
