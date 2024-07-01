package order

const (
	New        = "NEW"        // created in local db, not registered in accrual
	Registered = "REGISTERED" // accrual not calculated yet
	Invalid    = "INVALID"    // registered, accrual will not be calculated. final
	Processing = "PROCESSING"
	Processed  = "PROCESSED" // final
)
