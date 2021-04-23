package types

import "time"

// Money предоставляет собой денежную сумму в минимальных единицах (центы, копецки, дирамы и т.д.).
type Money int64

// PaymentCategory представляет собой категорию, в которой был совершен платеж (авто, аптеки, рестораны и т.д.).
type PaymentCategory string

// PaymentStatus представляем собой статус платежа
type PaymentStatus string

// Предопределенные  статусы платежей
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment предоставляет информацию о платеже
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Phone string

// Account предоставляет информацию о счете  пользователя
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}

//Информация о файле описывает файл и возвращается Stat и Lstat
type FileInfo interface {
	Name() string       //базовое имя файла
	Size() int64        //длина в байтах для обычных файлов, зависит от системы для остальных
	Mode() FileInfo     //биты файловых модов
	ModTime() time.Time //время модификации
	IsDir() bool
	Sys() interface{}
}

