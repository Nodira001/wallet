package wallet

import (
	"errors"
	"github.com/Nodira001/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}

	}
	s.nextAccountID++
	newAccount := &types.Account{
		ID:      s.nextAccountID,
		Balance: 0,
		Phone:   phone,
	}
	s.accounts= append(s.accounts, newAccount)
	return newAccount, nil

}
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
		}

	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}
