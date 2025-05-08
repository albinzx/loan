package repository

import (
	"context"
	"database/sql"

	"github.com/albinzx/loan/entity"
)

type LoanRepository interface {
	// InsertLoan inserts a new loan into the repository.
	InsertLoan(ctx context.Context, loan entity.Loan) (entity.Loan, error)
	// GetLoan retrieves a loan by its ID from the repository.
	GetLoan(ctx context.Context, id int64) (entity.Loan, error)
	// UpdateState updates an existing loan state in the repository.
	UpdateState(ctx context.Context, loan entity.Loan, previousState entity.State) error
	// UpdateLoanInvestment inserts a new investment into the repository.
	InsertLoanInvestment(ctx context.Context, investment entity.Investment) error
	// UpdateLoanInvestment inserts a new investment into the repository.
	UpdateLoanInvestmentAndState(ctx context.Context, loan entity.Loan, investment entity.Investment, previousState entity.State) error
	// UpdateLoanApproval inserts a new approval into the repository.
	UpdateLoanApproval(ctx context.Context, loan entity.Loan, approval entity.Approval, previousState entity.State) error
	// GetLoansByState retrieves loans by their state from the repository.
	GetLoansByState(ctx context.Context, state string) ([]entity.Loan, error)
	// GetLoansByBorrower retrieves loans by their borrower ID from the repository.
	GetLoansByBorrower(ctx context.Context, borrowerID int64) ([]entity.Loan, error)
	// GetLoansByInvestor retrieves loans by their investor ID from the repository.
	GetLoansByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error)
}

const (
	insertLoan               = "INSERT INTO loan(amount, rate, roi, borrower_id, agreement_letter_url, state) VALUES (?,?,?,?,?,?)"
	insertInvestment         = "INSERT INTO loan_investment(loan_id, investor_id, amount) VALUES (?,?,?)"
	insertApproval           = "INSERT INTO loan_approval(loan_id, employee_id, approval_date, action, document_url) VALUES (?,?,?,?,?)"
	updateLoanState          = "UPDATE LOAN SET state = ? WHERE id = ? AND state = ?"
	selectLoanByID           = "SELECT amount, rate, roi, borrower_id, agreement_letter_url, state FROM loan WHERE id = ?"
	selectInvestmentByLoanID = "SELECT investor_id, amount FROM loan_investment WHERE loan_id = ?"
)

type MysqlLoanRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *MysqlLoanRepository {
	return &MysqlLoanRepository{db: db}
}

// InsertLoan inserts a new loan into the repository.
func (r *MysqlLoanRepository) InsertLoan(ctx context.Context, loan entity.Loan) (entity.Loan, error) {
	res, err := r.db.ExecContext(ctx, insertLoan, loan.Amount, loan.Rate, loan.ROI,
		loan.BorrowerID, loan.AgreementLetterURL, loan.State.String())
	if err != nil {
		return loan, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return loan, err
	}

	loan.ID = id

	return loan, nil
}

// GetLoan retrieves a loan by its ID from the repository.
func (r *MysqlLoanRepository) GetLoan(ctx context.Context, id int64) (entity.Loan, error) {
	loan := entity.Loan{ID: id}
	var state string
	err := r.db.QueryRowContext(ctx, selectLoanByID, id).Scan(&loan.Amount, &loan.Rate, &loan.ROI,
		&loan.BorrowerID, &loan.AgreementLetterURL, &state)

	if err != nil {
		return loan, err
	}

	loan.State = entity.StateOf(state)

	rows, err := r.db.QueryContext(ctx, selectInvestmentByLoanID, id)
	if err != nil {
		return loan, err
	}
	defer rows.Close()

	invs := make([]entity.Investment, 0)
	for rows.Next() {
		inv := entity.Investment{LoanID: id}
		rows.Scan(&id, &inv.InvestorID, &inv.Amount)
		invs = append(invs, inv)
	}
	if err != nil {
		return loan, err
	}

	loan.Investments = invs

	return loan, nil
}

// UpdateState updates an existing loan state in the repository.
func (r *MysqlLoanRepository) UpdateState(ctx context.Context, loan entity.Loan, previousState entity.State) error {
	_, err := r.db.ExecContext(ctx, updateLoanState, loan.State.String(), loan.ID, previousState.String())

	return err
}

// InsertLoanInvestment inserts a new investment into the repository.
func (r *MysqlLoanRepository) InsertLoanInvestment(ctx context.Context, investment entity.Investment) error {
	_, err := r.db.ExecContext(ctx, insertInvestment, investment.LoanID, investment.InvestorID, investment.Amount)

	return err
}

// UpdateLoanInvestmentAndState inserts a new investment and update loan state into the repository.
func (r *MysqlLoanRepository) UpdateLoanInvestmentAndState(ctx context.Context, loan entity.Loan,
	investment entity.Investment, previousState entity.State) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, insertInvestment, investment.LoanID, investment.InvestorID, investment.Amount)
	tx.ExecContext(ctx, updateLoanState, loan.State.String(), loan.ID, previousState.String())

	return tx.Commit()
}

// UpdateLoanApproval inserts a new approval into the repository.
func (r *MysqlLoanRepository) UpdateLoanApproval(ctx context.Context, loan entity.Loan,
	approval entity.Approval, previousState entity.State) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, insertApproval, approval.LoanID, approval.EmployeeID, approval.Date,
		approval.Action, approval.DocumentURL)
	tx.ExecContext(ctx, updateLoanState, loan.State.String(), loan.ID, previousState.String())

	return tx.Commit()
}

// GetLoansByState retrieves loans by their state from the repository.
func (r *MysqlLoanRepository) GetLoansByState(ctx context.Context, state string) ([]entity.Loan, error) {
	return nil, nil
}

// GetLoansByBorrower retrieves loans by their borrower ID from the repository.
func (r *MysqlLoanRepository) GetLoansByBorrower(ctx context.Context, borrowerID int64) ([]entity.Loan, error) {
	return nil, nil
}

// GetLoansByInvestor retrieves loans by their investor ID from the repository.
func (r *MysqlLoanRepository) GetLoansByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error) {
	return nil, nil
}
