package service

import (
	"context"
	"log"

	"github.com/albinzx/loan/entity"
	"github.com/albinzx/loan/repository"
)

type LoanEngine interface {
	// Create creates a new loan with the given details.
	Create(context.Context, *entity.Loan) (*entity.Loan, error)
	// Get retrieves a loan by its ID.
	Get(context.Context, int64) (*entity.Loan, error)
	// Approve approves a loan with the given ID and approval details.
	Approve(context.Context, int64, entity.Approval) (*entity.Loan, error)
	// Invest invests in a loan with the given ID and investment details.
	Invest(context.Context, int64, entity.Investment) (*entity.Loan, error)
	// Disburse disburses a loan with the given ID and disbursement details.
	Disburse(context.Context, int64, entity.Approval) (*entity.Loan, error)
	// GetByState retrieves loans based on state.
	GetByState(context.Context, string) ([]entity.Loan, error)
	// GetByBorrower retrieves loans by their borrower ID.
	GetByBorrower(context.Context, int64) ([]entity.Loan, error)
}

type LoanService struct {
	repo repository.LoanRepository
}

func New(repo repository.LoanRepository) *LoanService {
	return &LoanService{repo: repo}
}

func (l *LoanService) Create(ctx context.Context, loan *entity.Loan) (*entity.Loan, error) {

	// Calculate the ROI based on the loan amount and rate
	loan.ROI = (float64(loan.Rate) / 100) * float64(loan.Amount)

	// Set the loan state to Proposed
	loan.State = &entity.Proposed{}

	rloan, err := l.repo.InsertLoan(ctx, *loan)

	return &rloan, err
}

func (l *LoanService) Get(ctx context.Context, id int64) (*entity.Loan, error) {
	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("error while get loan, %v", err)
		return nil, err
	}

	return &loan, nil
}

func (l *LoanService) Approve(ctx context.Context, id int64, approval entity.Approval) (*entity.Loan, error) {
	if approval.Empty() {
		return nil, nil
	}

	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, err
	}

	prevState := loan.State

	if loan.State.Approve(&loan, approval) {
		approval.LoanID = id
		if err := l.repo.UpdateLoanApproval(ctx, loan, approval, prevState); err != nil {
			return nil, err
		}

		return &loan, nil
	}

	return nil, nil
}

func (l *LoanService) Invest(ctx context.Context, id int64, investment entity.Investment) (*entity.Loan, error) {

	if investment.Empty() {
		return nil, nil
	}

	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, err
	}

	prevState := loan.State

	if loan.State.Invest(&loan, investment) {
		investment.LoanID = id
		if prevState.String() == loan.State.String() {
			if err := l.repo.InsertLoanInvestment(ctx, investment); err != nil {
				return nil, err
			}

		} else {

			if err := l.repo.UpdateLoanInvestmentAndState(ctx, loan, investment, prevState); err != nil {
				return nil, err
			}
		}
		return &loan, nil
	}

	return nil, nil
}

func (l *LoanService) Disburse(ctx context.Context, id int64, disbursement entity.Approval) (*entity.Loan, error) {
	if disbursement.Empty() {
		return nil, nil
	}

	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, err
	}

	prevState := loan.State

	if loan.State.Disburse(&loan, disbursement) {
		disbursement.LoanID = id
		if err := l.repo.UpdateLoanApproval(ctx, loan, disbursement, prevState); err != nil {
			return nil, err
		}

		return &loan, nil
	}

	return nil, nil
}

func (l *LoanService) GetByState(ctx context.Context, state string) ([]entity.Loan, error) {
	return l.repo.GetLoansByState(ctx, state)
}

func (l *LoanService) GetByBorrower(ctx context.Context, borrowerID int64) ([]entity.Loan, error) {
	return l.repo.GetLoansByBorrower(ctx, borrowerID)
}

func (l *LoanService) GetLoansByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error) {
	return l.repo.GetLoansByInvestor(ctx, investorID)
}
