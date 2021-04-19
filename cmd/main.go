package main

import "github.com/Nodira001/wallet/pkg/wallet"

func main() {
	s := &wallet.Service{}
	acc, _ := s.RegisterAccount("+992004403883")
	s.Deposit(acc.ID, 10_000_00)
	s.ExportToFile("data/12345.txt")
}
