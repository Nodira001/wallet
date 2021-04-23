package wallet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Nodira001/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("payment is not enough")
var ErrAccountNotFond = errors.New("reject")
var ErrFavoriteNotFound = errors.New("payment is not enough")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) addAccount(data testAccount1) (*types.Account, []*types.Payment, error) {
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
func (s *Service) FindFavoritePayment(paymentID string) (*types.Favorite, error) {
	payment := types.Favorite{}
	for _, paymentf := range s.favorites {
		if paymentf.ID == paymentID {
			payment = *paymentf
			return &payment, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	newFavorite := types.Favorite{
		ID:        uuid.NewString(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, &newFavorite)
	return &newFavorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favoritePayment, err := s.FindFavoritePayment(favoriteID)
	if err != nil {
		return nil, err
	}
	acc, err := s.FindAccountByID(favoritePayment.AccountID)
	if err != nil {
		return nil, err
	}
	if acc.Balance < favoritePayment.Amount {
		return nil, ErrNotEnoughBalance
	}

	newPayment := types.Payment{
		ID:        uuid.NewString(),
		AccountID: favoritePayment.AccountID,
		Amount:    favoritePayment.Amount,
		Category:  favoritePayment.Category,
		Status:    types.PaymentStatusInProgress,
	}
	acc.Balance -= favoritePayment.Amount
	s.payments = append(s.payments, &newPayment)
	return &newPayment, nil
}
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer file.Close()
	for _, account := range s.accounts {
		row := fmt.Sprint(account.ID) + ";" + fmt.Sprint(account.Phone) + ";" + fmt.Sprint(account.Balance) + "|"
		_, err = file.Write([]byte(row))
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
	  if favorite.ID == favoriteID {
		return favorite, nil
	  }
	}
  
	return nil, ErrFavoriteNotFound
  }
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	buf := make([]byte, 1)
	content := make([]byte, 0)
	for {
		read, err := file.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := strings.Split(string(content), "|")
	var row []string
	for _, v := range data {
		row = append(row, strings.ReplaceAll(v, ";", " "))
	}
	for _, acc := range row {
		if len(acc) == 0 {
			continue
		}
		accountSplit := strings.Split(acc, " ")

		id, err := strconv.ParseInt(accountSplit[0], 10, 64)
		if err != nil {

			return err
		}
		balance, err := strconv.ParseInt(accountSplit[2], 10, 64)

		if err != nil {
			return err
		}
		account := &types.Account{ID: id, Balance: types.Money(balance), Phone: types.Phone(accountSplit[1])}
		s.accounts = append(s.accounts, account)

	}
	return nil
}
func (s *Service) Import(dir string) (importError error) {

	_, err := os.Stat(dir + "/accounts.dump")
	
	if err == nil {
	  accountsFile, err := os.Open(dir + "/accounts.dump")
	  if err != nil {
		return err
	  }
	  defer accountsFile.Close()
	  accountsReader := bufio.NewReader(accountsFile)
	  for {
		line, err := accountsReader.ReadString('\n')
		if err == io.EOF {
		  break
		}
		if err != nil {
		  return err
		}
		account := strings.Split(line, ";")
  
		accountID, err := strconv.ParseInt(account[0], 10, 64)
		if err != nil {
		  return err
		}
		accountPhone := account[1]
		accountBalance, err := strconv.ParseInt(strings.ReplaceAll(account[2], "\n", ""), 10, 64)
		if err != nil {
		  return err
		}
		accountBackUp := &types.Account{
		  ID:      accountID,
		  Phone:   types.Phone(accountPhone),
		  Balance: types.Money(accountBalance),
		}
		_, err = s.FindAccountByID(accountID)
		if err == ErrAccountNotFound {
		  s.accounts = append(s.accounts, accountBackUp)
		  s.nextAccountID = int64(len(s.accounts))
		}
  
	  }
	}
	_, err = os.Stat(dir + "/payments.dump")
	if err == nil {
	  paymentsFile, err := os.Open(dir + "/payments.dump")
	  if err != nil {
		return err
	  }
  
	  defer paymentsFile.Close()
	  paymentsReader := bufio.NewReader(paymentsFile)
	  for {
		line, err := paymentsReader.ReadString('\n')
		if err == io.EOF {
		  break
		}
		if err != nil {
		  return err
		}
		payment := strings.Split(line, ";")
  
		paymentID := payment[0]
		paymentAccountID, err := strconv.ParseInt(payment[1], 10, 64)
		if err != nil {
  
		  return err
		}
		paymentAmount, err := strconv.ParseInt(payment[2], 10, 64)
		if err != nil {
		  importError = err
		  return importError
		}
		paymentCategory := payment[3]
		paymentStatus := strings.ReplaceAll(payment[4], "\n", "")
		paymentBackUp := &types.Payment{
		  ID:        paymentID,
		  AccountID: paymentAccountID,
		  Amount:    types.Money(paymentAmount),
		  Category:  types.PaymentCategory(paymentCategory),
		  Status:    types.PaymentStatus(paymentStatus),
		}
		_, err = s.FindPaymentByID(paymentID)
		if err == ErrPaymentNotFound {
		  s.payments = append(s.payments, paymentBackUp)
		}
  
	  }
	}
  
	_, err = os.Stat(dir + "/favorites.dump")
	if err == nil {
	  favoritesFile, err := os.Open(dir + "/favorites.dump")
	  if err != nil {
		return err
	  }
	  defer favoritesFile.Close()
	  favoritesReader := bufio.NewReader(favoritesFile)
  
	  for {
		line, err := favoritesReader.ReadString('\n')
		if err == io.EOF {
		  break
		}
		if err != nil {
		  return err
		}
		favorite := strings.Split(line, ";")
		favoriteID := favorite[0]
		favoriteAccountID, err := strconv.ParseInt(favorite[1], 10, 64)
		if err != nil {
		  return err
		}
		favoriteAmount, err := strconv.ParseInt(favorite[2], 10, 64)
		if err != nil {
		  return err
		}
		favoriteName := favorite[3]
		favoriteCategory := strings.ReplaceAll(favorite[4], "\n", "")
		favoriteBackUp := &types.Favorite{
		  ID:        favoriteID,
		  AccountID: favoriteAccountID,
		  Amount:    types.Money(favoriteAmount),
		  Name:      favoriteName,
		  Category:  types.PaymentCategory(favoriteCategory),
		}
		_, err = s.FindFavoriteByID(favoriteID)
		if err == ErrFavoriteNotFound {
		  s.favorites = append(s.favorites, favoriteBackUp)
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
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
  if len(payments) > 0 {
    if len(payments) <= records {
      file, _ := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
      defer file.Close()

      var str string
      for _, v := range payments {
        str += fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
      }
      file.WriteString(str)
      return nil
    }  

      var row string
      k := 0
      count := 1
      var file *os.File
      for _, v := range payments {
        if k == 0 {
          file, _ = os.OpenFile(dir+"/payments"+fmt.Sprint(count)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
        }
        k++
        row = fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
        _, err := file.WriteString(row)
        if err != nil {
          return err
        }
        if k == records {
          row = ""
          count++
          k = 0
          file.Close()
        }
      }

   
  }

  return nil
}