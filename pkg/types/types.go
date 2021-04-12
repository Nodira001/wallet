package types

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
