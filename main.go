package main

import (
	"context"
	"fmt"
	"github.com/devpayments/common/errors"
	"github.com/devpayments/core/payments"
	"github.com/devpayments/corleone/config"
	"github.com/devpayments/corleone/datastore"
	"github.com/devpayments/corleone/datastore/db"
	_ "github.com/lib/pq"
)

func main() {
	initCtx := context.TODO()

	dbConfig := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		User:     "remi",
		Password: "root1234",
		SSLMode:  "disable",
		Name:     "payments",
		Port:     "5432",
	}
	d := db.New(dbConfig)
	dbCon, err := d.GetDbConnection()
	if err != nil {
		panic(err)
	}
	defer dbCon.Close()
	store := datastore.NewStore(*dbCon)

	paymentService := payments.NewPaymentService(store.Payments, store.Transactions)
	payment, err := paymentService.Initiate(initCtx)
	errors.PanicIfNecessary(err)

	//err = paymentService.CompleteAuthorization(initCtx, uuid.MustParse(payment.ID), "")
	//errors.PanicIfNecessary(err)
	//
	//err = paymentService.Complete(initCtx, uuid.MustParse(payment.ID))
	//errors.PanicIfNecessary(err)

	fmt.Printf("%+v\n", payment)
	fmt.Printf("%+v\n", err)

	//fmt.Println(res)
}
