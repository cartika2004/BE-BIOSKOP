package database

// Biar otomatis ada isinya dbnya
import (
	"fmt"
	"log"
	"project-bioskop/models"
	"time"
)

// SeedAll memanggil semua seeder
func SeedAll() {
	seedMovies()
	seedStudiosAndSeats()
	seedSchedules() 
}

func seedMovies() {
	var count int64
	DB.Model(&models.Movie{}).Count(&count)
	if count > 0 {
		return
	}

	movies := []models.Movie{
		{Title: "Avengers: Secret Wars", Duration: 180, Genre: "Action", Poster: "avengers.jpg", Description: "Perang Multiverse dimulai."},
		{Title: "Avatar 3: The Seed Bearer", Duration: 190, Genre: "Sci-Fi", Poster: "avatar.jpg", Description: "Kembali ke Pandora."},
	}

	DB.Create(&movies)
	log.Println("Data Film Berhasil Ditanam!")
}

func seedStudiosAndSeats() {
	var count int64
	DB.Model(&models.Studio{}).Count(&count)
	if count > 0 {
		return
	}

	// 1. Bikin Studio 1
	studio1 := models.Studio{Name: "Studio 1 (Regular)", Capacity: 20}
	DB.Create(&studio1)
	generateSeatsForStudio(studio1.ID, []string{"A", "B"}, 10)

	// 2. Bikin Studio 2
	studio2 := models.Studio{Name: "Studio 2 (VIP)", Capacity: 10}
	DB.Create(&studio2)
	generateSeatsForStudio(studio2.ID, []string{"A"}, 10)

	log.Println("Data Studio & Kursi Berhasil Ditanam!")
}

func generateSeatsForStudio(studioID uint, rows []string, cols int) {
	var seats []models.Seat
	for _, row := range rows {
		for i := 1; i <= cols; i++ {
			seatNumber := fmt.Sprintf("%s%d", row, i)
			seats = append(seats, models.Seat{
				StudioID:   studioID,
				SeatNumber: seatNumber,
			})
		}
	}
	DB.Create(&seats)
}

// --- FUNGSI BARU UNTUK JADWAL ---
func seedSchedules() {
	var count int64
	DB.Model(&models.Schedule{}).Count(&count)
	if count > 0 {
		return // Kalau sudah ada jadwal (termasuk yg kamu buat manual), skip aja biar gak dobel
	}

	// Ambil data Film & Studio yg sudah ada
	var movie1, movie2 models.Movie
	DB.First(&movie1, 1) // ID 1: Avengers
	DB.First(&movie2, 2) // ID 2: Avatar

	var studio1, studio2 models.Studio
	DB.First(&studio1, 1) // ID 1: Regular
	DB.First(&studio2, 2) // ID 2: VIP

	// Tentukan tanggal: misal BESOK dan LUSA (biar valid buat dipesan)
	tomorrow := time.Now().Add(24 * time.Hour)
	dayAfter := time.Now().Add(48 * time.Hour)

	schedules := []models.Schedule{
		// Jadwal 1: Avengers di Studio 1 (Besok jam 13:00)
		{
			MovieID:   movie1.ID,
			StudioID:  studio1.ID,
			StartTime: time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 13, 0, 0, 0, time.Local),
			Price:     40000,
		},
		// Jadwal 2: Avengers di Studio 1 (Besok jam 16:00)
		{
			MovieID:   movie1.ID,
			StudioID:  studio1.ID,
			StartTime: time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 16, 0, 0, 0, time.Local),
			Price:     40000,
		},
		// Jadwal 3: Avatar di Studio 2 VIP (Lusa jam 19:00)
		{
			MovieID:   movie2.ID,
			StudioID:  studio2.ID,
			StartTime: time.Date(dayAfter.Year(), dayAfter.Month(), dayAfter.Day(), 19, 0, 0, 0, time.Local),
			Price:     75000,
		},
	}

	DB.Create(&schedules)
	log.Println("Data Jadwal (Schedules) Berhasil Ditanam!")
}