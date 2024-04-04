package repo

import (
	"context"
	"fmt"
	"strconv"

	"contact-management-service/internal/contacts/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) HealthCheck(ctx context.Context) error {
	alive := false
	res := r.db.WithContext(ctx).Raw("SELECT true").Scan(&alive)
	if res.Error != nil {
		log.Debug().Msgf("database is not reachable: %v", res.Error)
		return fmt.Errorf("db.Raw: %w", res.Error)
	}
	return nil
}

func (r *Repo) RunTransaction(ctx context.Context, fn func(tx *gorm.DB) error, onCommitFail func() error) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return onCommitFail()
	}
	return nil
}

func (r *Repo) StoreDataInDB(ctx context.Context, tx *gorm.DB, user *model.Contacts) error {
	err := tx.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("tx.Create failed to insert user data: %v", err)
	}

	return nil
}

func (r *Repo) UpdateDataInDB(ctx context.Context, tx *gorm.DB, userId string, user *model.Contacts) error {
	var contact model.Contacts
	if len(userId) > 0 {
		newUserId, _ := strconv.Atoi(userId)
		if err := tx.WithContext(ctx).First(&contact, newUserId).Error; err != nil {
			return fmt.Errorf("tx.First contact not found: %v", err)
		} else {
			if err := tx.WithContext(ctx).Where("id=?", userId).Updates(&user).Error; err != nil {
				return fmt.Errorf("tx.Save failed to update user contact info: %v", err)
			}
		}
	}

	return nil
}

func (r *Repo) GetContactByID(ctx context.Context, id string) (data []model.Contacts, err error) {
	var userData []model.Contacts

	if len(id) > 0 {
		newUserId, err := strconv.Atoi(id)
		if err != nil {
			return userData, fmt.Errorf("errr converting id string to int: %w", err)
		}

		err = r.db.Raw("SELECT * FROM contacts WHERE id=?", newUserId).Scan(&userData).Error
	} else {
		err = r.db.Raw("SELECT * FROM contacts ORDER BY id ASC").Scan(&userData).Error
	}

	if err != nil {
		return userData, err
	}

	return userData, nil
}

func (r *Repo) DeleteContact(ctx context.Context, tx *gorm.DB, userId string) error {
	var contact model.Contacts
	if len(userId) > 0 {
		newUserId, _ := strconv.Atoi(userId)
		if err := tx.WithContext(ctx).First(&contact, newUserId).Error; err != nil {
			return fmt.Errorf("tx.First contact not found: %v", err)
		} else {
			err = r.db.WithContext(ctx).Where("id = ?", userId).Delete(&model.Contacts{}).Error
		}
	}

	return nil
}
