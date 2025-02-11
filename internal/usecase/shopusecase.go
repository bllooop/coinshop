package usecase

import (
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

func (s *ShopUsecase) SendCoin(input domain.Transactions) (int, error) {
	return s.repo.SendCoin(input)
}

func (s *ShopUsecase) BuyItem(userid int, name string) (int, error) {
	return s.repo.BuyItem(userid, name)
}
