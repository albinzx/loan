package entity

import (
	"reflect"
	"testing"
)

func TestProposed_String(t *testing.T) {
	tests := []struct {
		name string
		p    *Proposed
		want string
	}{
		{
			name: "Test proposed to string",
			want: STATE_PROPOSED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proposed{}
			if got := p.String(); got != tt.want {
				t.Errorf("Proposed.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProposed_Approve(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		p    *Proposed
		args args
		want bool
	}{
		{
			name: "Test approve from proposed state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proposed{}
			if got := p.Approve(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Proposed.Approve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProposed_Invest(t *testing.T) {
	type args struct {
		loan       *Loan
		investment Investment
	}
	tests := []struct {
		name string
		p    *Proposed
		args args
		want bool
	}{
		{
			name: "Test invest from proposed state",
			args: args{
				loan:       &Loan{},
				investment: Investment{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proposed{}
			if got := p.Invest(tt.args.loan, tt.args.investment); got != tt.want {
				t.Errorf("Proposed.Invest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProposed_Disburse(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		p    *Proposed
		args args
		want bool
	}{
		{
			name: "Test disburse from proposed state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proposed{}
			if got := p.Disburse(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Proposed.Disburse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproved_String(t *testing.T) {
	tests := []struct {
		name string
		a    *Approved
		want string
	}{
		{
			name: "Test approved to string",
			want: STATE_APPROVED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Approved{}
			if got := a.String(); got != tt.want {
				t.Errorf("Approved.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproved_Approve(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		a    *Approved
		args args
		want bool
	}{
		{
			name: "Test approve from approved state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Approved{}
			if got := a.Approve(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Approved.Approve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproved_Invest(t *testing.T) {
	type args struct {
		loan       *Loan
		investment Investment
	}
	tests := []struct {
		name string
		a    *Approved
		args args
		want bool
	}{
		{
			name: "Test invest from approved state, under total",
			args: args{
				loan:       &Loan{Amount: 1000},
				investment: Investment{Amount: 500},
			},
			want: true,
		},
		{
			name: "Test invest from approved state, equals total",
			args: args{
				loan:       &Loan{Amount: 1000},
				investment: Investment{Amount: 1000},
			},
			want: true,
		},
		{
			name: "Test invest from approved state, exceed total",
			args: args{
				loan:       &Loan{Amount: 1000},
				investment: Investment{Amount: 1100},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Approved{}
			if got := a.Invest(tt.args.loan, tt.args.investment); got != tt.want {
				t.Errorf("Approved.Invest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproved_Disburse(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		a    *Approved
		args args
		want bool
	}{
		{
			name: "Test disburse from approved state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Approved{}
			if got := a.Disburse(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Approved.Disburse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvested_String(t *testing.T) {
	tests := []struct {
		name string
		i    *Invested
		want string
	}{
		{
			name: "Test invested to string",
			want: STATE_INVESTED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Invested{}
			if got := i.String(); got != tt.want {
				t.Errorf("Invested.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvested_Approve(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		i    *Invested
		args args
		want bool
	}{
		{
			name: "Test approve from invested state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Invested{}
			if got := i.Approve(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Invested.Approve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvested_Invest(t *testing.T) {
	type args struct {
		loan       *Loan
		investment Investment
	}
	tests := []struct {
		name string
		i    *Invested
		args args
		want bool
	}{
		{
			name: "Test invest from invested state",
			args: args{
				loan:       &Loan{},
				investment: Investment{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Invested{}
			if got := i.Invest(tt.args.loan, tt.args.investment); got != tt.want {
				t.Errorf("Invested.Invest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvested_Disburse(t *testing.T) {
	type args struct {
		loan         *Loan
		disbursement Approval
	}
	tests := []struct {
		name string
		i    *Invested
		args args
		want bool
	}{
		{
			name: "Test disburse from invested state",
			args: args{
				loan:         &Loan{},
				disbursement: Approval{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Invested{}
			if got := i.Disburse(tt.args.loan, tt.args.disbursement); got != tt.want {
				t.Errorf("Invested.Disburse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisbursed_String(t *testing.T) {
	tests := []struct {
		name string
		d    *Disbursed
		want string
	}{
		{
			name: "Test disbursed to string",
			want: STATE_DISBURSED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disbursed{}
			if got := d.String(); got != tt.want {
				t.Errorf("Disbursed.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisbursed_Approve(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		d    *Disbursed
		args args
		want bool
	}{
		{
			name: "Test approve from disbursed state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disbursed{}
			if got := d.Approve(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Disbursed.Approve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisbursed_Invest(t *testing.T) {
	type args struct {
		loan       *Loan
		investment Investment
	}
	tests := []struct {
		name string
		d    *Disbursed
		args args
		want bool
	}{
		{
			name: "Test invest from disbursed state",
			args: args{
				loan:       &Loan{},
				investment: Investment{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disbursed{}
			if got := d.Invest(tt.args.loan, tt.args.investment); got != tt.want {
				t.Errorf("Disbursed.Invest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisbursed_Disburse(t *testing.T) {
	type args struct {
		loan     *Loan
		approval Approval
	}
	tests := []struct {
		name string
		d    *Disbursed
		args args
		want bool
	}{
		{
			name: "Test disburse from disbursed state",
			args: args{
				loan:     &Loan{},
				approval: Approval{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disbursed{}
			if got := d.Disburse(tt.args.loan, tt.args.approval); got != tt.want {
				t.Errorf("Disbursed.Disburse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateOf(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want State
	}{
		{
			name: "Test convertion to proposed state",
			args: args{
				code: "PROPOSED",
			},
			want: &Proposed{},
		},
		{
			name: "Test convertion to approved state",
			args: args{
				code: "APPROVED",
			},
			want: &Approved{},
		},
		{
			name: "Test convertion to invested state",
			args: args{
				code: "INVESTED",
			},
			want: &Invested{},
		},
		{
			name: "Test convertion to disbursed state",
			args: args{
				code: "DISBURSED",
			},
			want: &Disbursed{},
		},
		{
			name: "Test invalid state",
			args: args{
				code: "REJECTED",
			},
			want: nil,
		},
		{
			name: "Test empty state",
			args: args{
				code: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StateOf(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
