package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/albinzx/loan/entity"
	"github.com/albinzx/loan/pkg/mailer"
	"github.com/albinzx/loan/repository"
)

type MockMailer struct{}

func (mm *MockMailer) Send(context.Context, mailer.Email) error {
	return nil
}

var errRepo = errors.New("error repository")

// mock repository
//
// use loan id
// 1 to get proposed loan
// 2 to get approved loan
// 3 to get invested loan
// 4 to get disbursed loan
// 99 for get loan error
//
// use employee id or investor id
// 88 for other query error
type MockRepository struct{}

func (mr *MockRepository) InsertLoan(ctx context.Context, loan *entity.Loan) (int64, error) {
	if loan.ID == 99 {
		return 0, errRepo
	}
	return 1, nil
}

func (mr *MockRepository) GetLoan(ctx context.Context, id int64) (*entity.Loan, error) {
	if id == 99 {
		return nil, errRepo
	}
	if id == 1 {
		return &entity.Loan{ID: 1, Amount: 1000, State: &entity.Proposed{}}, nil
	}
	if id == 2 {
		return &entity.Loan{ID: 2, Amount: 1000, State: &entity.Approved{}}, nil
	}
	if id == 3 {
		return &entity.Loan{ID: 3, Amount: 1000, State: &entity.Invested{}}, nil
	}
	if id == 4 {
		return &entity.Loan{ID: 4, Amount: 1000, State: &entity.Disbursed{}}, nil
	}

	return &entity.Loan{ID: id, Amount: 1000, State: &entity.Proposed{}}, nil
}

func (mr *MockRepository) UpdateState(ctx context.Context, loan *entity.Loan, previousState entity.State) error {
	if loan.ID == 99 {
		return errRepo
	}
	return nil
}

func (mr *MockRepository) InsertLoanInvestment(ctx context.Context, investment *entity.Investment) error {
	if investment.InvestorID == 88 {
		return errRepo
	}
	return nil
}

func (mr *MockRepository) UpdateLoanInvestmentAndState(ctx context.Context, loan *entity.Loan,
	investment *entity.Investment, previousState entity.State) error {
	if investment.InvestorID == 88 {
		return errRepo
	}
	return nil
}
func (mr *MockRepository) UpdateLoanApproval(ctx context.Context, loan *entity.Loan,
	approval *entity.Approval, previousState entity.State) error {
	if approval.EmployeeID == 88 {
		return errRepo
	}
	return nil
}

func (mr *MockRepository) GetLoansByState(ctx context.Context, state entity.State) ([]entity.Loan, error) {
	return nil, nil
}

func (mr *MockRepository) GetLoansByBorrower(ctx context.Context, borrowerID int64) ([]entity.Loan, error) {
	return nil, nil
}

func (mr *MockRepository) GetLoansByInvestor(ctx context.Context, investorID int64) ([]entity.Loan, error) {
	return nil, nil
}

func (mr *MockRepository) GetInvestorByLoanID(ctx context.Context, loanID int64) ([]entity.Investor, error) {
	return nil, nil
}

func TestLoanService_Create(t *testing.T) {
	type fields struct {
		repo repository.LoanRepository
		mail mailer.Mailer
	}
	type args struct {
		ctx  context.Context
		loan *entity.Loan
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Loan
		wantErr bool
	}{
		{
			name: "Test create success",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:  context.Background(),
				loan: &entity.Loan{ID: 1, Amount: 1000, Rate: 10},
			},
			want:    &entity.Loan{ID: 1, Amount: 1000, Rate: 10, ROI: 100, State: &entity.Proposed{}},
			wantErr: false,
		},
		{
			name: "Test create failed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:  context.Background(),
				loan: &entity.Loan{ID: 99},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoanService{
				repo: tt.fields.repo,
				mail: tt.fields.mail,
			}
			got, err := l.Create(tt.args.ctx, tt.args.loan)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoanService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoanService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoanService_Get(t *testing.T) {
	type fields struct {
		repo repository.LoanRepository
		mail mailer.Mailer
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Loan
		wantErr bool
	}{
		{
			name: "Test get success",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &entity.Loan{ID: 1, Amount: 1000, State: &entity.Proposed{}},
			wantErr: false,
		},
		{
			name: "Test get failed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx: context.Background(),
				id:  99,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoanService{
				repo: tt.fields.repo,
				mail: tt.fields.mail,
			}
			got, err := l.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoanService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoanService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoanService_Approve(t *testing.T) {
	type fields struct {
		repo repository.LoanRepository
		mail mailer.Mailer
	}
	type args struct {
		ctx      context.Context
		id       int64
		approval *entity.Approval
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Loan
		wantErr bool
	}{
		{
			name: "Test approve success",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:      context.Background(),
				id:       1,
				approval: &entity.Approval{},
			},
			want: &entity.Loan{
				ID:        1,
				Amount:    1000,
				State:     &entity.Approved{},
				Approvals: []entity.Approval{{Action: entity.APPROVE}},
			},
			wantErr: false,
		},
		{
			name: "Test approve do nothing due to state already approved",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:      context.Background(),
				id:       2,
				approval: &entity.Approval{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test approve do nothing due to state already invested",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:      context.Background(),
				id:       3,
				approval: &entity.Approval{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test approve do nothing due to state already disbursed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:      context.Background(),
				id:       4,
				approval: &entity.Approval{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test approve failed due to get loan error",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:      context.Background(),
				id:       99,
				approval: &entity.Approval{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test approve failed due to update loan approval error",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:      context.Background(),
				id:       1,
				approval: &entity.Approval{EmployeeID: 88},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoanService{
				repo: tt.fields.repo,
				mail: tt.fields.mail,
			}
			got, err := l.Approve(tt.args.ctx, tt.args.id, tt.args.approval)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoanService.Approve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoanService.Approve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoanService_Invest(t *testing.T) {
	type fields struct {
		repo repository.LoanRepository
		mail mailer.Mailer
	}
	type args struct {
		ctx        context.Context
		id         int64
		investment *entity.Investment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Loan
		wantErr bool
	}{
		{
			name: "Test invest success, total investment still under principal amount, state still approved",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         2,
				investment: &entity.Investment{Amount: 500},
			},
			want: &entity.Loan{
				ID:          2,
				Amount:      1000,
				State:       &entity.Approved{},
				Investments: []entity.Investment{{Amount: 500}},
			},
			wantErr: false,
		},
		{
			name: "Test invest success, state change to invested",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         2,
				investment: &entity.Investment{Amount: 1000},
			},
			want: &entity.Loan{
				ID:          2,
				Amount:      1000,
				State:       &entity.Invested{},
				Investments: []entity.Investment{{Amount: 1000}},
			},
			wantErr: false,
		},
		{
			name: "Test invest do nothing due to state still proposed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         1,
				investment: &entity.Investment{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test invest do nothing due to state already invested",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         3,
				investment: &entity.Investment{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test invest do nothing due to state already disbursed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         4,
				investment: &entity.Investment{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test invest failed due to get loan error",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         99,
				investment: &entity.Investment{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test invest failed due to update loan investment error",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:        context.Background(),
				id:         2,
				investment: &entity.Investment{InvestorID: 88},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoanService{
				repo: tt.fields.repo,
				mail: tt.fields.mail,
			}
			got, err := l.Invest(tt.args.ctx, tt.args.id, tt.args.investment)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoanService.Invest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoanService.Invest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoanService_Disburse(t *testing.T) {
	type fields struct {
		repo repository.LoanRepository
		mail mailer.Mailer
	}
	type args struct {
		ctx          context.Context
		id           int64
		disbursement *entity.Approval
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Loan
		wantErr bool
	}{
		{
			name: "Test disburse success",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:          context.Background(),
				id:           3,
				disbursement: &entity.Approval{},
			},
			want: &entity.Loan{
				ID:        3,
				Amount:    1000,
				State:     &entity.Disbursed{},
				Approvals: []entity.Approval{{Action: entity.DISBURSE}},
			},
			wantErr: false,
		},
		{
			name: "Test disburse do nothing due to state still proposed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:          context.Background(),
				id:           1,
				disbursement: &entity.Approval{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test disburse do nothing due to state still approved",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:          context.Background(),
				id:           2,
				disbursement: &entity.Approval{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test disburse do nothing due to state already disbursed",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:          context.Background(),
				id:           4,
				disbursement: &entity.Approval{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test disburse failed due to get loan error",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:          context.Background(),
				id:           99,
				disbursement: &entity.Approval{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test disburse failed due to update loan approval error",
			fields: fields{
				repo: &MockRepository{},
				mail: &MockMailer{},
			},
			args: args{
				ctx:          context.Background(),
				id:           3,
				disbursement: &entity.Approval{EmployeeID: 88},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoanService{
				repo: tt.fields.repo,
				mail: tt.fields.mail,
			}
			got, err := l.Disburse(tt.args.ctx, tt.args.id, tt.args.disbursement)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoanService.Disburse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoanService.Disburse() = %v, want %v", got, tt.want)
			}
		})
	}
}
