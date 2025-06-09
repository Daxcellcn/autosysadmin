// frontend/src/context/AuthContext.tsx
import { create } from 'zustand';
import { login as authLogin } from '../services/auth';

interface AuthState {
  user: {
    id: string;
    email: string;
    roles: string[];
    settings?: any;
  } | null;
  token: string | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => void;
}

const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  loading: true,
  login: async (email, password) => {
    set({ loading: true });
    try {
      const response = await authLogin(email, password);
      if (response.success) {
        set({
          user: response.data.user,
          token: response.data.token,
          loading: false,
        });
        localStorage.setItem('authToken', response.data.token);
      } else {
        set({ loading: false });
        throw new Error(response.error);
      }
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },
  logout: () => {
    localStorage.removeItem('authToken');
    set({ user: null, token: null });
  },
  checkAuth: () => {
    const token = localStorage.getItem('authToken');
    if (token) {
      // In a real app, you would validate the token with the server
      set({
        token,
        user: {
          id: 'user-id',
          email: 'user@example.com',
          roles: ['user'],
        },
        loading: false,
      });
    } else {
      set({ loading: false });
    }
  },
}));

export default useAuthStore;