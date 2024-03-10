package usecase

//
//import (
//	"context"
//	"errors"
//	"my-app/module/user/domain"
//	"testing"
//)
//
//type mockHashes struct{}
//
//func (mockHashes) RandomStr(length int) (string, error) {
//	return "abcd", nil
//}
//
//func (mockHashes) HashPassword(salt, password string) (string, error) {
//	return "iahweiufwef", nil
//}
//
//type mockUserRepo struct{}
//
//func (mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
//	if email == "existed@gmail.com" {
//		return &domain.User{}, nil
//	}
//
//	if email == "error@gmail.com" {
//		return nil, errors.New("cannot get record")
//	}
//
//	return &domain.User{}, nil
//}
//
//func (mockUserRepo) Create(ctx context.Context, data *domain.User) error {
//	return nil
//}
//
//func TestUseCase_Register(t *testing.T) {
//	uc := NewUseCase(mockUserRepo{}, mockHashes{})
//
//	type testData struct {
//		Input    EmailPasswordRegistrationDTO
//		Expected error
//	}
//
//	table := []testData{
//		{
//			Input: EmailPasswordRegistrationDTO{
//				FirstName: "Thien",
//				LastName:  "Trương Công",
//				Email:     "existed@gmail.com",
//				Password:  "123456",
//			},
//			Expected: domain.ErrEmailHasExisted,
//		},
//		{
//			Input: EmailPasswordRegistrationDTO{
//				FirstName: "Thien",
//				LastName:  "Trương Công",
//				Email:     "error@gmail.com",
//				Password:  "123456",
//			},
//			Expected: errors.New("cannot get record"),
//		},
//	}
//
//	for i := range table {
//		actualEror := uc.Register(context.Background(), table[i].Input)
//
//		if actualEror.Error() != table[i].Expected.Error() {
//			t.Errorf("Register failed. Expected %s, but actual is %s", table[i].Expected.Error(), actualEror.Error())
//		}
//	}
//
//}
