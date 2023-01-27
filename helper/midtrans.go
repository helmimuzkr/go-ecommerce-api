package helper

import (
	"e-commerce-api/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewSnapMidtrans() snap.Client {
	s := snap.Client{}
	s.New(config.SERVER_KEY, midtrans.Sandbox)

	return s
}
