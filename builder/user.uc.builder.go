package builder

import (
	"gorm.io/gorm"
	"my-app/common"
	"my-app/module/user/infras/repository"
	"my-app/module/user/usecase"
)

type simpleBuilder struct {
	db *gorm.DB
	tp usecase.TokenProvider
}

func NewSimpleBuilder(db *gorm.DB, tp usecase.TokenProvider) simpleBuilder {
	return simpleBuilder{db: db, tp: tp}
}

func (s simpleBuilder) BuildUserQueryRepo() usecase.UserQueryRepository {
	return repository.NewUserRepo(s.db)
}

func (s simpleBuilder) BuildUserCmdRepo() usecase.UserCommandRepository {
	return repository.NewUserRepo(s.db)
}

func (simpleBuilder) BuildHashes() usecase.Hashes {
	return &common.Hasher{}
}

func (s simpleBuilder) BuildTokenProvider() usecase.TokenProvider {
	return s.tp
}

func (s simpleBuilder) BuildSessionQueryRepo() usecase.SessionQueryRepository {
	return repository.NewSessionMySQLRepo(s.db)
}

func (s simpleBuilder) BuildSessionCmdRepo() usecase.SessionCommandRepository {
	return repository.NewSessionMySQLRepo(s.db)
}
