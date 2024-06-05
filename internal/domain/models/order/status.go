package order

//const (
//	OrderStatusRegistered = "REGISTERED" // accrual not calculated yet
//	OrderStatusInvalid    = "INVALID"    // registered, accrual will not be calculated. final
//	OrderStatusProcessing = "PROCESSING"
//	OrderStatusProcessed  = "PROCESSED" // final
//)

const (
	Registered = "REGISTERED" // accrual not calculated yet
	Invalid    = "INVALID"    // registered, accrual will not be calculated. final
	Processing = "PROCESSING"
	Processed  = "PROCESSED" // final
)

type Status string
