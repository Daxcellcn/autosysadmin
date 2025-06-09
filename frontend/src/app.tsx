// frontend/src/App.tsx
import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useAuthStore } from './context/AuthContext';
import { AlertProvider } from './context/AlertContext';
import { BillingProvider } from './context/BillingContext';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Servers from './pages/Servers';
import Billing from './pages/Billing';
import Settings from './pages/Settings';
import Navbar from './components/Navbar';
import Sidebar from './components/Sidebar';
import LoadingSpinner from './components/LoadingSpinner';
import AlertBanner from './components/AlertBanner';
import CookieConsent from './components/CookieConsent';
import './styles/main.css';
import './styles/dark.css';
import './styles/light.css';

const App: React.FC = () => {
  const { user, loading, checkAuth } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  if (loading) {
    return <LoadingSpinner fullPage />;
  }

  return (
    <BrowserRouter>
      <AlertProvider>
        <BillingProvider>
          <div className="app-container">
            {user && <Navbar />}
            <div className="content-container">
              {user && <Sidebar />}
              <main className="main-content">
                <AlertBanner />
                <Routes>
                  <Route path="/login" element={<Login />} />
                  <Route path="/" element={user ? <Dashboard /> : <Login />} />
                  <Route path="/servers" element={user ? <Servers /> : <Login />} />
                  <Route path="/billing" element={user ? <Billing /> : <Login />} />
                  <Route path="/settings" element={user ? <Settings /> : <Login />} />
                </Routes>
              </main>
            </div>
            <CookieConsent />
            <ToastContainer position="bottom-right" autoClose={5000} />
          </div>
        </BillingProvider>
      </AlertProvider>
    </BrowserRouter>
  );
};

export default App;