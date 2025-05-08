package transport

import "github.com/albinzx/loan/entity"

type LoanTransport interface {
	// Create creates a new loan with the given details.
	Create(loan *entity.Loan) (*entity.Loan, error)
	// Get retrieves a loan by its ID.
	Get(id int64) (*entity.Loan, error)
	// ApproveLoan approves a loan with the given ID and approval details.
	Approve(int64, entity.Approval) (*entity.Loan, error)
	// InvestLoan invests in a loan with the given ID and investment details.
	Invest(int64, entity.Investment) (*entity.Loan, error)
	// DisburseLoan disburses a loan with the given ID and disbursement details.
	Disburse(int64, entity.Approval) (*entity.Loan, error)
	// GetAll retrieves all loans.
	GetByState(state string) ([]entity.Loan, error)
	// GetByBorrower retrieves loans by their borrower ID.
	GetByBorrower(borrowerID int64) ([]entity.Loan, error)
}
