package usecase

import (
	"context"
	"errors"
	"my-app/common"
	"my-app/module/user/domain"
)

type UseCase interface {
	Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error
	LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error)
}

type Hashes interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
	CompareHashPassword(hashedPassword, salt, password string) bool
}

type TokenProvider interface {
	IssueToken(ctx context.Context, id, sub string) (token string, err error)
	TokenExpireInSeconds() int
	RefreshExpireInSeconds() int
}

type useCase struct {
	uerRepo       UserRepository
	sessionRepo   SessionRepository
	hashes        Hashes
	tokenProvider TokenProvider
}

func NewUseCase(userRepo UserRepository, sessionRepo SessionRepository, hashes Hashes, tokenProvider TokenProvider) UseCase {
	return &useCase{uerRepo: userRepo, sessionRepo: sessionRepo, hashes: hashes, tokenProvider: tokenProvider}
}

func (uc *useCase) Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error {
	// 1. Find user by email
	// 1.1 Found: return error (email has existed)
	// 2. Generate salt
	// 3. hash password + salt
	// 4. Create user entity

	user, err := uc.uerRepo.FindByEmail(ctx, dto.Email)
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

	if err := uc.uerRepo.Create(ctx, userEntity); err != nil {
		return err
	}

	return nil
}

type UserRepository interface {
	//Find(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, data *domain.User) error
	//Update(ctx context.Context, data *domain.User) error
	//Delete(ctx context.Context, data *domain.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, data *domain.Session) error
}
