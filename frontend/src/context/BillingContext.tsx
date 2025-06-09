// frontend/src/context/BillingContext.tsx
import React, { createContext, useContext, useEffect, useState } from 'react';
import {
  getBillingPlans,
  getCurrentPlan,
  getPaymentHistory,
  processPayment as apiProcessPayment,
} from '../services/billing';

interface BillingContextType {
  plans: any[];
  currentPlan: any;
  payments: any[];
  fetchPlans: () => Promise<void>;
  fetchCurrentPlan: () => Promise<void>;
  fetchPayments: () => Promise<void>;
  processPayment: (cardDetails: any) => Promise<void>;
}

const BillingContext = createContext<BillingContextType>({
  plans: [],
  currentPlan: null,
  payments: [],
  fetchPlans: async () => {},
  fetchCurrentPlan: async () => {},
  fetchPayments: async () => {},
  processPayment: async () => {},
});

export const BillingProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [plans, setPlans] = useState<any[]>([]);
  const [currentPlan, setCurrentPlan] = useState<any>(null);
  const [payments, setPayments] = useState<any[]>([]);

  const fetchPlans = async () => {
    const response = await getBillingPlans();
    if (response.success) {
      setPlans(response.plans);
    }
  };

  const fetchCurrentPlan = async () => {
    const response = await getCurrentPlan();
    if (response.success) {
      setCurrentPlan(response.plan);
    }
  };

  const fetchPayments = async () => {
    const response = await getPaymentHistory();
    if (response.success) {
      setPayments(response.payments);
    }
  };

  const processPayment = async (cardDetails: any) => {
    const response = await apiProcessPayment(cardDetails);
    if (response.success) {
      await fetchPayments(); // Refresh payment history
      return response.result;
    } else {
      throw new Error(response.error);
    }
  };

  useEffect(() => {
    fetchPlans();
    fetchCurrentPlan();
    fetchPayments();
  }, []);

  return (
    <BillingContext.Provider
      value={{
        plans,
        currentPlan,
        payments,
        fetchPlans,
        fetchCurrentPlan,
        fetchPayments,
        processPayment,
      }}
    >
      {children}
    </BillingContext.Provider>
  );
};

export const useBilling = () => useContext(BillingContext);