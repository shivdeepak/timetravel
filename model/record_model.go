package model

import (
	"reflect"
	"time"

	"github.com/gobeam/stringy"
	"gorm.io/gorm"
)

type Record struct {
	ID        uint           `gorm:"primaryKey;autoIncrement:false" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `gorm:"primaryKey" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Dob        time.Time `json:"dob"`
	Phone      string    `json:"phone"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Zip        string    `json:"zip"`
	Country    string    `json:"country"`
}

type RecordJSON struct {
	ID   uint                   `json:"id"`
	Data map[string]interface{} `json:"data"`
}

func (r Record) ToJSON() (RecordJSON, error) {
	v := reflect.ValueOf(r)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	result := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Name == "ID" || field.Name == "DeletedAt" {
			continue
		}
		fieldKey := stringy.New(field.Name).SnakeCase().ToLower()
		result[fieldKey] = v.Field(i).Interface()
	}

	recordJson := RecordJSON{}
	recordJson.ID = r.ID
	recordJson.Data = result

	return recordJson, nil
}

func (r Record) MutableFields() []string {
	return []string{
		"first_name",
		"middle_name",
		"last_name",
		"email",
		"dob",
		"phone",
		"street",
		"city",
		"state",
		"zip",
		"country",
	}
}

func (r Record) SanitizePayload(unsafeData map[string]interface{}, preserveNil bool) map[string]interface{} {
	safeData := make(map[string]interface{})

	for _, safeField := range r.MutableFields() {
		value, ok := unsafeData[safeField]
		if ok {
			if preserveNil {
				safeData[safeField] = value
			} else {
				if value != nil {
					safeData[safeField] = value
				}
			}
		}
	}

	return safeData
}
