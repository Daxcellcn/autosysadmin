// frontend/src/components/AlertBanner.tsx
import React from 'react';
import { useAlert } from '../context/AlertContext';

const AlertBanner: React.FC = () => {
  const { alert } = useAlert();

  if (!alert) return null;

  return (
    <div className={`alert-banner ${alert.type}`}>
      <div className="alert-content">
        <span>{alert.message}</span>
      </div>
    </div>
  );
};

export default AlertBanner;