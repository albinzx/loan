package repository

import (
	"context"

	"github.com/albinzx/loan/entity"
)

type LoanRepository interface {
	// InsertLoan inserts a new loan into the repository.
	InsertLoan(ctx context.Context, loan *entity.Loan) (int64, error)
	// GetLoan retrieves a loan by its ID from the repository.
	GetLoan(ctx context.Context, id int64) (*entity.Loan, error)
	// UpdateState updates an existing loan state in the repository.
	UpdateState(ctx context.Context, loan *entity.Loan, previousState entity.State) error
	// UpdateLoanInvestment inserts a new investment into the repository.
	InsertLoanInvestment(ctx context.Context, investment *entity.Investment) error
	// UpdateLoanInvestment inserts a new investment into the repository.
	UpdateLoanInvestmentAndState(ctx context.Context, loan *entity.Loan,
		investment *entity.Investment, previousState entity.State) error
	// UpdateLoanApproval inserts a new approval into the repository.
	UpdateLoanApproval(ctx context.Context, loan *entity.Loan,
		approval *entity.Approval, previousState entity.State) error
	// GetLoansByState retrieves loans by their state from the repository.
	GetLoansByState(ctx context.Context, state entity.State) ([]entity.Loan, error)
	// GetLoansByBorrower retrieves loans by their borrower ID from the repository.
	GetLoansByBorrower(ctx context.Context, borrowerID int64) ([]entity.Loan, error)
	// GetLoansByInvestor retrieves loans by their investor ID from the repository.
	GetLoansByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error)
	// GetInvestorEmailByLoanID retrieves investor emails by loan ID theirs invested.
	GetInvestorByLoanID(ctx context.Context, loanID int64) ([]entity.Investor, error)
}
