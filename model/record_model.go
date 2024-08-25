package model

import (
	"reflect"
	"time"

	"github.com/gobeam/stringy"
	"github.com/rs/zerolog/log"
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

func (r Record) ExtractChangedData(safeData map[string]interface{}) map[string]interface{} {
	changedData := make(map[string]interface{})

	reflectRecord := reflect.ValueOf(r)
	for field, newValue := range safeData {
		fieldKey := stringy.New(field).CamelCase("?", "").UcFirst()
		rField := reflectRecord.FieldByName(fieldKey)
		if rField.IsValid() {
			currentValue := rField.Interface().(interface{})

			if !areEqual(newValue, currentValue) {
				log.Debug().Msgf("Field: %s", fieldKey)
				log.Debug().Msgf("Changed")
				changedData[field] = newValue
			}
		}
	}

	return changedData
}

func areEqual(newValue interface{}, currentValue interface{}) bool {
	if reflect.TypeOf(newValue) == reflect.TypeOf(currentValue) {
		return reflect.DeepEqual(newValue, currentValue)
	}

	if reflect.TypeOf(currentValue) == reflect.TypeOf(time.Time{}) {
		parsedTime, err := time.Parse(time.RFC3339, newValue.(string))
		if err != nil {
			log.Error().Err(err).Msg("Couldn't Convert to Time")
			return false
		}
		return reflect.DeepEqual(parsedTime, currentValue)
	}

	return false
}

func (r Record) GetData() map[string]interface{} {
	data := make(map[string]interface{})

	reflectRecord := reflect.ValueOf(r)
	for _, field := range r.MutableFields() {
		fieldKey := stringy.New(field).CamelCase("?", "").UcFirst()
		rField := reflectRecord.FieldByName(fieldKey)
		if rField.IsValid() {
			data[field] = reflectRecord.FieldByName(fieldKey).
				Interface().(interface{})
		}
	}

	return data
}

func (r Record) MergeData(changedData map[string]interface{}) map[string]interface{} {
	mergedData := r.GetData()

	for field, newValue := range changedData {
		mergedData[field] = newValue
	}

	return mergedData
}
