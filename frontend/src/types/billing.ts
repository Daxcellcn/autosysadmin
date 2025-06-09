// frontend/src/types/billing.ts
export interface BillingPlan {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  features: string[];
}

export interface Subscription {
  id: string;
  planId: string;
  status: 'active' | 'canceled' | 'expired';
  startDate: string;
  endDate: string;
  renewalDate: string;
}

export interface Payment {
  id: string;
  amount: number;
  currency: string;
  status: 'completed' | 'pending' | 'failed';
  date: string;
  invoiceUrl: string;
}