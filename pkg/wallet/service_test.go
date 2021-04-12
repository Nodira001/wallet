package wallet

import (
	"reflect"
	"testing"

	"github.com/Nodira001/wallet/pkg/types"
)

func TestService_AddAccountByID_success(t *testing.T) {
	svc := &Service{}

	_, err := svc.RegisterAccount("+992900880306")
	if err != nil {
		t.Errorf("RegistAccount error %v,", err)

	}

}

func TestService_AddAccountByID_fail(t *testing.T) {
	svc := &Service{}

	_, _ = svc.RegisterAccount("+992900880306")
	_, err := svc.RegisterAccount("+992900880306")

	if err != ErrPhoneRegistered {
		t.Errorf("RegistAccount error %v,", err)

	}
}
func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	newAccount, _ := svc.RegisterAccount("+992900800306")
	foundAccount, err := svc.FindAccountByID(newAccount.ID)
	if err != nil {
		t.Errorf("FindAccountByID error %v,", err)
	}
	if !reflect.DeepEqual(newAccount, foundAccount) {
		t.Error("error")
	}

}
func TestService_Reject_succes(t *testing.T) {
	svc := &Service{}

	acc, err := svc.RegisterAccount("+992900800306")
	if err != nil {
		t.Error("test")
	}
	svc.payments = append(svc.payments, &types.Payment{
		ID:        "10",
		AccountID: acc.ID,
		Amount:    1000,
		Category:  "food",
		Status:    types.PaymentStatusInProgress})
	err = svc.Reject("10")
	if err != nil {
		t.Error("err 54 ", err, svc)
	}
}
func TestService_Reject_fail(t *testing.T) {
	svc := &Service{}

	acc, err := svc.RegisterAccount("+992900800306")
	if err != nil {
		t.Error("test")
	}
	svc.payments = append(svc.payments, &types.Payment{
		ID:        "10",
		AccountID: acc.ID,
		Amount:    1000,
		Category:  "food",
		Status:    types.PaymentStatusInProgress})
	err = svc.Reject("0")
	if err != ErrPaymentNotFound {
		t.Error("err 54 ", err, svc)
	}
}

func TestService_Repeat_succes(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error("asd")
		return
	}
	pament := payments[0]
	_, err = s.Repeat(pament.ID)
	if err != nil {
		t.Error("asd")
		return
	}
}
func TestService_Repeat_fail(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defultTestAccount)

	if err != nil {
		t.Error("asd")
		return
	}

	pament := payments[0]
	
	pament.Amount += 9_000_00
	
	_, err = s.Repeat(pament.ID)
	
	if err == nil {
		t.Error("asd")
		return
	}

}
