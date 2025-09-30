package service

import "github.com/KatenkaKet/wallet/pkg/repository"

type Wallet interface {
}

type Service struct {
	Wallet
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
