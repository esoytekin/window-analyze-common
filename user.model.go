package common

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Uuid     string         `json:"uuid" gorm:"primaryKey"`
	Password string         `json:"-"`
	Role     pq.StringArray `json:"role" gorm:"type:text[]" swaggertype:"array,string"`
	Data     UserData       `json:"data"`
	Deleted  gorm.DeletedAt `json:"-"`
}

func (r *User) BeforeCreate(tx *gorm.DB) error {
	if r.Uuid == "" {
		r.Uuid = uuid.New().String()
	}
	return nil
}

type UserData struct {
	UserID      string         `json:"-" gorm:"foreignKey;primaryKey"`
	DisplayName string         `json:"displayName"`
	Email       string         `json:"email"`
	PhotoURL    string         `json:"photoURL"`
	Name        string         `json:"name"`
	Surname     string         `json:"surname"`
	Created     time.Time      `json:"createDate"`
	Deleted     gorm.DeletedAt `json:"-"`
}
