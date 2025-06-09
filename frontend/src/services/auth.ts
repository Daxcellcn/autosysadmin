// frontend/src/services/auth.ts
import { login as apiLogin } from './api';

export const login = async (email: string, password: string) => {
  try {
    const data = await apiLogin(email, password);
    return {
      success: true,
      data,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Login failed',
    };
  }
};

export const refreshToken = async () => {
  // Implementation would call the refresh token endpoint
};

export const logout = async () => {
  // Implementation would call the logout endpoint
};