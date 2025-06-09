// frontend/src/context/AlertContext.tsx
import React, { createContext, useContext, useState } from 'react';

interface Alert {
  message: string;
  type: 'success' | 'error' | 'info' | 'warning';
}

interface AlertContextType {
  alert: Alert | null;
  showAlert: (message: string, type: Alert['type']) => void;
  clearAlert: () => void;
}

const AlertContext = createContext<AlertContextType>({
  alert: null,
  showAlert: () => {},
  clearAlert: () => {},
});

export const AlertProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [alert, setAlert] = useState<Alert | null>(null);

  const showAlert = (message: string, type: Alert['type']) => {
    setAlert({ message, type });
    setTimeout(() => {
      setAlert(null);
    }, 5000);
  };

  const clearAlert = () => {
    setAlert(null);
  };

  return (
    <AlertContext.Provider value={{ alert, showAlert, clearAlert }}>
      {children}
    </AlertContext.Provider>
  );
};

export const useAlert = () => useContext(AlertContext);