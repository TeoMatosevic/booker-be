package protocol

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Protocol messages for booking service
type BookingMessage struct {
	ID         string `json:"id"`
	CreatedBy  string `json:"created_by"`
	PropertyID string `json:"property_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	GuestName  string `json:"guest_name"`
	Adults     int    `json:"adults"`
	Children   int    `json:"children"`
}

type CreateBookingMessage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GuestName string `json:"guest_name"`
	Adults    int    `json:"adults"`
	Children  int    `json:"children"`
}

type UpdateBookingMessage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GuestName string `json:"guest_name"`
	Adults    int    `json:"adults"`
	Children  int    `json:"children"`
}

// Protocol messages for user service
type UserMessage struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type CreateUserMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Protocol messages for group code service
type GroupCodeMessage struct {
	GroupID string `json:"group_id"`
}

// Protocol messages for group service
type GroupCreateMessage struct {
	Name string `json:"name"`
}

// Protocol messages for property service
type CreatePropertyMessage struct {
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
}

type UpdatePropertyMessage struct {
	Color string `json:"color"`
}

type GroupMessage struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
}

// HashPassword generates a bcrypt hash of the password
// Cost factor of 14 provides good security while maintaining reasonable performance
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a bcrypt hashed password with its plaintext version
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateID() string {
	return uuid.New().String()
}

func GetCurrentTime() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func IsValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
