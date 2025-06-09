package stripe

import (
	"encoding/json"
	"net/http"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type WebhookHandler struct {
	webhookSecret string
}

func NewWebhookHandler(secret string) *WebhookHandler {
	return &WebhookHandler{
		webhookSecret: secret,
	}
}

func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	const maxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, maxBodyBytes)

	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		http.Error(w, "Webhook signature verification failed", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			http.Error(w, "Error parsing payment intent", http.StatusBadRequest)
			return
		}
		h.handlePaymentIntentSucceeded(paymentIntent)
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			http.Error(w, "Error parsing invoice", http.StatusBadRequest)
			return
		}
		h.handleInvoicePaymentSucceeded(invoice)
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			http.Error(w, "Error parsing subscription", http.StatusBadRequest)
			return
		}
		h.handleSubscriptionDeleted(subscription)
	default:
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) handlePaymentIntentSucceeded(payment stripe.PaymentIntent) {
	// Implement your payment success logic here
}

func (h *WebhookHandler) handleInvoicePaymentSucceeded(invoice stripe.Invoice) {
	// Implement your invoice payment logic here
}

func (h *WebhookHandler) handleSubscriptionDeleted(subscription stripe.Subscription) {
	// Implement your subscription cancellation logic here
}