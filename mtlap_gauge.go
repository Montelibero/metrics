package metrics

import (
	"log/slog"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

const (
	defaultLimit      = 20
	mtlapAsset        = "MTLAP"
	mtlapIssuer       = "GCNVDZIHGX473FEI7IXCUAEXUJ4BGCKEMHF36VYP5EMS7PX2QBLAMTLA"
	mtlapAssetRequest = "MTLAP:GCNVDZIHGX473FEI7IXCUAEXUJ4BGCKEMHF36VYP5EMS7PX2QBLAMTLA"
)

type MTLAPGauge struct {
	l  Logger
	cl Accounter
	m  Metricer
}

func NewMTLAPGauge(l Logger, cl Accounter, m Metricer) *MTLAPGauge {
	return &MTLAPGauge{
		l:  l,
		cl: cl,
		m:  m,
	}
}

func (g *MTLAPGauge) fetchAllAccounts() ([]horizon.Account, error) {
	var allAccounts []horizon.Account
	accs, err := g.cl.Accounts(horizonclient.AccountsRequest{
		Asset: mtlapAssetRequest,
		Limit: defaultLimit,
	})
	if err != nil {
		return nil, err
	}

	if len(accs.Embedded.Records) < defaultLimit {
		return accs.Embedded.Records, nil
	}

	for {
		allAccounts = append(allAccounts, accs.Embedded.Records...)
		accs, err = g.cl.NextAccountsPage(accs)
		if err != nil {
			return nil, err
		}
		if len(accs.Embedded.Records) == 0 {
			break
		}
	}
	return allAccounts, nil
}

func (g *MTLAPGauge) Update() {
	accounts, err := g.fetchAllAccounts()
	if err != nil {
		g.l.Error("[mtlap_total] failed to get accounts", slog.Any("error", err))
		return
	}

	g.m.MTLAPGaugeReset()

	for _, acc := range accounts {
		p := MTLAPGaugeParams{
			isCouncilReady: isDataEqual(acc, dataKeyMTLACouncil, dataValueReady),
			isCouncilDelegation: !isDataEqual(acc, dataKeyMTLACouncil, dataValueReady) &&
				isDataExist(acc, dataKeyMTLACouncil),
			isAssemblyDelegation: isDataExist(acc, dataKeyMTLAAssembly),
			isFundDelegation:     isDataExist(acc, dataKeyMTLFund),
			isBSNBasicFilled: isDataExist(acc, dataKeyName) &&
				isDataExist(acc, dataKeyAbout) &&
				isDataExist(acc, dataKeyWebsite),
			mtlapCount: acc.GetCreditBalance(mtlapAsset, mtlapIssuer),
		}

		g.m.MTLAPGaugeInc(p)
	}
}
