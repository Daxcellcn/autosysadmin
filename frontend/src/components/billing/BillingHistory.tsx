// frontend/src/components/billing/BillingHistory.tsx
import React from 'react';
import { useBilling } from '../../context/BillingContext';

const BillingHistory: React.FC = () => {
  const { payments } = useBilling();

  return (
    <div className="billing-history">
      <h3>Payment History</h3>
      {payments.length === 0 ? (
        <p>No payment history found</p>
      ) : (
        <table className="payments-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Amount</th>
              <th>Status</th>
              <th>Invoice</th>
            </tr>
          </thead>
          <tbody>
            {payments.map((payment) => (
              <tr key={payment.id}>
                <td>{new Date(payment.date).toLocaleDateString()}</td>
                <td>${payment.amount.toFixed(2)}</td>
                <td>
                  <span className={`status-badge ${payment.status}`}>
                    {payment.status}
                  </span>
                </td>
                <td>
                  <button className="invoice-btn">Download</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default BillingHistory;