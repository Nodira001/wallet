package wallet

import (
	"reflect"
	"testing"
)

func TestService_AddAccountByID_success(t *testing.T) {
	svc := &Service{}

	_, err := svc.RegistAccount("+992900880306")
	if err != nil {
		t.Errorf("RegistAccount error %v,", err)

	}

}

func TestService_AddAccountByID_fail(t *testing.T) {
	svc := &Service{}

	_, _ = svc.RegistAccount("+992900880306")
	_, err := svc.RegistAccount("+992900880306")

	if err != ErrPhoneRegistered {
		t.Errorf("RegistAccount error %v,", err)

	}
}
func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	newAccount, _ := svc.RegistAccount("+992900800306")
	foundAccount, err := svc.FindAccountByID(newAccount.ID)
	if err != nil {
		t.Errorf("FindAccountByID error %v,", err)
	}
	if !reflect.DeepEqual(newAccount, foundAccount) {
		t.Error("error")
	}

}
func TestService_FindAccountByID_fail(t *testing.T) {
	svc := &Service{}
	_, _ = svc.RegistAccount("+992900800306")
	_, err := svc.FindAccountByID(10)
	if err != ErrAccountNotFound {
		t.Errorf("FindAccountByID error %v,", err)
	}

}
