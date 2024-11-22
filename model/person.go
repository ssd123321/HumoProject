package model

import (
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
