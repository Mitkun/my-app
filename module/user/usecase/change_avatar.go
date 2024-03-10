package usecase

import (
	"context"
	"github.com/viettranx/service-context/core"
	"my-app/common"
	"my-app/module/user/domain"
)

type changeAvtUC struct {
	userQueryRepo UserQueryRepository
	userCmdRepo   UserCommandRepository
	imgRepo       ImageRepository
}

func NewChangeAvtUC(userQueryRepo UserQueryRepository, userCmdRepo UserCommandRepository, imgRepo ImageRepository) *changeAvtUC {
	return &changeAvtUC{userQueryRepo: userQueryRepo, userCmdRepo: userCmdRepo, imgRepo: imgRepo}
}

func (uc *changeAvtUC) ChangeAvatar(ctx context.Context, dto SingleImageDTO) error {
	userEntity, err := uc.userQueryRepo.Find(ctx, dto.Requester.UserId())
	if err != nil {
		return core.ErrBadRequest.WithError(domain.ErrCannotChangeAvatar.Error()).WithDebug(err.Error())
	}

	img, err := uc.imgRepo.Find(ctx, dto.ImageId)
	if err != nil {
		return core.ErrBadRequest.WithError(domain.ErrCannotChangeAvatar.Error()).WithDebug(err.Error())
	}

	if err := userEntity.ChangeAvatar(img.FileName); err != nil {
		return core.ErrBadRequest.WithError(domain.ErrCannotChangeAvatar.Error()).WithDebug(err.Error())
	}

	if err := uc.userCmdRepo.Update(ctx, userEntity); err != nil {
		return core.ErrBadRequest.WithError(domain.ErrCannotChangeAvatar.Error()).WithDebug(err.Error())
	}

	go func() {
		defer common.Recover()
		//TODO retry if fail
		_ = uc.imgRepo.SetImageStatusActivated(ctx, dto.ImageId)
	}()

	return nil
}