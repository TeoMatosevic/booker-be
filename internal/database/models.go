package database

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

type Group struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
}

type Property struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	GroupID   string `json:"group_id"`
	Name      string `json:"name"`
}

type Booking struct {
	ID         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	CreatedBy  string `json:"created_by"`
	PropertyID string `json:"property_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	GuestName  string `json:"guest_name"`
}

type GroupUser struct {
	ID      string `json:"id"`
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

type GroupCode struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	Code     string `json:"code"`
	ActiveTo string `json:"active_to"`
}
