package metrics

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

const (
	dataKeyMTLFund      = "mtl_delegate"
	dataKeyMTLAAssembly = "mtla_a_delegate"
	dataKeyMTLACouncil  = "mtla_c_delegate"
	dataKeyName         = "Name"
	dataKeyAbout        = "About"
	dataKeyWebsite      = "Website"
	dataValueReady      = "ready"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Accounter interface {
	Accounts(request horizonclient.AccountsRequest) (accounts horizon.AccountsPage, err error)
	NextAccountsPage(page horizon.AccountsPage) (accounts horizon.AccountsPage, err error)
}

type Metricer interface {
	MTLAPGaugeReset()
	MTLAPGaugeInc(p MTLAPGaugeParams)
}
