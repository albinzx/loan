package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/albinzx/loan/entity"
)

const (
	insertLoan               = "INSERT INTO loan(amount, rate, roi, borrower_id, agreement_letter_url, state) VALUES (?,?,?,?,?,?)"
	insertInvestment         = "INSERT INTO loan_investment(loan_id, investor_id, amount) VALUES (?,?,?)"
	insertApproval           = "INSERT INTO loan_approval(loan_id, employee_id, approval_date, action, document_url) VALUES (?,?,?,?,?)"
	updateLoanState          = "UPDATE LOAN SET state = ? WHERE id = ? AND state = ?"
	selectLoanByID           = "SELECT amount, rate, roi, borrower_id, agreement_letter_url, state FROM loan WHERE id = ?"
	selectInvestmentByLoanID = "SELECT investor_id, amount FROM loan_investment WHERE loan_id = ?"
	selectApprovalByLoanID   = "SELECT employee_id, approval_date, action, document_url FROM loan_approval WHERE loan_id = ?"
	selectLoanByState        = "SELECT id, amount, rate, roi, borrower_id, agreement_letter_url FROM loan WHERE state = ?"
	selectLoanByBorrowerID   = "SELECT id, amount, rate, roi, agreement_letter_url, state FROM loan WHERE borrower_id = ?"
	selectLoanByInvestorID   = "SELECT id, amount, rate, roi, borrower_id, agreement_letter_url, state FROM loan l WHERE EXISTS " +
		"(SELECT 1 FROM loan_investment i WHERE i.loan_id = l.id and i.investor_id = ?)"
	selectInvestor = "SELECT name, email FROM investor i JOIN loan_investment li ON i.id = li.investor_id WHERE li.loan_id = ?"
)

type LoanMysqlRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *LoanMysqlRepository {
	return &LoanMysqlRepository{db: db}
}

// InsertLoan inserts a new loan into the repository.
func (r *LoanMysqlRepository) InsertLoan(ctx context.Context, loan *entity.Loan) (int64, error) {
	res, err := r.db.ExecContext(ctx, insertLoan, loan.Amount, loan.Rate, loan.ROI,
		loan.BorrowerID, loan.AgreementLetterURL, loan.State.String())
	if err != nil {
		log.Printf("error while inserting loan, %v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error while getting loan ID, %v", err)
		return 0, err
	}

	loan.ID = id

	return id, nil
}

// GetLoan retrieves a loan by its ID from the repository.
func (r *LoanMysqlRepository) GetLoan(ctx context.Context, id int64) (*entity.Loan, error) {
	loan := entity.Loan{ID: id}
	var state string
	if err := r.db.QueryRowContext(ctx, selectLoanByID, id).Scan(&loan.Amount, &loan.Rate, &loan.ROI,
		&loan.BorrowerID, &loan.AgreementLetterURL, &state); err != nil {

		log.Printf("error while getting loan, %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	loan.State = entity.StateOf(state)

	rows, err := r.db.QueryContext(ctx, selectInvestmentByLoanID, id)
	if err != nil {
		log.Printf("error while getting loan investments, %v", err)
		return nil, err
	}
	defer rows.Close()

	invs := make([]entity.Investment, 0)
	for rows.Next() {
		inv := entity.Investment{LoanID: id}
		rows.Scan(&inv.InvestorID, &inv.Amount)
		invs = append(invs, inv)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error while scanning loan investments rows, %v", err)
		return nil, err
	}

	loan.Investments = invs

	rows, err = r.db.QueryContext(ctx, selectApprovalByLoanID, id)
	if err != nil {
		log.Printf("error while getting loan approvals, %v", err)
		return nil, err
	}
	defer rows.Close()

	approvals := make([]entity.Approval, 0)
	for rows.Next() {
		app := entity.Approval{LoanID: id}
		rows.Scan(&app.EmployeeID, &app.Date, &app.Action, &app.DocumentURL)
		approvals = append(approvals, app)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error while scanning loan approvals rows, %v", err)
		return nil, err
	}

	loan.Approvals = approvals

	return &loan, nil
}

// UpdateState updates an existing loan state in the repository.
func (r *LoanMysqlRepository) UpdateState(ctx context.Context, loan *entity.Loan, previousState entity.State) error {
	if _, err := r.db.ExecContext(ctx, updateLoanState, loan.State.String(), loan.ID, previousState.String()); err != nil {
		log.Printf("error while updating loan state, %v", err)
		return err
	}

	return nil
}

// InsertLoanInvestment inserts a new investment into the repository.
func (r *LoanMysqlRepository) InsertLoanInvestment(ctx context.Context, investment *entity.Investment) error {
	if _, err := r.db.ExecContext(ctx, insertInvestment, investment.LoanID, investment.InvestorID, investment.Amount); err != nil {
		log.Printf("error while inserting loan investment, %v", err)
		return err
	}

	return nil
}

// UpdateLoanInvestmentAndState inserts a new investment and update loan state into the repository.
func (r *LoanMysqlRepository) UpdateLoanInvestmentAndState(ctx context.Context, loan *entity.Loan,
	investment *entity.Investment, previousState entity.State) error {

	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("error while creating db transaction, %v", err)
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, insertInvestment, investment.LoanID, investment.InvestorID, investment.Amount); err != nil {
		log.Printf("error while inserting loan investment, %v", err)
	}

	if _, err = tx.ExecContext(ctx, updateLoanState, loan.State.String(), loan.ID, previousState.String()); err != nil {
		log.Printf("error while updating loan state, %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("error while commiting db transaction, %v", err)
	}

	return nil
}

// UpdateLoanApproval inserts a new approval into the repository.
func (r *LoanMysqlRepository) UpdateLoanApproval(ctx context.Context, loan *entity.Loan,
	approval *entity.Approval, previousState entity.State) error {

	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("error while creating db transaction, %v", err)
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, insertApproval, approval.LoanID, approval.EmployeeID, approval.Date,
		approval.Action, approval.DocumentURL); err != nil {
		log.Printf("error while inserting loan approval, %v", err)
	}

	if _, err = tx.ExecContext(ctx, updateLoanState, loan.State.String(), loan.ID, previousState.String()); err != nil {
		log.Printf("error while updating loan state, %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("error while commiting db transaction, %v", err)
	}

	return nil
}

// GetLoansByState retrieves loans by their state from the repository.
func (r *LoanMysqlRepository) GetLoansByState(ctx context.Context, state entity.State) ([]entity.Loan, error) {

	rows, err := r.db.QueryContext(ctx, selectLoanByState, state.String())
	if err != nil {
		log.Printf("error while getting loans by state, %v", err)
		return nil, err
	}
	defer rows.Close()

	loans := make([]entity.Loan, 0)
	for rows.Next() {
		loan := entity.Loan{State: state}
		rows.Scan(&loan.ID, &loan.Amount, &loan.Rate, &loan.ROI, &loan.BorrowerID, &loan.AgreementLetterURL)
		loans = append(loans, loan)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error while scanning loans rows, %v", err)
		return nil, err
	}

	return loans, nil
}

// GetLoansByBorrower retrieves loans by their borrower ID from the repository.
func (r *LoanMysqlRepository) GetLoansByBorrower(ctx context.Context, borrowerID int64) ([]entity.Loan, error) {
	rows, err := r.db.QueryContext(ctx, selectLoanByBorrowerID, borrowerID)
	if err != nil {
		log.Printf("error while getting loans by borrower ID, %v", err)
		return nil, err
	}
	defer rows.Close()

	loans := make([]entity.Loan, 0)
	for rows.Next() {
		loan := entity.Loan{BorrowerID: borrowerID}
		var state string

		rows.Scan(&loan.ID, &loan.Amount, &loan.Rate, &loan.ROI, &loan.AgreementLetterURL, &state)
		loan.State = entity.StateOf(state)

		loans = append(loans, loan)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error while scanning loans rows, %v", err)
		return nil, err
	}

	return loans, nil
}

// GetLoansByInvestor retrieves loans by their investor ID from the repository.
func (r *LoanMysqlRepository) GetLoansByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error) {
	rows, err := r.db.QueryContext(ctx, selectLoanByInvestorID, investorID)
	if err != nil {
		log.Printf("error while getting loans by investor ID, %v", err)
		return nil, err
	}
	defer rows.Close()

	loans := make([]entity.Loan, 0)
	for rows.Next() {
		loan := entity.Loan{}
		var state string

		rows.Scan(&loan.ID, &loan.Amount, &loan.Rate, &loan.ROI, &loan.BorrowerID, &loan.AgreementLetterURL, &state)
		loan.State = entity.StateOf(state)

		loans = append(loans, loan)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error while scanning loans rows, %v", err)
		return nil, err
	}

	return loans, nil
}

func (r *LoanMysqlRepository) GetInvestorByLoanID(ctx context.Context, loanID int64) ([]entity.Investor, error) {
	rows, err := r.db.QueryContext(ctx, selectInvestor, loanID)
	if err != nil {
		log.Printf("error while getting investor emails, %v", err)
		return nil, err
	}
	defer rows.Close()

	investors := make([]entity.Investor, 0)
	for rows.Next() {
		var inv entity.Investor
		rows.Scan(&inv.Name, &inv.Email)
		investors = append(investors, inv)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error while scanning investor rows, %v", err)
		return nil, err
	}

	return investors, nil
}
