package order

const (
	Registered = "REGISTERED" // accrual not calculated yet
	Invalid    = "INVALID"    // registered, accrual will not be calculated. final
	Processing = "PROCESSING"
	Processed  = "PROCESSED" // final
)
