// frontend/src/components/CookieConsent.tsx
import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';

const CookieConsent: React.FC = () => {
  const [consentGiven, setConsentGiven] = useState<boolean>(() => {
    return localStorage.getItem('cookieConsent') === 'given';
  });

  useEffect(() => {
    if (!consentGiven) {
      const timer = setTimeout(() => {
        toast.info(
          <div>
            <p>We use cookies to enhance your experience.</p>
            <button 
              onClick={() => {
                localStorage.setItem('cookieConsent', 'given');
                setConsentGiven(true);
                toast.dismiss();
              }}
              className="cookie-btn"
            >
              I Understand
            </button>
          </div>,
          {
            position: 'bottom-right',
            autoClose: false,
            closeOnClick: false,
            draggable: false,
            closeButton: false,
          }
        );
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [consentGiven]);

  return null;
};

export default CookieConsent;