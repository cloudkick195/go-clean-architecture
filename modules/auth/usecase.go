package auth

import (
	"context"
	"time"

	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/modules/member"
)

type IUsecase interface {
	Register(ctx context.Context, input *member.CreateInput) error
	Login(ctx context.Context, input *LoginInput) (*commons.Token, error)
}

type Usecase struct {
	member member.IUsecase
}

func NewUsecase(member member.IUsecase) IUsecase {
	return &Usecase{
		member: member,
	}
}
func (a *Usecase) Register(ctx context.Context, input *member.CreateInput) error {
	return a.member.Create(ctx, input)
}

func (a *Usecase) Login(ctx context.Context, input *LoginInput) (*commons.Token, error) {
	result, _ := a.member.Get(ctx, &models.Member{SupId: input.SupId})
	if result == nil || !result.ComparePassword(input.Password) {
		return nil, member.ErrIdOrPasswordInvalid
	}

	if result.Status != models.MemberStatusActive {
		return nil, member.ErrMemberBlocked
	}

	payload := commons.TokenPayload{
		Id: result.ID,
	}

	jwtProvider := commons.NewJwtProvider()
	token, err := jwtProvider.Generate(payload, 5*time.Hour)
	if err != nil {
		return nil, err
	}
	updateM := map[string]interface{}{
		"Token":       token.Token,
		"TokenExpiry": token.Expiry,
	}
	if err := a.member.Update(ctx, nil, result, updateM); err != nil {
		return nil, err
	}
	return token, nil

}
