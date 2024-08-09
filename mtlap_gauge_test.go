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
				isCouncilReady: true,
				mtlapCount:     defaultBalance,
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
				isCouncilDelegation: true,
				mtlapCount:          defaultBalance,
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
				isAssemblyDelegation: true,
				mtlapCount:           defaultBalance,
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
				isFundDelegation: true,
				mtlapCount:       defaultBalance,
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
				isBSNBasicFilled: true,
				mtlapCount:       defaultBalance,
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
				mtlapCount: "23.0000000",
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
