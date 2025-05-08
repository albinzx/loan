package entity

// Loan represents a loan entity in the system.
// It contains information about the loan amount, rate, borrower, and associated investments and approvals.
// The Loan struct is used to manage the loan lifecycle, including approval, investment, and disbursement processes.
// It also includes methods for calculating the return on investment (ROI) and managing the loan state.
type Loan struct {
	ID                 int64        `json:"id"`
	Amount             int64        `json:"amount"`
	Rate               float64      `json:"rate"`
	ROI                float64      `json:"roi"`
	BorrowerID         int64        `json:"borrower_id"`
	AgreementLetterURL string       `json:"agreement_letter_url"`
	Investments        []Investment `json:"investments"`
	Approvals          []Approval   `json:"approvals"`
	State              State        `json:"state"`
}

// Investment represents an investment made by an investor in a loan.
// It contains information about the loan ID, investor ID, investment amount, and return on investment (ROI).
// The Investment struct is used to track the investments made in a loan and calculate the ROI for each investment.
type Investment struct {
	LoanID     int64 `json:"loan_id"`
	InvestorID int64 `json:"investor_id"`
	Amount     int64 `json:"amount"`
}

// Approval represents an approval action taken by an employee on a loan.
// It contains information about the loan ID, employee ID, date of action, action taken, and a URL to the prove or approval document.
// The Approval is used to track the approval history of a loan and the actions taken by employees in the system.
// It helps in maintaining a record of the approval process and the documents associated with it.
type Approval struct {
	LoanID      int64  `json:"loan_id"`
	EmployeeID  int64  `json:"employee_id"`
	Date        string `json:"date"`
	Action      string `json:"action"`
	DocumentURL string `json:"document_url"`
}

// Employee represents an employee in the system.
// It contains information about the employee ID, name, and role.
type Employee struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Borrower represents a borrower in the system.
// It contains information about the borrower ID, name, and credit limit.
type Borrower struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Limit int64  `json:"limit"`
}

// Investor represents an investor in the system.
// It contains information about the investor ID and name.
type Investor struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (a Approval) Empty() bool {
	return a.EmployeeID == 0 || a.Date == "" && a.Action == "" && a.DocumentURL == ""
}

func (i Investment) Empty() bool {
	return i.InvestorID == 0 || i.Amount == 0
}

// SumInvestment calculates the total amount of investments made in the loan.
// It iterates through the Investments slice and sums up the Amount of each investment.
func (l *Loan) SumInvestment() int64 {
	var sum int64
	for _, investment := range l.Investments {
		sum += investment.Amount
	}
	return sum
}

// AddInvestment adds an investment to the loan.
// It checks if the total amount of investments does not exceed the loan amount.
// If the investment is valid, it appends the investment to the Investments slice.
// It returns two boolean values:
// 1. A boolean indicating whether the investment was successfully added.
// 2. A boolean indicating whether the total amount of investments is equal to the loan amount.
func (l *Loan) AddInvestment(i Investment) (bool, bool) {
	total := l.SumInvestment() + i.Amount
	if total > l.Amount {
		return false, false
	}

	l.Investments = append(l.Investments, i)
	if total == l.Amount {
		return true, true
	}

	return true, false
}

// AddApproval adds an approval to the loan.
// It appends the approval to the Approvals slice.
func (l *Loan) AddApproval(a Approval) {
	l.Approvals = append(l.Approvals, a)
}
