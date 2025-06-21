package protocol

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Protocol messages for booking service
type BookingMessage struct {
	ID         string `json:"id"`
	CreatedBy  string `json:"created_by"`
	PropertyID string `json:"property_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	GuestName  string `json:"guest_name"`
}

type CreateBookingMessage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GuestName string `json:"guest_name"`
}

type UpdateBookingMessage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GuestName string `json:"guest_name"`
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

type GroupMessage struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
}

func Sha256Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
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
