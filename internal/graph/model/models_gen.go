// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Auth struct {
	Token string `json:"token"`
}

type CreateCustomerInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Address  string `json:"address"`
}

type CreateProductInput struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type Customer struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"Address"`
	UserID  int64    `json:"userId"`
	Orders  []*Order `json:"orders"`
}

type Order struct {
	ID         int64       `json:"id"`
	CustomerID int64       `json:"customerId"`
	Status     OrderStatus `json:"status"`
	TotalSum   int64       `json:"totalSum"`
	Products   []*Product  `json:"products"`
}

type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type OrderStatus string

const (
	OrderStatusCreated  OrderStatus = "CREATED"
	OrderStatusDone     OrderStatus = "DONE"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

var AllOrderStatus = []OrderStatus{
	OrderStatusCreated,
	OrderStatusDone,
	OrderStatusCanceled,
}

func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusCreated, OrderStatusDone, OrderStatusCanceled:
		return true
	}
	return false
}

func (e OrderStatus) String() string {
	return string(e)
}

func (e *OrderStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderStatus", str)
	}
	return nil
}

func (e OrderStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleAdmin    Role = "ADMIN"
	RoleCustomer Role = "CUSTOMER"
)

var AllRole = []Role{
	RoleAdmin,
	RoleCustomer,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleCustomer:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}