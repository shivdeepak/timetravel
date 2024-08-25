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

	// GetRecordAt will retrieve a record's version at a given timestamp
	GetRecordAt(ctx context.Context, id uint, at time.Time) (model.Record, error)

	// GetVersions will retrieve all versions of record.
	GetVersions(ctx context.Context, id uint) ([]model.Record, error)

	// CreateRecord will insert a new record.
	//
	// If it a record with that id already exists it will fail.
	CreateRecord(ctx context.Context, id uint, unsafeData map[string]interface{}) (model.Record, error)

	// UpdateRecord will change the internal `Map` values of the record if they exist.
	// if the update[key] is null it will delete that key from the record's Map.
	//
	// UpdateRecord will error if id <= 0 or the record does not exist with that id.
	UpdateRecord(ctx context.Context, prevRecord model.Record, unsafeData map[string]interface{}) (model.Record, error)
}

// SQLiteRecordService is a SQLite implementation of RecordService.
type SQLiteRecordService struct{}

func NewSQLiteRecordService() SQLiteRecordService {
	return SQLiteRecordService{}
}

func (s *SQLiteRecordService) GetRecord(ctx context.Context, id uint) (model.Record, error) {
	return s.GetRecordAt(ctx, id, time.Now())
}

func (s *SQLiteRecordService) GetRecordAt(ctx context.Context, id uint, at time.Time) (model.Record, error) {
	db := model.GetDb()

	var record model.Record
	result := db.Order("updated_at desc").
		Where("updated_at <= ?", at.Format(time.RFC3339)).
		First(&record, id)
	if result.Error != nil {
		return model.Record{}, result.Error
	}

	return record, nil
}

func (s *SQLiteRecordService) GetVersions(ctx context.Context, id uint) ([]model.Record, error) {
	db := model.GetDb()

	var records []model.Record
	result := db.Order("updated_at desc").Find(&records, id)
	if result.Error != nil {
		return []model.Record{}, result.Error
	}

	return records, nil
}

func (s *SQLiteRecordService) CreateRecord(ctx context.Context, id uint, unsafeData map[string]interface{}) (model.Record, error) {
	log.Debug().Msg("CreateRecord")

	safeData := model.Record{}.SanitizePayload(unsafeData, false)

	numSafeFields := len(safeData)

	db := model.GetDb()
	if numSafeFields > 0 {
		log.Debug().Msg("Running Create")
		safeData["id"] = id
		safeData["created_at"] = time.Now().Format(time.RFC3339)
		safeData["updated_at"] = time.Now().Format(time.RFC3339)
		result := db.Model(&model.Record{}).Create(safeData)
		if result.Error != nil {
			logging.LogError(result.Error)
			return model.Record{}, result.Error
		} else {
			log.Debug().Msg("Record Created")
			return s.GetRecord(ctx, id)
		}
	} else {
		log.Debug().Msg("Skipped Create, Nothing to Create!")
		return model.Record{}, errors.New("No Fields to Create Record!")
	}
}

func (s *SQLiteRecordService) UpdateRecord(ctx context.Context, prevRecord model.Record, unsafeData map[string]interface{}) (model.Record, error) {
	log.Debug().Msg("UpdateRecord")

	safeData := model.Record{}.SanitizePayload(unsafeData, false)

	numSafeFields := len(safeData)

	db := model.GetDb()
	if numSafeFields > 0 {
		log.Debug().Msg("Running Updated")
		safeData["id"] = prevRecord.ID
		safeData["created_at"] = prevRecord.CreatedAt.Format(time.RFC3339)
		safeData["updated_at"] = time.Now().Format(time.RFC3339)
		result := db.Model(&model.Record{}).Create(safeData)
		if result.Error != nil {
			logging.LogError(result.Error)
			return model.Record{}, result.Error
		} else {
			log.Debug().Msg("Record Updated")
			return s.GetRecord(ctx, prevRecord.ID)
		}
	} else {
		log.Debug().Msg("Skipped Update, Nothing to Update!")
		return prevRecord, nil
	}
}
