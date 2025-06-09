// frontend/src/types/auth.ts
export interface User {
  id: string;
  email: string;
  roles: string[];
  settings?: UserSettings;
}

export interface UserSettings {
  theme: 'light' | 'dark';
  notifications: boolean;
  twoFactor: boolean;
}

export interface AuthResponse {
  user: User;
  token: string;
  refreshToken: string;
}