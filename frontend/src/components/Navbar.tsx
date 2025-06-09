// frontend/src/components/Navbar.tsx
import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../context/AuthContext';

const Navbar: React.FC = () => {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="navbar">
      <div className="navbar-brand">Autosysadmin</div>
      <div className="navbar-items">
        <div className="navbar-user">
          <span>{user?.email}</span>
          <button onClick={handleLogout} className="logout-btn">
            Logout
          </button>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;