package service

import (
	"context"
	"errors"
	"time"

	"github.com/rainbowmga/timetravel/concern/logging"
	"github.com/rainbowmga/timetravel/model"
	"github.com/rs/zerolog/log"
)

var ErrRecordDoesNotExist = errors.New("record with that id does not exist")
var ErrRecordIDInvalid = errors.New("record id must >= 0")
var ErrRecordAlreadyExists = errors.New("record already exists")

// Implements method to get, create, and update record data.
type RecordService interface {

	// GetRecord will retrieve an record.
	GetRecord(ctx context.Context, id uint) (model.Record, error)

	// CreateRecord will insert a new record.
	//
	// If it a record with that id already exists it will fail.
	CreateRecord(ctx context.Context, id uint, unsafeData map[string]interface{}) (model.Record, error)

	// UpdateRecord will change the internal `Map` values of the record if they exist.
	// if the update[key] is null it will delete that key from the record's Map.
	//
	// UpdateRecord will error if id <= 0 or the record does not exist with that id.
	UpdateRecord(ctx context.Context, record model.Record, unsafeData map[string]interface{}) (model.Record, error)
}

// SQLiteRecordService is a SQLite implementation of RecordService.
type SQLiteRecordService struct{}

func NewSQLiteRecordService() SQLiteRecordService {
	return SQLiteRecordService{}
}

func (s *SQLiteRecordService) GetRecord(ctx context.Context, id uint) (model.Record, error) {
	db := model.GetDb()

	var record model.Record
	result := db.First(&record, id)
	if result.Error != nil {
		return model.Record{}, result.Error
	}

	return record, nil
}

func (s *SQLiteRecordService) CreateRecord(ctx context.Context, id uint, unsafeData map[string]interface{}) (model.Record, error) {
	log.Debug().Msg("CreateRecord")

	var newRecord model.Record
	newRecord.ID = id
	safeData := newRecord.SanitizePayload(unsafeData, false)

	numSafeFields := len(safeData)

	db := model.GetDb()
	if numSafeFields > 0 {
		log.Debug().Msg("Running Create")
		safeData["created_at"] = time.Now().Format(time.RFC3339)
		safeData["updated_at"] = time.Now().Format(time.RFC3339)
		result := db.Model(&newRecord).Create(safeData)
		if result.Error != nil {
			logging.LogError(result.Error)
			return model.Record{}, result.Error
		} else {
			log.Debug().Msg("Record Created")
			return s.GetRecord(ctx, id)
		}
	} else {
		return model.Record{}, errors.New("No Fields to Update")
	}
}

func (s *SQLiteRecordService) UpdateRecord(ctx context.Context, record model.Record, unsafeData map[string]interface{}) (model.Record, error) {
	log.Debug().Msg("UpdateRecord")

	safeData := record.SanitizePayload(unsafeData, true)

	numSafeFields := len(safeData)

	db := model.GetDb()
	if numSafeFields > 0 {
		log.Debug().Msg("Running Updated")
		safeData["updated_at"] = time.Now().Format(time.RFC3339)
		result := db.Model(&record).Updates(safeData)
		if result.Error != nil {
			logging.LogError(result.Error)
			return model.Record{}, result.Error
		} else {
			log.Debug().Msg("Record Updated")
			return s.GetRecord(ctx, record.ID)
		}
	} else {
		return model.Record{}, errors.New("No Fields to Update")
	}
}
