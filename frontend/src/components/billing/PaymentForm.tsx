// frontend/src/components/billing/PaymentForm.tsx
import React, { useState } from 'react';
import { toast } from 'react-toastify';
import { useBilling } from '../../context/BillingContext';

const PaymentForm: React.FC = () => {
  const { processPayment } = useBilling();
  const [cardDetails, setCardDetails] = useState({
    number: '',
    expiry: '',
    cvc: '',
    name: '',
  });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await processPayment(cardDetails);
      toast.success('Payment processed successfully');
      setCardDetails({
        number: '',
        expiry: '',
        cvc: '',
        name: '',
      });
    } catch (error) {
      toast.error('Payment failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCardDetails(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  return (
    <form onSubmit={handleSubmit} className="payment-form">
      <div className="form-group">
        <label htmlFor="cardNumber">Card Number</label>
        <input
          type="text"
          id="cardNumber"
          name="number"
          value={cardDetails.number}
          onChange={handleChange}
          placeholder="1234 5678 9012 3456"
          required
        />
      </div>
      <div className="form-group">
        <label htmlFor="cardName">Name on Card</label>
        <input
          type="text"
          id="cardName"
          name="name"
          value={cardDetails.name}
          onChange={handleChange}
          placeholder="John Doe"
          required
        />
      </div>
      <div className="form-row">
        <div className="form-group">
          <label htmlFor="cardExpiry">Expiry Date</label>
          <input
            type="text"
            id="cardExpiry"
            name="expiry"
            value={cardDetails.expiry}
            onChange={handleChange}
            placeholder="MM/YY"
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="cardCvc">CVC</label>
          <input
            type="text"
            id="cardCvc"
            name="cvc"
            value={cardDetails.cvc}
            onChange={handleChange}
            placeholder="123"
            required
          />
        </div>
      </div>
      <button type="submit" disabled={loading} className="submit-payment-btn">
        {loading ? 'Processing...' : 'Submit Payment'}
      </button>
    </form>
  );
};

export default PaymentForm;