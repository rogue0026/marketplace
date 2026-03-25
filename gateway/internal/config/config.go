package config

import (
	"errors"
	"os"
)

var (
	ErrUsersAddrEmpty      = errors.New("env USERS_ADDR is empty")
	ErrProductsAddrEmpty   = errors.New("env PRODUCTS_ADDR is empty")
	ErrOrdersAddrEmpty     = errors.New("env ORDERS_ADDR is empty")
	ErrHttpServerAddrEmpty = errors.New("env HTTP_SERVER_ADDR is empty")
)

type GRPCClientsConfig struct {
	UsersClientAddr    string
	ProductsClientAddr string
	OrdersClientAddr   string
}

func LoadGRPCClientsConfig() (*GRPCClientsConfig, error) {
	users := os.Getenv("USERS_ADDR")
	if len(users) == 0 {
		return nil, ErrUsersAddrEmpty
	}

	products := os.Getenv("PRODUCTS_ADDR")
	if len(products) == 0 {
		return nil, ErrProductsAddrEmpty
	}

	orders := os.Getenv("ORDERS_ADDR")
	if len(orders) == 0 {
		return nil, ErrOrdersAddrEmpty
	}

	cfg := &GRPCClientsConfig{
		UsersClientAddr:    users,
		ProductsClientAddr: products,
		OrdersClientAddr:   orders,
	}

	return cfg, nil
}

type HTTPServerConfig struct {
	Addr string
}

func LoadHttpServerConfig() (*HTTPServerConfig, error) {
	addr := os.Getenv("HTTP_SERVER_ADDR")
	if len(addr) == 0 {
		return nil, ErrHttpServerAddrEmpty
	}

	cfg := &HTTPServerConfig{
		Addr: addr,
	}

	return cfg, nil
}
