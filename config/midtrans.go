package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewSnapMidtrans(c AppConfig) snap.Client {
	s := snap.Client{}
	s.New(c.SERVER_KEY, midtrans.Sandbox)

	return s
}
