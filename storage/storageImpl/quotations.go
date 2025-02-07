package storageImpl

import (
	"chatbot/config"
	"chatbot/storage"
	"chatbot/storage/models"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var qute *quotations

type Quotations interface {
	GetRandomOne(t string) (string, error)
	GetOne(qute string) (string, error)
	AddOne(t string, text string) error
}

type quotations struct {
	db *gorm.DB
}

func InitQuotations() (Quotations, error) {
	if qute != nil {
		return qute, nil
	}
	db := storage.InitWithConfig(config.SqlDBConfig{
		Name: config.GlobalConfig.Storage.SqlDB.Quotations,
		Path: config.GlobalConfig.Storage.SqlDB.Path})
	if db == nil {
		log.Error().Msg("failed to init database")
		return nil, errors.New("failed to init database")
	}
	q := models.Quotation{}
	if err := db.AutoMigrate(q); err != nil {
		log.Error().Err(err).Msg("failed to auto migrate quotation table")
		return nil, err
	}
	qute = &quotations{db}
	return qute, nil
}

func (q quotations) GetRandomOne(t string) (string, error) {
	r := models.Quotation{}
	if err := q.db.Where("level = ?", t).Order("RANDOM()").First(&r).Error; err != nil {
		log.Error().Err(err).Msg("failed to get quotation")
		return "", err
	}
	return r.Text, nil
}

func (q quotations) AddOne(t string, text string) error {
	r := &models.Quotation{
		Text:  text,
		Level: t,
	}
	if err := q.db.Create(r).Error; err != nil {
		return err
	}
	log.Debug().Str(t, text).Msg("success add quotation")
	return nil
}

func (q quotations) GetOne(t string) (string, error) {
	r := models.Quotation{}
	if err := q.db.Where("text = ?", t).Error; err != nil {
		log.Error().Err(err).Msg("failed to get quotation")
		return "", err
	}
	return r.Text, nil
}

func (q *quotations) GetAllType() ([]string, error) {

	// r := models.Quotation{}
	return nil, nil
}
