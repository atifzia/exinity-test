package api

import (
	"net/http"
	"payment-gateway/internal/models/request"
	"payment-gateway/internal/util"
	"strconv"
	"strings"
)

// DepositHandler handles deposit requests (feel free to update how user is passed to the request)
// Sample Request (POST /deposit):
//
//	{
//	    "amount": 100.00,
//	    "user_id": 1,
//	    "currency": "EUR"
//	}
func (a *API) DepositHandler(w http.ResponseWriter, r *http.Request) {
	var req request.Transaction

	// decode JSON payload
	if err := util.DecodeRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process the deposit request
	response, err := a.svc.ISvcTx.ProcessTransaction(req, a.svc.ISvcGateway, "deposit")
	if err != nil {
		if strings.Contains(err.Error(), "no gateways available for the specified country") {
			http.Error(w, err.Error(), http.StatusGatewayTimeout)

			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode and send the response
	util.SendEncodedResponse(w, response, http.StatusOK)
}

// WithdrawalHandler handles withdrawal requests (feel free to update how user is passed to the request)
// Sample Request (POST /deposit):
//
//	{
//	    "amount": 100.00,
//	    "user_id": 1,
//	    "currency": "EUR"
//	}
func (a *API) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	var req request.Transaction

	// decode JSON payload
	if err := util.DecodeRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process the withdrawal request
	response, err := a.svc.ISvcTx.ProcessTransaction(req, a.svc.ISvcGateway, "withdrawal")
	if err != nil {
		if strings.Contains(err.Error(), "no gateways available for the specified country") {
			http.Error(w, err.Error(), http.StatusGatewayTimeout)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode and send the response
	util.SendEncodedResponse(w, response, http.StatusOK)
}

// CallBackHandler handles updating of transaction status via callback
// Sample Request (GET /call_back?tx_id=101&status=completed)
func (a *API) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	txIdStr := r.URL.Query().Get("tx_id")
	txIdInt, err := strconv.ParseInt(txIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}
	status := r.URL.Query().Get("status")

	a.svc.ISvcTx.ProcessCallBack(txIdInt, status)

}
