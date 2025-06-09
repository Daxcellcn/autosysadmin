// frontend/src/components/LoadingSpinner.tsx
import React from 'react';

interface LoadingSpinnerProps {
  fullPage?: boolean;
  size?: 'small' | 'medium' | 'large';
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ fullPage = false, size = 'medium' }) => {
  const sizeClasses = {
    small: 'w-6 h-6 border-2',
    medium: 'w-8 h-8 border-4',
    large: 'w-12 h-12 border-4',
  };

  if (fullPage) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
        <div className={`animate-spin rounded-full border-t-transparent ${sizeClasses[size]}`}></div>
      </div>
    );
  }

  return <div className={`animate-spin rounded-full border-t-transparent ${sizeClasses[size]}`}></div>;
};

export default LoadingSpinner;