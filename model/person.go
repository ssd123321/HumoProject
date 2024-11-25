package model

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"time"
)

type Content struct {
	Name      string `json:"name,omitempty"`
	Age       int    `json:"age,omitempty"`
	Dimension struct {
		Weight float64 `json:"weight,omitempty"`
		Height float64 `json:"height,omitempty"`
	} `json:"dimension,embedded" gorm:"embedded"`
}

/*
	func (c *Content) Value() (driver.Value, error) {
		return json.Marshal(c)
	}

	func (c *Content) Scan(value interface{}) error {
		b, ok := value.([]byte)
		if !ok {
			return errors.New("type assertion to byte")
		}
		return json.Unmarshal(b, &c)
	}
*/
type Person struct {
	ID        int        `gorm:"primaryKey:id"`
	Content   Content    `json:"content,omitempty" gorm:"type:jsonb"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Cache     bool       `json:"cache" gorm:"column:cache"`
}
type DBPerson struct {
	ID        int        `gorm:"primaryKey:id"`
	Content   string     `json:"content,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Cache     bool       `json:"cache" gorm:"column:cache"`
}
type PeopleNoContent struct {
	ID        int        `gorm:"primaryKey:id"`
	Content   Content    `json:"content,omitempty" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Cache     bool       `json:"cache" gorm:"column:cache"`
}
type CardContent struct {
	CardNumber   int    `json:"card_number"`
	Logotype     string `json:"logotype"`
	BankName     string `json:"bankname"`
	DateOfExpire string `json:"date_of_expire"`
}
type Card struct {
	ID        int         `json:"id,omitempty" gorm:"primaryKey:id"`
	PersonID  int         `json:"person_id,omitempty" gorm:"column:person_id"`
	Content   CardContent `json:"content,omitempty" gorm:"type:jsonb"`
	CreatedAt *time.Time  `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time  `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (c *CardContent) Scan(src interface{}) error {
	log.Println(1)
	var data []byte
	switch v := src.(type) {
	case []uint8:
		data = v
	case string:
		data = []byte(v)
	}
	return json.Unmarshal(data, c)
}

func (c *CardContent) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Content) Scan(src interface{}) error {
	log.Println(1)
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
