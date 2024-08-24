package app

import (
	"fmt"

	"route256/cart/internal/app/definitions"
)

type (
	Options struct {
		Addr, ProductToken, ProductAddr, LOMSAddr, JaegerAddr string
	}

	configProductService struct {
		productToken, productAddr string
	}

	path struct {
		cartItemAdd, cartItemDelete, cartDelete, cartList, cartCheckout, metrics string
	}

	Config struct {
		addr string
		configProductService
		lomsAddr   string
		jaegerAddr string
		path       path
	}
)

func NewConfig(opts *Options) *Config {
	return &Config{
		addr: opts.Addr,
		configProductService: configProductService{
			productToken: opts.ProductToken,
			productAddr:  opts.ProductAddr,
		},
		lomsAddr:   opts.LOMSAddr,
		jaegerAddr: opts.JaegerAddr,
		path: path{
			cartItemAdd:    fmt.Sprintf("POST /user/{%s}/cart/{%s}", definitions.ParamUserID, definitions.ParamSkuID),
			cartItemDelete: fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", definitions.ParamUserID, definitions.ParamSkuID),
			cartDelete:     fmt.Sprintf("DELETE /user/{%s}/cart/", definitions.ParamUserID),
			cartList:       fmt.Sprintf("GET /cart/{%s}/list/", definitions.ParamUserID),
			cartCheckout:   "POST /cart/checkout",
			metrics:        "GET /metrics",
		},
	}
}
