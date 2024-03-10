package usecase

import (
	"context"
	"github.com/google/uuid"
	"my-app/common"
	"my-app/module/user/domain"
)

type UseCase interface {
	Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error
	LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponseDTO, error)
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
	*loginEmailPasswordUC
	*registerUC
	*refreshTokenUC
}

type Builder interface {
	BuildUserQueryRepo() UserQueryRepository
	BuildUserCmdRepo() UserCommandRepository
	BuildHashes() Hashes
	BuildTokenProvider() TokenProvider
	BuildSessionQueryRepo() SessionQueryRepository
	BuildSessionCmdRepo() SessionCommandRepository
	BuildSessionRepo() SessionRepository
}

func UseCaseWithBuilder(b Builder) UseCase {
	return &useCase{
		loginEmailPasswordUC: NewLoginEmailPasswordUC(b.BuildUserQueryRepo(), b.BuildSessionCmdRepo(), b.BuildTokenProvider(), b.BuildHashes()),
		registerUC:           NewRegisterUC(b.BuildUserQueryRepo(), b.BuildUserCmdRepo(), b.BuildHashes()),
		refreshTokenUC:       NewRefreshTokenUC(b.BuildUserQueryRepo(), b.BuildSessionRepo(), b.BuildTokenProvider(), b.BuildHashes()),
	}
}

//func NewUseCase(userRepo UserRepository, sessionRepo SessionRepository, hashes Hashes, tokenProvider TokenProvider) UseCase {
//	return &useCase{
//		loginEmailPasswordUC: NewLoginEmailPasswordUC(userRepo, sessionRepo, tokenProvider, hashes),
//		changeAvtUC:           NewChangeAvtUC(userRepo, userRepo, hashes),
//	}
//}

type UserRepository interface {
	UserQueryRepository
	UserCommandRepository
}

type UserQueryRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Find(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type UserCommandRepository interface {
	Create(ctx context.Context, data *domain.User) error
	Update(ctx context.Context, data *domain.User) error
}

type SessionRepository interface {
	SessionQueryRepository
	SessionCommandRepository
}

type SessionQueryRepository interface {
	Find(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	FindByRefreshToken(ctx context.Context, rt string) (*domain.Session, error)
}

type SessionCommandRepository interface {
	Create(ctx context.Context, data *domain.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ImageRepository interface {
	Find(ctx context.Context, id uuid.UUID) (*common.Image, error)
	SetImageStatusActivated(ctx context.Context, id uuid.UUID) error
}
