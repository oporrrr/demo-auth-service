package repository

import (
	"demo-auth-center/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Upsert(user *entity.User) error {
	return r.db.
		Where(entity.User{AccountID: user.AccountID}).
		Assign(entity.User{
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			CountryCode:   user.CountryCode,
			PrefixName:    user.PrefixName,
			Gender:        user.Gender,
			DateOfBirth:   user.DateOfBirth,
			AccountStatus: user.AccountStatus,
			CisNumber:     user.CisNumber,
		}).
		FirstOrCreate(user).Error
}

func (r *UserRepository) FindByAccountID(accountID string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateRole(accountID, role string) error {
	return r.db.Model(&entity.User{}).
		Where("account_id = ?", accountID).
		Update("role", role).Error
}

