package utils

import (
	"math/big"
)

var (
	Wei      = big.NewInt(1)                   //nolint:nolintlint,gochecknoglobals,gomnd
	Kwei     = big.NewInt(1000)                //nolint:nolintlint,gochecknoglobals,gomnd
	Mwei     = big.NewInt(1000000)             //nolint:nolintlint,gochecknoglobals,gomnd
	Gwei     = big.NewInt(1000000000)          //nolint:nolintlint,gochecknoglobals,gomnd
	Microeth = big.NewInt(1000000000000)       //nolint:nolintlint,gochecknoglobals,gomnd
	Millieth = big.NewInt(1000000000000000)    //nolint:nolintlint,gochecknoglobals,gomnd
	Eth      = big.NewInt(1000000000000000000) //nolint:nolintlint,gochecknoglobals,gomnd
	Eth_1K   = Eth.Mul(Kwei, Eth)              //nolint:nolintlint,gochecknoglobals,gomnd
)
