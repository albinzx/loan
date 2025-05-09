package service

import (
	"context"
	"fmt"
	"log"

	"github.com/albinzx/loan/entity"
	"github.com/albinzx/loan/pkg/mailer"
	"github.com/albinzx/loan/repository"
)

const (
	EMAIL_SUBJECT = "Loan %d is invested"
	EMAIL_MSG     = "Dear %s,\nPlease find the link of aggrement letter:%s\nThanks"
)

type LoanEngine interface {
	// Create creates a new loan with the given details.
	Create(context.Context, *entity.Loan) (*entity.Loan, error)
	// Get retrieves a loan by its ID.
	Get(context.Context, int64) (*entity.Loan, error)
	// Approve approves a loan with the given ID and approval details.
	Approve(context.Context, int64, *entity.Approval) (*entity.Loan, error)
	// Invest invests in a loan with the given ID and investment details.
	Invest(context.Context, int64, *entity.Investment) (*entity.Loan, error)
	// Disburse disburses a loan with the given ID and disbursement details.
	Disburse(context.Context, int64, *entity.Approval) (*entity.Loan, error)
	// GetLoansByInvestor retrieves loans by investor ID
	GetByInvestor(context.Context, int64) ([]entity.Loan, error)
	// GetByStateOrBorrower retrieves loans based on state or borrower id.
	GetByStateOrBorrower(context.Context, entity.State, int64) ([]entity.Loan, error)
}

type LoanService struct {
	repo repository.LoanRepository
	mail mailer.Mailer
}

func New(repo repository.LoanRepository, mail mailer.Mailer) *LoanService {
	return &LoanService{repo: repo, mail: mail}
}

func (l *LoanService) Create(ctx context.Context, loan *entity.Loan) (*entity.Loan, error) {

	// Calculate the ROI based on the loan amount and rate
	loan.ROI = (float64(loan.Rate) / 100) * float64(loan.Amount)

	// Set the loan state to Proposed
	loan.State = &entity.Proposed{}

	// store loan entity
	id, err := l.repo.InsertLoan(ctx, loan)
	if err != nil {
		log.Printf("error while creating loan, %v", err)
		return nil, err
	}

	loan.ID = id

	return loan, nil
}

func (l *LoanService) Get(ctx context.Context, id int64) (*entity.Loan, error) {
	// get loan entity
	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("error while getting loan, %v", err)
		return nil, err
	}

	return loan, nil
}

func (l *LoanService) Approve(ctx context.Context, id int64, approval *entity.Approval) (*entity.Loan, error) {
	// get loan entity
	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("error while getting loan, %v", err)
		return nil, err
	}

	approval.Action = entity.APPROVE
	prevState := loan.State

	if loan.State.Approve(loan, *approval) {
		// loan can be approved, update state and insert approval data
		approval.LoanID = id
		if err := l.repo.UpdateLoanApproval(ctx, loan, approval, prevState); err != nil {
			log.Printf("error while approving loan, %v", err)
			return nil, err
		}

		return loan, nil
	}

	return nil, nil
}

func (l *LoanService) Invest(ctx context.Context, id int64, investment *entity.Investment) (*entity.Loan, error) {
	// get loan entity
	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("error while getting loan, %v", err)
		return nil, err
	}

	prevState := loan.State

	if loan.State.Invest(loan, *investment) {
		investment.LoanID = id

		// loan state still approved, total invested hasn't reach the principal amount, just store investment data
		if prevState.String() == loan.State.String() {
			if err := l.repo.InsertLoanInvestment(ctx, investment); err != nil {
				log.Printf("error while adding loan investment, %v", err)
				return nil, err
			}

		} else {
			// loan state changed to invested, store investment data, change state and send email to investor
			if err := l.repo.UpdateLoanInvestmentAndState(ctx, loan, investment, prevState); err != nil {
				log.Printf("error while investing loan, %v", err)
				return nil, err
			}

			// send email asynchronosly
			go func() {
				// get investor data
				investors, err := l.repo.GetInvestorByLoanID(ctx, id)
				if err != nil {
					log.Printf("error while getting investors, %v", err)
					return
				}

				// send email to each investor
				subject := fmt.Sprintf(EMAIL_SUBJECT, id)
				for _, investor := range investors {
					l.mail.Send(ctx, mailer.Email{
						Recipients: []string{investor.Email},
						Subject:    subject,
						Message:    fmt.Sprintf(EMAIL_MSG, investor.Name, loan.AgreementLetterURL),
					})
				}

			}()
		}
		return loan, nil
	}

	return nil, nil
}

func (l *LoanService) Disburse(ctx context.Context, id int64, disbursement *entity.Approval) (*entity.Loan, error) {
	// get loan entity
	loan, err := l.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("error while getting loan, %v", err)
		return nil, err
	}

	disbursement.Action = entity.DISBURSE
	prevState := loan.State

	if loan.State.Disburse(loan, *disbursement) {
		// loan state change to disbursed, insert disbursement data and change state
		disbursement.LoanID = id
		if err := l.repo.UpdateLoanApproval(ctx, loan, disbursement, prevState); err != nil {
			log.Printf("error while disbursing loan, %v", err)
			return nil, err
		}

		return loan, nil
	}

	return nil, nil
}

func (l *LoanService) GetByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error) {
	return l.repo.GetLoansByInvestor(ctx, investorID)
}

func (l *LoanService) GetByStateOrBorrower(ctx context.Context, state entity.State, borrowerID int64) ([]entity.Loan, error) {
	return l.repo.GetLoansByStateOrBorrower(ctx, state, borrowerID)
}
