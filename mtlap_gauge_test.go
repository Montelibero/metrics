package metrics

import (
	"testing"

	horizonclient "github.com/stellar/go/clients/horizonclient"
	horizon "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/base"
)

func TestMTLAPGauge_Update(t *testing.T) {
	const defaultBalance = "1.0000000"

	defaultAccountBalances := []horizon.Balance{
		{
			Balance: defaultBalance,
			Asset: base.Asset{
				Code:   mtlapAsset,
				Issuer: mtlapIssuer,
			},
		},
	}

	tests := []struct {
		expectedGaugeParams MTLAPGaugeParams
		name                string
		accs                []horizon.Account
	}{
		{
			name: "CouncilReady",
			accs: []horizon.Account{
				{
					Data: map[string]string{
						dataKeyMTLACouncil: base64Encode(dataValueReady),
					},
					Balances: defaultAccountBalances,
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				IsCouncilReady: true,
				MTLAPCount:     defaultBalance,
			},
		},
		{
			name: "CouncilDelegation",
			accs: []horizon.Account{
				{
					Data: map[string]string{
						dataKeyMTLACouncil: base64Encode(mtlapIssuer),
					},
					Balances: defaultAccountBalances,
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				IsCouncilDelegation: true,
				MTLAPCount:          defaultBalance,
			},
		},
		{
			name: "AssemblyDelegation",
			accs: []horizon.Account{
				{
					Data: map[string]string{
						dataKeyMTLAAssembly: base64Encode(mtlapIssuer),
					},
					Balances: defaultAccountBalances,
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				IsAssemblyDelegation: true,
				MTLAPCount:           defaultBalance,
			},
		},
		{
			name: "FundDelegation",
			accs: []horizon.Account{
				{
					Data: map[string]string{
						dataKeyMTLFund: base64Encode(mtlapIssuer),
					},
					Balances: defaultAccountBalances,
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				IsFundDelegation: true,
				MTLAPCount:       defaultBalance,
			},
		},
		{
			name: "BSNBasicFilled",
			accs: []horizon.Account{
				{
					Data: map[string]string{
						dataKeyName:    base64Encode("Stas Karkavin"),
						dataKeyAbout:   base64Encode("Software Engineer"),
						dataKeyWebsite: base64Encode("https://xdefrag.dev"),
					},
					Balances: defaultAccountBalances,
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				IsBSNBasicFilled: true,
				MTLAPCount:       defaultBalance,
			},
		},
		{
			name: "BSNPartialyFilled",
			accs: []horizon.Account{
				{
					Data: map[string]string{
						dataKeyName: base64Encode("Stas Karkavin"),
					},
					Balances: defaultAccountBalances,
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				IsBSNPartialyFilled: true,
				MTLAPCount:          defaultBalance,
			},
		},
		{
			name: "MTLAPCount",
			accs: []horizon.Account{
				{
					Balances: []horizon.Balance{
						{
							Balance: "23.0000000",
							Asset: base.Asset{
								Code:   mtlapAsset,
								Issuer: mtlapIssuer,
							},
						},
					},
				},
			},
			expectedGaugeParams: MTLAPGaugeParams{
				MTLAPCount: "23.0000000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewMockLogger(t)
			cl := NewMockAccounter(t)
			m := NewMockMetricer(t)

			cl.EXPECT().Accounts(horizonclient.AccountsRequest{
				Asset: mtlapAssetRequest,
				Limit: defaultLimit,
			}).Return(horizon.AccountsPage{
				Embedded: struct {
					Records []horizon.Account "json:\"records\""
				}{
					Records: tt.accs,
				},
			}, nil)

			m.EXPECT().MTLAPGaugeReset()
			m.EXPECT().MTLAPGaugeInc(tt.expectedGaugeParams)

			mtlapGauge := NewMTLAPGauge(l, cl, m)
			mtlapGauge.Update()
		})
	}
}
