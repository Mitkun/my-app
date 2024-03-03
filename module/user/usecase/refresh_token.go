package usecase

import (
	"context"
	"errors"
	"my-app/common"
	"my-app/module/user/domain"
	"time"
)

type refreshTokenUC struct {
	userRepo      UserQueryRepository
	sessionRepo   SessionRepository
	tokenProvider TokenProvider
	hashes        Hashes
}

func NewRefreshTokenUC(userRepo UserQueryRepository, sessionRepo SessionRepository, tokenProvider TokenProvider, hashes Hashes) *refreshTokenUC {
	return &refreshTokenUC{userRepo: userRepo, sessionRepo: sessionRepo, tokenProvider: tokenProvider, hashes: hashes}
}

func (uc *refreshTokenUC) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponseDTO, error) {
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if session.RefreshExpAt().Before(time.Now().UTC()) {
		return nil, errors.New("refresh token has expired")
	}

	user, err := uc.userRepo.Find(ctx, session.UserId())
	if err != nil {
		return nil, err
	}

	if user.Status() == "banned" {
		return nil, errors.New("user has been banned")
	}

	userId := user.Id()
	sessionId := common.GenUUID()

	// 3. Gen JWT
	accessToken, err := uc.tokenProvider.IssueToken(ctx, sessionId.String(), userId.String())

	if err != nil {
		return nil, err
	}

	// 4. Insert session into DB
	newRefreshToken, _ := uc.hashes.RandomStr(16)
	tokenExpAt := time.Now().UTC().Add(time.Second * time.Duration(uc.tokenProvider.TokenExpireInSeconds()))
	refreshExpAt := time.Now().UTC().Add(time.Second * time.Duration(uc.tokenProvider.RefreshExpireInSeconds()))

	newSession := domain.NewSession(sessionId, userId, newRefreshToken, tokenExpAt, refreshExpAt)

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	go func() {
		_ = uc.sessionRepo.Delete(ctx, session.Id())
	}()

	// 5. Return token response dto

	return &TokenResponseDTO{
		AccessToken:       accessToken,
		AccessTokenExpIn:  uc.tokenProvider.TokenExpireInSeconds(),
		RefreshToken:      newRefreshToken,
		RefreshTokenExpIn: uc.tokenProvider.RefreshExpireInSeconds(),
	}, nil
}
