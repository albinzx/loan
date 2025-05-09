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

func New(service svc.LoanEngine) *LoanHTTPTransport {
	return &LoanHTTPTransport{service: service}
}

func (h *LoanHTTPTransport) Handler(path string) http.Handler {
	router := httprouter.New()

	router.POST(fmt.Sprintf("%s/loans", path), h.Create)
	router.GET(fmt.Sprintf("%s/loans/:id", path), h.Get)
	router.PATCH(fmt.Sprintf("%s/loans/:id/approve", path), h.Approve)
	router.PATCH(fmt.Sprintf("%s/loans/:id/invest", path), h.Invest)
	router.PATCH(fmt.Sprintf("%s/loans/:id/disburse", path), h.Disburse)
	router.GET(fmt.Sprintf("%s/loans-state/:state", path), h.GetByState)
	router.GET(fmt.Sprintf("%s/loans-borrower/:id", path), h.GetByBorrower)
	router.GET(fmt.Sprintf("%s/loans-investor/:id", path), h.GetByInvestor)

	return router
}

func (h *LoanHTTPTransport) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	loan := &model.Loan{}
	if err := json.NewDecoder(r.Body).Decode(loan); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(r.Context(), loan.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loan, err := h.service.Get(r.Context(), loanID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if loan == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*loan)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Approve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	approval := &model.Approval{}
	if err := json.NewDecoder(r.Body).Decode(approval); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if approval.Empty() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Approve(r.Context(), loanID, approval.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Invest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	investment := &model.Investment{}
	if err := json.NewDecoder(r.Body).Decode(investment); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if investment.Empty() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Invest(r.Context(), loanID, investment.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) Disburse(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	loanID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	approval := &model.Approval{}
	if err := json.NewDecoder(r.Body).Decode(approval); err != nil {
		log.Printf("error while decoding request, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if approval.Empty() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Disburse(r.Context(), loanID, approval.ToEntity())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(model.ToLoanModel(*res)); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) GetByState(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	state := params.ByName("state")

	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loanEntities, err := h.service.GetByState(r.Context(), entity.StateOf(state))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loans := make([]model.Loan, len(loanEntities))
	for i, l := range loanEntities {
		loans[i] = model.ToLoanModel(l)
	}

	if err = json.NewEncoder(w).Encode(loans); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) GetByBorrower(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	borrowerID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loanEntities, err := h.service.GetByBorrower(r.Context(), borrowerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loans := make([]model.Loan, len(loanEntities))
	for i, l := range loanEntities {
		loans[i] = model.ToLoanModel(l)
	}

	if err = json.NewEncoder(w).Encode(loans); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

func (h *LoanHTTPTransport) GetByInvestor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	investorID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while parsing parameter id, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loanEntities, err := h.service.GetByInvestor(r.Context(), investorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loans := make([]model.Loan, len(loanEntities))
	for i, l := range loanEntities {
		loans[i] = model.ToLoanModel(l)
	}

	if err = json.NewEncoder(w).Encode(loans); err != nil {
		log.Printf("error while encoding response, %v", err)
	}
}

// Serve serves OAuth2 http server for specified address
func (h *LoanHTTPTransport) Serve(address, path string) error {
	handler := h.Handler(path)

	server := &http.Server{Addr: address,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
