// frontend/src/pages/Billing.tsx
import React, { useEffect } from 'react';
import { useBilling } from '../context/BillingContext';
import PlanCard from '../components/billing/PlanCard';
import PaymentForm from '../components/billing/PaymentForm';
import BillingHistory from '../components/billing/BillingHistory';

const Billing: React.FC = () => {
  const { plans, currentPlan, fetchPlans, fetchCurrentPlan, fetchPayments } = useBilling();

  useEffect(() => {
    fetchPlans();
    fetchCurrentPlan();
    fetchPayments();
  }, [fetchPlans, fetchCurrentPlan, fetchPayments]);

  const handlePlanSelect = async (planId: string) => {
    // In a real app, this would call the API to change the plan
    console.log('Selected plan:', planId);
  };

  return (
    <div className="billing-page">
      <h1>Billing & Plans</h1>
      <div className="billing-content">
        <div className="plans-section">
          <h2>Available Plans</h2>
          <div className="plans-grid">
            {plans.map((plan) => (
              <PlanCard
                key={plan.id}
                plan={plan}
                currentPlan={currentPlan?.id === plan.id}
                onSelect={handlePlanSelect}
              />
            ))}
          </div>
        </div>
        <div className="payment-section">
          <h2>Payment Method</h2>
          <PaymentForm />
        </div>
        <div className="history-section">
          <BillingHistory />
        </div>
      </div>
    </div>
  );
};

export default Billing;