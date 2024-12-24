package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var NotFound = fmt.Errorf("record not found")

type StringArray []string

type ReturnError struct {
	Status  string
	Message string
}
type JwtPair struct {
	AccessJWT  string `json:"access_jwt,omitempty"`
	RefreshJWT string `json:"refresh_jwt,omitempty"`
}
type Content struct {
	Name      string `json:"name,omitempty"`
	Age       int    `json:"age,omitempty"`
	Dimension struct {
		Weight float64 `json:"weight,omitempty"`
		Height float64 `json:"height,omitempty"`
	} `json:"dimension,embedded" gorm:"embedded"`
}
type Person struct {
	ID        int        `gorm:"primaryKey:id"`
	Password  string     `json:"password,omitempty" gorm:"column:current_password"`
	Login     string     `json:"login,omitempty" gorm:"column:login"`
	Content   Content    `json:"content,omitempty" gorm:"type:jsonb"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Cache     bool       `json:"cache" gorm:"column:cache"`
}
type SigningRequest struct {
	Login    string `json:"login,omitempty" gorm:"column:login"`
	Password string `json:"password,omitempty" gorm:"column:current_password"`
}
type PeopleNoContent struct {
	ID        int        `gorm:"primaryKey:id"`
	Content   Content    `json:"content,omitempty" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Cache     bool       `json:"cache" gorm:"column:cache"`
}
type Card struct {
	ID           int        `json:"id,omitempty" gorm:"primaryKey:id"`
	PersonID     int        `json:"person_id,omitempty" gorm:"column:person_id"`
	DateOfExpire string     `json:"date_of_expire,omitempty" gorm:"column:date_of_expire"`
	CardNumber   int        `json:"card_number,omitempty" gorm:"column:card_number"`
	Logotype     string     `json:"logotype,omitempty" gorm:"column:logotype"`
	Money        float64    `json:"money,omitempty" gorm:"column:money"`
	BankName     string     `json:"bank_name,omitempty" gorm:"column:bank_name"`
	CreatedAt    *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}
type PersonNew struct {
	ID              int         `gorm:"primaryKey:id"`
	Passwords       StringArray `gorm:"column:passwords;type:text[]"`
	CurrentPassword string      `json:"password,omitempty" gorm:"column:current_password"`
	Login           string      `json:"login,omitempty" gorm:"column:login"`
	Content         Content     `json:"content,omitempty" gorm:"type:jsonb"`
	CreatedAt       *time.Time  `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt       *time.Time  `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt       *time.Time  `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Cache           bool        `json:"cache" gorm:"column:cache"`
}
type NewPassword struct {
	Id              int    `json:"id" gorm:"column:id"`
	NewPass         string `json:"new_password"`
	CurrentPassword string `json:"current_password" gorm:"current_password"`
}

var CorporativeCard Card = Card{
	ID: 2,
}

var Sender Card = Card{}

func (c *Content) Scan(src interface{}) error {
	var data []byte
	switch v := src.(type) {
	case []uint8:
		data = v
	case string:
		data = []byte(v)
	}
	return json.Unmarshal(data, c)
}

func (c *Content) Value() (driver.Value, error) {
	return json.Marshal(c)
}
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	switch v := value.(type) {
	case string:
		str := strings.Trim(v, "{}")
		*s = strings.Split(str, ",")
	case []byte:
		str := strings.Trim(string(v), "{}")
		*s = strings.Split(str, ",")
	default:
		return errors.New("failed to scan StringArray: invalid type")
	}

	return nil
}
