package wallet

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
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
	_, payments, err := s.addAccount(defultTestAccount1)
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
	_, payments, err := s.addAccount(defultTestAccount1)

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
func TestService_FavoritePayment_succes(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defultTestAccount1)

	if err != nil {
		t.Error("asd")
		return
	}
	_, err = s.FavoritePayment(payments[0].ID, "Jenya")
	if err != nil {
		t.Error("FavoriPayment(): error")
	}
}
func TestService_PayFromFavorite_succes(t *testing.T) {
	s := &Service{}
	acc, payments, err := s.addAccount(defultTestAccount1)

	if err != nil {
		t.Error("asd")
		return
	}

	fav, err := s.FavoritePayment(payments[0].ID, "Jenya")
	if err != nil {
		t.Error("FavoriPayment(): error", acc.Balance)
	}
	_, err = s.PayFromFavorite(fav.ID)
	if err != nil {
		t.Error("sdfgh")
		return

	}

}
func (s *Service) Export(dir string) error {
  
	if len(s.accounts) > 0 {
  
	  accountsFile, err := os.Create(dir + "/accounts.dump")
  
	  if err != nil {
		log.Println(err)
		return err
	  }
  
	  defer func() {
		if accErr := accountsFile.Close(); accErr != nil {
		  log.Print(accErr)
		  return
		}
	  }()
  
	  for _, account := range s.accounts {
  
		accountsRow := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + string(strconv.FormatInt(int64(account.Balance), 10)) + "\n"
		_, err = accountsFile.Write([]byte(accountsRow))
		if err != nil {
		  log.Print(err)
		  return err
		}
	  }
	}
	if len(s.payments) > 0 {
	  paymentsFile, err := os.Create(dir + "/payments.dump")
	  if err != nil {
		log.Println(err)
		return err
	  }
  
	  defer func() {
		if payErr := paymentsFile.Close(); payErr != nil {
		  log.Print(payErr)
		  return
		}
	  }()
  
	  for _, payment := range s.payments {
  
		paymentsRow := payment.ID + ";" + strconv.FormatInt(payment.AccountID, 10) + ";" + strconv.FormatInt(int64(payment.Amount), 10) + ";" + string(payment.Category) + ";" + string(payment.Status) + "\n"
		_, err = paymentsFile.Write([]byte(paymentsRow))
		if err != nil {
		  log.Print(err)
		  return err
		}
	  }
	}
  
	if len(s.favorites) > 0 {
	  favoritesFile, err := os.Create(dir + "/favorites.dump")
  
	  if err != nil {
		log.Println(err)
		return err
	  }
  
	  defer func() {
		if favErr := favoritesFile.Close(); favErr != nil {
		  log.Print(favErr)
		  return
		}
	  }()
  
	  for _, favorite := range s.favorites {
  
		favoriteRow := favorite.ID + ";" + strconv.FormatInt(int64(favorite.AccountID), 10) + ";" + strconv.FormatInt(int64(favorite.Amount), 10) + ";" + favorite.Name + ";" + string(favorite.Category) + "\n"
		_, err = favoritesFile.Write([]byte(favoriteRow))
		if err != nil {
		  log.Print(err)
		  return err
		}
	  }
	}
	fmt.Println("nextAccountID", s.nextAccountID, "accounts->>", len(s.accounts), "payments->>", len(s.payments), "favorites->>", len(s.favorites))
	fmt.Println("start")
	for _, v := range s.accounts {
	  fmt.Println(v)
	}
	for _, v := range s.payments {
	  fmt.Println(v)
	}
	for _, v := range s.favorites {
	  fmt.Println(v)
	}
	fmt.Println("stop")
	return nil
  }
func TestService_FullExport(t *testing.T) {
	s := &Service{}
  
	acc, err := s.RegisterAccount("+992004403883")
	if err != nil {
	  fmt.Print(err)
	  return
	}
  
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
	  fmt.Print(err)
	  return
	}
  
	payment, err := s.Pay(acc.ID, 10_000, "auto")
	if err != nil {
	  fmt.Print(err)
	  return
	}
  
	_, err = s.FavoritePayment(payment.ID, "Auto")
	if err != nil {
	  fmt.Print(err)
	  return
	}
  
	err = s.Export("data")
	if err != nil {
	  fmt.Print(err)
	  return
	}
  }