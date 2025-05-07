package service

import (
	"github.com/albinzx/loan/entity"
	"github.com/albinzx/loan/repository"
)

type LoanEngine interface {
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

type loanEngine struct {
	repo repository.LoanRepository
}

func New(repo repository.LoanRepository) LoanEngine {
	return &loanEngine{repo: repo}
}

func (l *loanEngine) Create(loan *entity.Loan) (*entity.Loan, error) {

	// Calculate the ROI based on the loan amount and rate
	loan.ROI = (float64(loan.Rate) / 100) * float64(loan.Amount)

	// Set the loan state to Proposed
	loan.State = &entity.Proposed{}

	l.repo.InsertLoan(*loan)

	return loan, nil
}

func (l *loanEngine) Get(id int64) (*entity.Loan, error) {
	loan, err := l.repo.GetLoan(id)
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (l *loanEngine) Approve(id int64, approval entity.Approval) (*entity.Loan, error) {
	if approval.Empty() {
		return nil, nil
	}

	loan, err := l.repo.GetLoan(id)
	if err != nil {
		return nil, err
	}

	if loan.State.Approve(&loan, approval) {
		if err := l.repo.UpdateLoanApproval(loan, approval); err != nil {
			return nil, err
		}

		return &loan, nil
	}

	return nil, nil
}

func (l *loanEngine) Invest(id int64, investment entity.Investment) (*entity.Loan, error) {

	if investment.Empty() {
		return nil, nil
	}

	loan, err := l.repo.GetLoan(id)
	if err != nil {
		return nil, err
	}

	if loan.State.Invest(&loan, investment) {
		if err := l.repo.UpdateLoanInvestment(loan, investment); err != nil {
			return nil, err
		}

		return &loan, nil
	}

	return nil, nil
}

func (l *loanEngine) Disburse(id int64, disbursement entity.Approval) (*entity.Loan, error) {
	if disbursement.Empty() {
		return nil, nil
	}

	loan, err := l.repo.GetLoan(id)
	if err != nil {
		return nil, err
	}

	if loan.State.Disburse(&loan, disbursement) {
		if err := l.repo.UpdateLoanApproval(loan, disbursement); err != nil {
			return nil, err
		}

		return &loan, nil
	}

	return nil, nil
}

func (l *loanEngine) GetByState(state string) ([]entity.Loan, error) {
	return l.repo.GetLoansByState(state)
}

func (l *loanEngine) GetByBorrower(borrowerID int64) ([]entity.Loan, error) {
	return l.repo.GetLoansByBorrower(borrowerID)
}

func (l *loanEngine) GetLoansByInvestor(investorID int64) ([]entity.Loan, error) {
	return l.repo.GetLoansByInvestor(investorID)
}
