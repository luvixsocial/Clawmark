package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	AvatarURL string `gorm:"default:''"`
	Bio       string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Posts     []Post         `gorm:"foreignKey:UserID"`
	Comments  []Comment      `gorm:"foreignKey:UserID"`
	Likes     []Like         `gorm:"foreignKey:UserID"`
	Dislikes  []Dislike      `gorm:"foreignKey:UserID"`
	Follows   []Follow       `gorm:"foreignKey:FollowerID"`
}

type Post struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"`
	Content     string `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	User        User           `gorm:"foreignKey:UserID"`
	Comments    []Comment      `gorm:"foreignKey:PostID"`
	Likes       []Like         `gorm:"foreignKey:PostID"`
	Dislikes    []Dislike      `gorm:"foreignKey:PostID"`
	PostPlugins []PostPlugin   `gorm:"foreignKey:PostID"`
}

type PostPlugin struct {
	ID        uint   `gorm:"primaryKey"`
	PostID    uint   `gorm:"not null;index"`
	Type      string `gorm:"not null"` // e.g., "image", "gif", "sticker"
	URL       string `gorm:"not null"`
	HTML      string
	CreatedAt time.Time
	Post      Post `gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	PostID    uint   `gorm:"not null;index"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignKey:UserID"`
	Post      Post           `gorm:"foreignKey:PostID"`
}

type Like struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null;index"`
	PostID    uint `gorm:"not null;index"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

type Dislike struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null;index"`
	PostID    uint `gorm:"not null;index"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
}

type Follow struct {
	ID          uint `gorm:"primaryKey"`
	FollowerID  uint `gorm:"not null;index"`
	FollowingID uint `gorm:"not null;index"`
	CreatedAt   time.Time
	Follower    User `gorm:"foreignKey:FollowerID"`
	Following   User `gorm:"foreignKey:FollowingID"`
}
