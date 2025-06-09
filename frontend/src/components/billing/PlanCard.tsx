// frontend/src/components/billing/PlanCard.tsx
import React from 'react';

interface Plan {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  features: string[];
}

interface PlanCardProps {
  plan: Plan;
  currentPlan?: boolean;
  onSelect?: (planId: string) => void;
}

const PlanCard: React.FC<PlanCardProps> = ({ plan, currentPlan = false, onSelect }) => {
  return (
    <div className={`plan-card ${currentPlan ? 'current' : ''}`}>
      <h3>{plan.name}</h3>
      <div className="plan-price">
        {plan.price > 0 ? (
          <>
            <span className="amount">{plan.currency}{plan.price}</span>
            <span className="period">/month</span>
          </>
        ) : (
          <span className="amount">Free</span>
        )}
      </div>
      <p className="plan-description">{plan.description}</p>
      <ul className="plan-features">
        {plan.features.map((feature, index) => (
          <li key={index}>{feature}</li>
        ))}
      </ul>
      {onSelect && (
        <button
          onClick={() => onSelect(plan.id)}
          className="select-plan-btn"
          disabled={currentPlan}
        >
          {currentPlan ? 'Current Plan' : 'Select Plan'}
        </button>
      )}
    </div>
  );
};

export default PlanCard;