package usecase

import (
	"context"
	"errors"
	"my-app/common"
	"my-app/module/user/domain"
)

type registerUC struct {
	userQueryRepo UserQueryRepository
	userCmdRepo   UserCommandRepository
	hashes        Hashes
}

func NewRegisterUC(userQueryRepo UserQueryRepository, userCmdRepo UserCommandRepository, hashes Hashes) *registerUC {
	return &registerUC{userQueryRepo: userQueryRepo, userCmdRepo: userCmdRepo, hashes: hashes}
}

func (uc *registerUC) Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error {
	// 1. Find user by email
	// 1.1 Found: return error (email has existed)
	// 2. Generate salt
	// 3. hash password + salt
	// 4. Create user entity

	user, err := uc.userQueryRepo.FindByEmail(ctx, dto.Email)
	if user != nil {
		return domain.ErrEmailHasExisted
	}

	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		return err
	}

	salt, err := uc.hashes.RandomStr(30)
	if err != nil {
		return err
	}

	hashedPassword, err := uc.hashes.HashPassword(salt, dto.Password)
	if err != nil {
		return err
	}

	userEntity, err := domain.NewUser(
		common.GenUUID(),
		dto.FirstName,
		dto.LastName,
		dto.Email,
		hashedPassword,
		salt,
		domain.RoleUser,
	)
	if err != nil {
		return err
	}

	if err := uc.userCmdRepo.Create(ctx, userEntity); err != nil {
		return err
	}

	return nil
}
