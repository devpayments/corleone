module github.com/devpayments/corleone

go 1.19

replace github.com/devpayments/common => ../common

replace github.com/devpayments/core => ../core

require (
	github.com/devpayments/core v0.0.0-20221101053241-ca462a5f42e7
	github.com/doug-martin/goqu/v9 v9.18.0
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.7
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/devpayments/common v0.0.0-20221101053312-1a187d1b6de1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
