package model

import (
	"time"

	"github.com/albinzx/loan/entity"
)

// Loan is a loan model in request or response.
type Loan struct {
	ID                 int64        `json:"id"`
	Amount             int64        `json:"amount"`
	Rate               float64      `json:"rate"`
	ROI                float64      `json:"roi"`
	BorrowerID         int64        `json:"borrower_id"`
	AgreementLetterURL string       `json:"agreement_letter_url"`
	Investments        []Investment `json:"investments"`
	Approvals          []Approval   `json:"approvals"`
	State              string       `json:"state"`
}

// Investment is an investment model in request or response.
type Investment struct {
	LoanID     int64 `json:"loan_id"`
	InvestorID int64 `json:"investor_id"`
	Amount     int64 `json:"amount"`
}

// Approval is an approval model in request or response.
type Approval struct {
	LoanID      int64  `json:"loan_id"`
	EmployeeID  int64  `json:"employee_id"`
	Date        string `json:"date"`
	DocumentURL string `json:"document_url"`
}

// Employee is an employee model in request or response.
type Employee struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Borrower is an borrower model in request or response.
type Borrower struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Limit int64  `json:"limit"`
}

// Investor is an investor model in request or response.
type Investor struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToEntity converts loan model to entity
func (l Loan) ToEntity() *entity.Loan {
	return &entity.Loan{
		ID:                 l.ID,
		Amount:             l.Amount,
		Rate:               l.Rate,
		ROI:                l.ROI,
		BorrowerID:         l.BorrowerID,
		AgreementLetterURL: l.AgreementLetterURL,
	}
}

// Empty checks if approval mandatory fields is empty
func (a Approval) Empty() bool {
	return a.EmployeeID == 0 || a.DocumentURL == ""
}

// ToEntity converts approval model to entity
func (a Approval) ToEntity() *entity.Approval {
	return &entity.Approval{
		LoanID:      a.LoanID,
		EmployeeID:  a.EmployeeID,
		Date:        time.Now(),
		DocumentURL: a.DocumentURL,
	}
}

// Empty checks if investment mandatory fields is empty
func (i Investment) Empty() bool {
	return i.InvestorID == 0 || i.Amount == 0
}

// ToEntity converts investment model to entity
func (i Investment) ToEntity() *entity.Investment {
	return &entity.Investment{
		LoanID:     i.LoanID,
		InvestorID: i.InvestorID,
		Amount:     i.Amount,
	}
}

// ToLoanModel converts loan entity to model
func ToLoanModel(l entity.Loan) Loan {
	var investments []Investment
	var approvals []Approval

	if l.Investments != nil {
		n := len(l.Investments)
		investments = make([]Investment, n)
		for i := range n {
			investments[i] = ToInvestmentModel(l.Investments[i])
		}
	}

	if l.Approvals != nil {
		n := len(l.Approvals)
		approvals = make([]Approval, n)
		for i := range n {
			approvals[i] = ToApprovalModel(l.Approvals[i])
		}
	}

	return Loan{
		ID:                 l.ID,
		Amount:             l.Amount,
		Rate:               l.Rate,
		ROI:                l.ROI,
		BorrowerID:         l.BorrowerID,
		AgreementLetterURL: l.AgreementLetterURL,
		Investments:        investments,
		Approvals:          approvals,
		State:              l.State.String(),
	}
}

// ToApprovalModel converts approval entity to model
func ToApprovalModel(a entity.Approval) Approval {
	return Approval{
		LoanID:      a.LoanID,
		EmployeeID:  a.EmployeeID,
		Date:        a.Date.String(),
		DocumentURL: a.DocumentURL,
	}
}

// ToInvestmentModel converts investment entity to model
func ToInvestmentModel(i entity.Investment) Investment {
	return Investment{
		LoanID:     i.LoanID,
		InvestorID: i.InvestorID,
		Amount:     i.Amount,
	}
}
