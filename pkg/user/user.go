package user

import (
	"github.com/xyield/xrpl-go/model/transactions/types"
)

type Publisher interface {
	Uid() types.Address
	PublicKey() string
	Host() string
}

type Consumer interface {
	Uid() types.Address
}
