package database

import (
	"alta-wedding/config"
	"alta-wedding/models"

	"golang.org/x/crypto/bcrypt"
)

// Fungsi untuk mengambil dan mencari data organizer by email di database
func FindOrganizerByEmail(email string) (*models.Organizer, error) {
	organizer := models.Organizer{}
	tx := config.DB.Where("email=?", email).Find(&organizer)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		return &organizer, nil
	}
	return nil, nil
}

// Fungsi untuk mengambil dan mencari data organizer by id di database
func FindOrganizerById(id int) (*models.Organizer, error) {
	organizer := models.Organizer{}
	tx := config.DB.Where("id=?", id).Find(&organizer)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		return &organizer, nil
	}
	return nil, nil
}

// Fungsi untuk menambahkan organizer ke dalam database
func InsertOrganizer(newOrganizer models.Organizer) (models.Organizer, error) {
	if err := config.DB.Save(&newOrganizer).Error; err != nil {
		return models.Organizer{}, err
	}
	return newOrganizer, nil
}

// Fungsi untuk login organizer berdasarkan data yang ada pada database
func LoginOrganizer(login models.LoginRequestBody) (*models.Organizer, error) {
	organizerData, err := FindOrganizerByEmail(login.Email)
	if organizerData == nil || err != nil {
		return nil, err
	}
	check := CheckPasswordHash(login.Password, organizerData.Password)
	if !check {
		return nil, nil
	}
	return organizerData, nil
}

// Fungsi untuk Edit Profile Organizer
func EditOrganizer(reqOrganizer models.Organizer, organizer_id int) (*models.Organizer, error) {
	organizer := models.Organizer{}
	tx := config.DB.Find(&organizer, organizer_id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected < 1 {
		return nil, nil
	}
	if err := config.DB.Model(&models.Organizer{}).Where("id=?", organizer_id).Updates(reqOrganizer).Error; err != nil {
		return nil, err
	}
	return &organizer, nil
}

// Fungsi untuk Edit Photo Profile Organizer
func EditPhotoOrganizer(url string, organizer_id int) (int64, error) {
	tx := config.DB.Model(&models.Organizer{}).Where("id=?", organizer_id).Update("logo", url)
	if tx.Error != nil {
		return -1, tx.Error
	}
	return tx.RowsAffected, nil
}

// Fungsi untuk get organizer by ID
func FindProfilOrganizer(id int) (*models.ProfileRespon, error) {
	organizer := models.ProfileRespon{}
	tx := config.DB.Table("organizers").Select(
		"organizers.id, organizers.wo_name, organizers.email, organizers.phone_number, organizers.about, organizers.web_url, organizers.status, organizers.logo, organizers.city, organizers.address").
		Where("organizers.deleted_at IS NULL AND organizers.id=?", id).Find(&organizer)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &organizer, nil
}

// Fungsi untuk enkripsi password organizer
func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Fungsi untuk compare password organizer dengan enkripsi password organizer
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
