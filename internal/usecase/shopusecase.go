package usecase

import (
	"time"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/bllooop/coinshop/internal/repository"
)

type ShopUsecase struct {
	repo repository.Shop
}

func NewShopUsecase(repo *repository.Repository) *ShopUsecase {
	return &ShopUsecase{
		repo: repo,
	}
}

func (s *ShopUsecase) SendCoin(userid int, input domain.Transactions) (int, error) {
	input.Source = &userid
	timestamp := time.Now()
	input.Timestamp = &timestamp
	return s.repo.SendCoin(input)
}

func (s *ShopUsecase) BuyItem(userid int, name string) (int, error) {
	return s.repo.BuyItem(userid, name)
}

func (s *ShopUsecase) GetUserSummary(userID int) (*domain.UserSummary, error) {
	return s.repo.GetUserSummary(userID)
}
