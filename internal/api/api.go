package api

import (
	"net/http"
	"payment-gateway/db"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/services"
	"payment-gateway/internal/services/gateway"
	"payment-gateway/internal/services/tx"

	"github.com/gorilla/mux"
)

type API struct {
	Router *mux.Router
	db     db.Idb
	svc    services.Service
}

func New(dbInst db.Idb) *API {
	return &API{db: dbInst, Router: mux.NewRouter()}
}

func (a *API) SetupServices(kafkaProducer kafka.IProducer) {
	a.svc.ISvcGateway = gateway.NewSvcGateway(a.db)
	a.svc.ISvcTx = tx.NewSvcTx(a.db, kafkaProducer)
}

func (a *API) SetupRoutes() {
	a.Router.Handle("/deposit", http.HandlerFunc(a.DepositHandler)).Methods("POST")
	a.Router.Handle("/withdrawal", http.HandlerFunc(a.WithdrawalHandler)).Methods("POST")
	a.Router.Handle("/call_back", http.HandlerFunc(a.CallBackHandler)).Methods("GET")
}
