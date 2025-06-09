// frontend/src/services/billing.ts
import {
  getBillingPlans as apiGetPlans,
  getCurrentPlan as apiGetCurrentPlan,
  getPaymentHistory as apiGetPaymentHistory,
  processPayment as apiProcessPayment,
} from './api';

export const getBillingPlans = async () => {
  try {
    const plans = await apiGetPlans();
    return {
      success: true,
      plans,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to fetch plans',
    };
  }
};

export const getCurrentPlan = async () => {
  try {
    const plan = await apiGetCurrentPlan();
    return {
      success: true,
      plan,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to fetch current plan',
    };
  }
};

export const getPaymentHistory = async () => {
  try {
    const payments = await apiGetPaymentHistory();
    return {
      success: true,
      payments,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to fetch payment history',
    };
  }
};

export const processPayment = async (cardDetails: any) => {
  try {
    const result = await apiProcessPayment(cardDetails);
    return {
      success: true,
      result,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Payment failed',
    };
  }
};