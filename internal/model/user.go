package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	// NOTE: Not using gorm.Model since it's properties cannot be accessed directly
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FirstName string
	LastName  string
	Bio       string
	Verified  bool

	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string

	XP      uint64 `gorm:"default:500"`
	Arkhoin uint64 `gorm:"default:1000"`

	// FIXME: Make the followers/following be a queryable table
	// WARN: NOT TESTED

	// e.g.: UserFollow that includes the followerId and the followingId
	// Not very optimal
	Followers []*User
	Following []*User

	Comments []struct {
		OwnerID  uint
		Message  string
		Likes    uint
		Dislikes uint
	} `gorm:"type:json"`

	Notification struct {
		InApp bool `gorm:"default:true"`
		Email bool `gorm:"default:false"`
		// SMS bool `gorm:"default:false"` // NOTE: Not sure if we want to send SMS, it can get expensive
	} `gorm:"type:json"`

	Privacy struct {
		Profile    map[string]any `gorm:"type:json"`
		Image      map[string]any `gorm:"type:json"`
		Comments   map[string]any `gorm:"type:json"`
		Favorites  map[string]any `gorm:"type:json"`
		Projects   map[string]any `gorm:"type:json"`
		Components map[string]any `gorm:"type:json"`
		Followers  map[string]any `gorm:"type:json"`
		Following  map[string]any `gorm:"type:json"`
		Inventory  map[string]any `gorm:"type:json"`
		Formations map[string]any `gorm:"type:json"`
	} `gorm:"type:json"`

	// TODO: Implement enums Language, Theme and Country (country shouldnt be an enum)
}

type follow struct {
	gorm.Model
	UserID     uint
	FollowedID uint
}

type UserRegister struct {
	FirstName string `validate:"required,min=3"                  json:"firstname"`
	LastName  string `validate:"required,min=3"                  json:"lastname"`
	Username  string `validate:"required,username,min=3,max=32"  json:"username"`
	Email     string `validate:"required,email"                  json:"email"`
	Password  string `validate:"required,password,min=12,max=48" json:"password"`
}

type UserLogin struct {
	Username string `validate:"omitempty,username,min=3,max=32" json:"username"`
	Email    string `validate:"omitempty,email"                 json:"email"`
	Password string `validate:"required,password,min=12,max=48" json:"password"`
}
