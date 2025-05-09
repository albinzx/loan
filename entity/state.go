package entity

import "fmt"

// State interface represents the state of a loan in the system.
// It defines methods for approving, investing, and disbursing loans.
// Each state implements these methods to define the behavior of the loan in that state.
// The State interface is used to manage the loan lifecycle and transition between different states.
// The states include Proposed, Approved, Invested, and Disbursed.
type State interface {
	fmt.Stringer
	Approve(*Loan, Approval) bool
	Invest(*Loan, Investment) bool
	Disburse(*Loan, Approval) bool
}

// Proposed represents the proposed state of a loan.
// It is initial state of a loan, in this state, the loan is proposed but not yet approved.
// The Proposed state allows for the approval of the loan by adding an approval action.
// Other actions like investing or disbursing are not allowed in this state.
type Proposed struct {
}

// String returns the string representation of the Proposed state.
func (p *Proposed) String() string {
	return "PROPOSED"
}

// Approve allows for the approval of a loan in the Proposed state.
// If an approval is provided, it adds the approval to the loan and transitions the loan to the Approved state.
// If no approval is provided, it returns false indicating that the approval was not successful.
func (p *Proposed) Approve(loan *Loan, approval Approval) bool {

	loan.AddApproval(approval)
	loan.State = &Approved{}

	return true
}

// Invest is not allowed in a loan in the Proposed state.
func (p *Proposed) Invest(l *Loan, i Investment) bool {
	return false
}

// Disburse is not allowed in a loan in the Proposed state.
func (p *Proposed) Disburse(l *Loan, a Approval) bool {
	return false
}

// Approved represents the approved state of a loan.
// In this state, the loan has been approved but not yet invested.
// The Approved state allows for the investment of the loan by adding an investment action.
// Other actions like approving or disbursing are not allowed in this state.
type Approved struct {
}

// String returns the string representation of the Approved state.
func (a *Approved) String() string {
	return "APPROVED"
}

// Approve is not allowed in a loan in the Approved state.
func (a *Approved) Approve(loan *Loan, approval Approval) bool {
	return false
}

// Invest allows for the investment of a loan in the Approved state.
// If an investment is provided, it adds the investment to the loan
// If total investments amount is equal with principal amount change state of the loan to the Invested state.
func (a *Approved) Invest(loan *Loan, investment Investment) bool {
	if added, completed := loan.AddInvestment(investment); added {
		if completed {
			loan.State = &Invested{}
		}
		return true
	}

	return false
}

// Disburse is not allowed in a loan in the Approved state.
func (a *Approved) Disburse(loan *Loan, approval Approval) bool {
	return false
}

// Invested represents the invested state of a loan.
// In this state, the loan has been invested but not yet disbursed.
// The Invested state allows for the disbursement of the loan by adding an disbursement action.
// Other actions like approving or investing are not allowed in this state.
type Invested struct {
}

// String returns the string representation of the Invested state.
func (i *Invested) String() string {
	return "INVESTED"
}

// Approve is not allowed in a loan in the Invested state.
func (i *Invested) Approve(loan *Loan, approval Approval) bool {
	return false
}

// Invest is not allowed in a loan in the Invested state.
func (i *Invested) Invest(loan *Loan, investment Investment) bool {
	return false
}

// Disburse allows for the disbursement of a loan in the Invested state.
// If an approval is provided, it adds the approval to the loan and transitions the loan to the Disbursed state.
// If no approval is provided, it returns false indicating that the disbursement was not successful.
// The disbursement action is used to finalize the loan process and make the funds available to the borrower.
func (i *Invested) Disburse(loan *Loan, disbursement Approval) bool {
	loan.AddApproval(disbursement)
	loan.State = &Disbursed{}
	return true
}

// Disbursed represents the disbursed state of a loan.
// In this state, the loan has been disbursed and is considered complete.
// The Disbursed state does not allow for any further actions like approving, investing, or disbursing.
// The Disbursed state indicates that the loan process is complete and the funds have been made available to the borrower.
// The Disbursed state is the final state in the loan lifecycle.
type Disbursed struct {
}

// String returns the string representation of the Disbursed state.
func (d *Disbursed) String() string {
	return "DISBURSED"
}

// Approve is not allowed in a loan in the Disbursed state.
func (d *Disbursed) Approve(loan *Loan, approval Approval) bool {
	return false
}

// Invest is not allowed in a loan in the Disbursed state.
func (d *Disbursed) Invest(loan *Loan, investment Investment) bool {
	return false
}

// Disburse is not allowed in a loan in the Disbursed state.
func (d *Disbursed) Disburse(loan *Loan, approval Approval) bool {
	return false
}

func StateOf(code string) State {
	switch code {
	case "PROPOSED":
		return &Proposed{}
	case "APPROVED":
		return &Approved{}
	case "INVESTED":
		return &Invested{}
	case "DISBURSED":
		return &Disbursed{}
	}

	return &Proposed{}
}
