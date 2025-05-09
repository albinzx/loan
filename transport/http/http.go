package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/albinzx/loan/entity"
	svc "github.com/albinzx/loan/service"
	"github.com/albinzx/loan/transport/model"
	"github.com/julienschmidt/httprouter"
)

type LoanHTTPTransport struct {
	service svc.LoanEngine
}

// New creates new loan http transport
func New(service svc.LoanEngine) *LoanHTTPTransport {
	return &LoanHTTPTransport{service: service}
}

// Handler returns http router and handler
func (h *LoanHTTPTransport) Handler(path string) http.Handler {
	router := httprouter.New()

	router.POST(fmt.Sprintf("%s/loans", path), h.Create)
	router.GET(fmt.Sprintf("%s/loans/:id", path), h.Get)
	router.GET(fmt.Sprintf("%s/loans", path), h.GetByStateOrBorrower)
	router.PATCH(fmt.Sprintf("%s/loans/:id/approve", path), h.Approve)
	router.PATCH(fmt.Sprintf("%s/loans/:id/invest", path), h.Invest)
	router.PATCH(fmt.Sprintf("%s/loans/:id/disburse", path), h.Disburse)
	router.GET(fmt.Sprintf("%s/investments/:id", path), h.GetByInvestor)

	return router
}

// Create is http handler for create loan
func (h *LoanHTTPTransport) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// decode request
	loan := &model.Loan{}
	if err := json.NewDecoder(r.Body).Decode(loan); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to create loan
	res, err := h.service.Create(r.Context(), loan.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response
	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// get path parameter
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to get loan
	loan, err := h.service.Get(r.Context(), loanID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// loan not found
	if loan == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// encode response
	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*loan)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Approve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// get path parameter
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decode request
	approval := &model.Approval{}
	if err := json.NewDecoder(r.Body).Decode(approval); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate mandatory fields
	if approval.Empty() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to approve loan
	res, err := h.service.Approve(r.Context(), loanID, approval.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// invalid state to be approved
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// encode response
	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Invest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// get path parameter
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decode request
	investment := &model.Investment{}
	if err := json.NewDecoder(r.Body).Decode(investment); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate mandatory fields
	if investment.Empty() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to invest loan
	res, err := h.service.Invest(r.Context(), loanID, investment.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// invalid state to be invested
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// encode response
	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Disburse(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// get path parameter
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decode request
	approval := &model.Approval{}
	if err := json.NewDecoder(r.Body).Decode(approval); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate mandatory fields
	if approval.Empty() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to disburse loan
	res, err := h.service.Disburse(r.Context(), loanID, approval.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// invalid state to be disbursed
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// encode response
	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) GetByInvestor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// get path parameter
	id := params.ByName("id")
	investorID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to get loan by investor id
	loanEntities, err := h.service.GetByInvestor(r.Context(), investorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response
	loans := make([]model.Loan, len(loanEntities))
	for i, l := range loanEntities {
		loans[i] = model.ToLoanModel(l)
	}

	if err = json.NewEncoder(w).Encode(loans); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) GetByStateOrBorrower(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get query parameters
	params := r.URL.Query()
	state := params.Get("state")
	borrower := params.Get("borrower")

	// validate parameter
	if state == "" && borrower == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate borrower id is integer
	var borrowerID int64
	var err error
	if borrower != "" {
		borrowerID, err = strconv.ParseInt(borrower, 10, 64)
		if err != nil {
			log.Printf("error while parsing parameter id, %v", err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// validate both parameter cannot be empty
	entState := entity.StateOf(state)
	if entState == nil && borrowerID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call service to get loans by state or borrower
	loanEntities, err := h.service.GetByStateOrBorrower(r.Context(), entState, borrowerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response
	loans := make([]model.Loan, len(loanEntities))
	for i, l := range loanEntities {
		loans[i] = model.ToLoanModel(l)
	}

	if err = json.NewEncoder(w).Encode(loans); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

// Serve serves http server for specified address and path
func (h *LoanHTTPTransport) Serve(address, path string) error {
	handler := h.Handler(path)

	server := &http.Server{Addr: address,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
