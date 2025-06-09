package stripe

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
)

type StripeClient struct {
	apiKey string
}

func NewStripeClient(apiKey string) *StripeClient {
	stripe.Key = apiKey
	return &StripeClient{
		apiKey: apiKey,
	}
}

func (c *StripeClient) CreateCustomer(email, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	return customer.New(params)
}

func (c *StripeClient) GetCustomer(customerID string) (*stripe.Customer, error) {
	return customer.Get(customerID, nil)
}

func (c *StripeClient) UpdateCustomer(customerID string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return customer.Update(customerID, params)
}

func (c *StripeClient) DeleteCustomer(customerID string) (*stripe.Customer, error) {
	return customer.Del(customerID, nil)
}

func (c *StripeClient) CreatePaymentMethod(customerID, paymentMethodID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	return paymentmethod.Attach(paymentMethodID, params)
}