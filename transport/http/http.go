package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/albinzx/loan/entity"
	svc "github.com/albinzx/loan/service"
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
	router.POST(fmt.Sprintf("%s/loan", path), h.Create)
	router.GET(fmt.Sprintf("%s/loan/:id", path), h.Get)

	return router
}

func (h *LoanHTTPTransport) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	loan := &entity.Loan{}
	if err := json.NewDecoder(r.Body).Decode(loan); err != nil {

		return
	}

	h.service.Create(r.Context(), loan)

	json.NewEncoder(w).Encode(loan)
}

func (h *LoanHTTPTransport) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	loanID, _ := strconv.ParseInt(id, 10, 64)

	loan, _ := h.service.Get(r.Context(), loanID)

	json.NewEncoder(w).Encode(loan)
}

func (h *LoanHTTPTransport) GetByState(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (h *LoanHTTPTransport) GetByBorrower(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
