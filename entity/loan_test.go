package entity

import (
	"testing"
)

func TestLoan_SumInvestment(t *testing.T) {
	type fields struct {
		ID                 int64
		Amount             int64
		Rate               float64
		ROI                float64
		BorrowerID         int64
		AgreementLetterURL string
		Investments        []Investment
		Approvals          []Approval
		State              State
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "Test empty",
			fields: fields{
				Investments: nil,
			},
			want: 0,
		},
		{
			name: "Test one investment",
			fields: fields{
				Investments: []Investment{
					{
						Amount: 10000,
					},
				},
			},
			want: 10000,
		},
		{
			name: "Test more than one investment",
			fields: fields{
				Investments: []Investment{
					{
						Amount: 10000,
					},
					{
						Amount: 5000,
					},
				},
			},
			want: 15000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loan{
				ID:                 tt.fields.ID,
				Amount:             tt.fields.Amount,
				Rate:               tt.fields.Rate,
				ROI:                tt.fields.ROI,
				BorrowerID:         tt.fields.BorrowerID,
				AgreementLetterURL: tt.fields.AgreementLetterURL,
				Investments:        tt.fields.Investments,
				Approvals:          tt.fields.Approvals,
				State:              tt.fields.State,
			}
			if got := l.SumInvestment(); got != tt.want {
				t.Errorf("Loan.SumInvestment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoan_AddInvestment(t *testing.T) {
	type fields struct {
		ID                 int64
		Amount             int64
		Rate               float64
		ROI                float64
		BorrowerID         int64
		AgreementLetterURL string
		Investments        []Investment
		Approvals          []Approval
		State              State
	}
	type args struct {
		i Investment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  bool
	}{
		{
			name: "Test investment under total",
			fields: fields{
				Amount: 10000,
			},
			args: args{
				Investment{
					Amount: 5000,
				},
			},
			want:  true,
			want1: false,
		},
		{
			name: "Test investment reach total",
			fields: fields{
				Amount: 10000,
			},
			args: args{
				Investment{
					Amount: 10000,
				},
			},
			want:  true,
			want1: true,
		},
		{
			name: "Test investment exceed total",
			fields: fields{
				Amount: 10000,
			},
			args: args{
				Investment{
					Amount: 10001,
				},
			},
			want:  false,
			want1: false,
		},
		{
			name: "Test multiple investment under total",
			fields: fields{
				Amount: 5000,
				Investments: []Investment{
					{
						Amount: 3000,
					},
				},
			},
			args: args{
				Investment{
					Amount: 1000,
				},
			},
			want:  true,
			want1: false,
		},
		{
			name: "Test multiple investment reach total",
			fields: fields{
				Amount: 5000,
				Investments: []Investment{
					{
						Amount: 3000,
					},
				},
			},
			args: args{
				Investment{
					Amount: 2000,
				},
			},
			want:  true,
			want1: true,
		},
		{
			name: "Test investment exceed total",
			fields: fields{
				Amount: 5000,
				Investments: []Investment{
					{
						Amount: 3000,
					},
				},
			},
			args: args{
				Investment{
					Amount: 2001,
				},
			},
			want:  false,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loan{
				ID:                 tt.fields.ID,
				Amount:             tt.fields.Amount,
				Rate:               tt.fields.Rate,
				ROI:                tt.fields.ROI,
				BorrowerID:         tt.fields.BorrowerID,
				AgreementLetterURL: tt.fields.AgreementLetterURL,
				Investments:        tt.fields.Investments,
				Approvals:          tt.fields.Approvals,
				State:              tt.fields.State,
			}
			got, got1 := l.AddInvestment(tt.args.i)
			if got != tt.want {
				t.Errorf("Loan.AddInvestment() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Loan.AddInvestment() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoan_AddApproval(t *testing.T) {
	type fields struct {
		ID                 int64
		Amount             int64
		Rate               float64
		ROI                float64
		BorrowerID         int64
		AgreementLetterURL string
		Investments        []Investment
		Approvals          []Approval
		State              State
	}
	type args struct {
		a Approval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Test add approval from empty",
			fields: fields{
				Approvals: nil,
			},
			args: args{
				Approval{
					EmployeeID: 1,
				},
			},
			want: 1,
		},
		{
			name: "Test add approval from existing one approval",
			fields: fields{
				Approvals: []Approval{
					{
						EmployeeID: 10,
					},
				},
			},
			args: args{
				Approval{
					EmployeeID: 1,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loan{
				ID:                 tt.fields.ID,
				Amount:             tt.fields.Amount,
				Rate:               tt.fields.Rate,
				ROI:                tt.fields.ROI,
				BorrowerID:         tt.fields.BorrowerID,
				AgreementLetterURL: tt.fields.AgreementLetterURL,
				Investments:        tt.fields.Investments,
				Approvals:          tt.fields.Approvals,
				State:              tt.fields.State,
			}
			l.AddApproval(tt.args.a)
			if got := len(l.Approvals); got != tt.want {
				t.Errorf("Loan.AddApproval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
