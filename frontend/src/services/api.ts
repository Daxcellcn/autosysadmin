// frontend/src/services/api.ts
import axios from 'axios';
import { useAuthStore } from '../context/AuthContext';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_URL,
});

// Add a request interceptor to include the auth token
api.interceptors.request.use((config) => {
  const { token } = useAuthStore.getState();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add a response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
    }
    return Promise.reject(error);
  }
);

export const login = async (email: string, password: string) => {
  const response = await api.post('/auth/login', { email, password });
  return response.data;
};

export const getServers = async () => {
  const response = await api.get('/agents');
  return response.data.agents;
};

export const runServerCommand = async (serverId: string, command: string) => {
  const response = await api.post(`/agents/${serverId}/command`, { command });
  return response.data;
};

export const getBillingPlans = async () => {
  const response = await api.get('/billing/plans');
  return response.data.plans;
};

export const getCurrentPlan = async () => {
  const response = await api.get('/billing/subscription');
  return response.data.subscription;
};

export const getPaymentHistory = async () => {
  const response = await api.get('/billing/history');
  return response.data.payments;
};

export const processPayment = async (cardDetails: any) => {
  const response = await api.post('/billing/payment', cardDetails);
  return response.data;
};

export const updateUserSettings = async (settings: any) => {
  const response = await api.put('/user/settings', settings);
  return response.data;
};

export default api;