package types

import (
	"time"

	"github.com/google/uuid"
)

// Base model to include ID as UUID
type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseModel
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	AvatarURL string `gorm:"default:''"`
	Bio       string `gorm:"type:text"`
	Posts     []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	BaseModel
	UserID      uuid.UUID `gorm:"not null;index"`
	Content     string    `gorm:"type:text;not null"`
	Tags 	  []string  `gorm:"type:text[]"`
	User        User      `gorm:"foreignKey:UserID"`
	Comments    []Comment `gorm:"foreignKey:PostID"`
	PostPlugins []PostPlugin `gorm:"foreignKey:PostID"`
}

type Comment struct {
	BaseModel
	UserID  uuid.UUID `gorm:"not null;index"`
	PostID  uuid.UUID `gorm:"not null;index"`
	Content string    `gorm:"type:text;not null"`
	User    User      `gorm:"foreignKey:UserID"`
	Post    Post      `gorm:"foreignKey:PostID"`
}

type PostPlugin struct {
	BaseModel
	PostID    uuid.UUID   `gorm:"not null;index"`
	Type      string `gorm:"not null"` // e.g., "image", "gif", "sticker"
	URL       string `gorm:"not null"`
	HTML      string
	Post      Post `gorm:"foreignKey:PostID"`
}

type Like struct {
	BaseModel
	UserID    uuid.UUID `gorm:"not null;index"`
	PostID    uuid.UUID `gorm:"not null;index"`
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

type Dislike struct {
	BaseModel
	UserID    uuid.UUID `gorm:"not null;index"`
	PostID    uuid.UUID `gorm:"not null;index"`
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

type Follow struct {
	BaseModel
	FollowerID  uuid.UUID `gorm:"not null;index"`
	FollowingID uuid.UUID `gorm:"not null;index"`
	Follower    User `gorm:"foreignKey:FollowerID"`
	Following   User `gorm:"foreignKey:FollowingID"`
}
