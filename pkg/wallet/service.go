package wallet

import (
	"errors"
	"fmt"

	"github.com/Nodira001/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("payment is not enough")
var ErrAccountNotFond = errors.New("reject")

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
	s.accounts = append(s.accounts, newAccount)
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

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments {
		if paymentID == pay.ID {
			payment = pay

		}

	}
	if payment == nil {
		return nil, ErrPaymentNotFound

	}
	return payment, nil

}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound

	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}
 

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return err

	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err

	}
	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFond
	}

	account.Balance += amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	paymentNew, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return paymentNew, nil

}

type testService1 struct {
	*Service
}

type testAccount1 struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defultTestAccount1 = testAccount1{
	phone:   "+992900880306",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *Service) addAccount(data  testAccount1) (*types.Account, []*types.Payment, error) {
	// тестируем пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error =%v", err)
	}
	// пополняем его счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error =%v", err)
	}
	// 	выполяняем  платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работает через индекс, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't register account, error =%v", err)
		}
	}
	return account, payments, nil
}
