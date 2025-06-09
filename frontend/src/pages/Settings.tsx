// frontend/src/pages/Settings.tsx
import React, { useState, useEffect } from 'react';
import { useAuthStore } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';
import { updateUserSettings } from '../services/api';

const Settings: React.FC = () => {
  const { user } = useAuthStore();
  const { showAlert } = useAlert();
  const [settings, setSettings] = useState({
    theme: 'light',
    notifications: true,
    twoFactor: false,
  });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (user?.settings) {
      setSettings(user.settings);
    }
  }, [user]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setSettings(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await updateUserSettings(settings);
      showAlert('Settings updated successfully', 'success');
    } catch (error) {
      showAlert('Failed to update settings', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="settings-page">
      <h1>Settings</h1>
      <form onSubmit={handleSubmit} className="settings-form">
        <div className="form-group">
          <label htmlFor="theme">Theme</label>
          <select
            id="theme"
            name="theme"
            value={settings.theme}
            onChange={handleChange}
          >
            <option value="light">Light</option>
            <option value="dark">Dark</option>
          </select>
        </div>
        <div className="form-group checkbox-group">
          <input
            type="checkbox"
            id="notifications"
            name="notifications"
            checked={settings.notifications}
            onChange={handleChange}
          />
          <label htmlFor="notifications">Enable Email Notifications</label>
        </div>
        <div className="form-group checkbox-group">
          <input
            type="checkbox"
            id="twoFactor"
            name="twoFactor"
            checked={settings.twoFactor}
            onChange={handleChange}
          />
          <label htmlFor="twoFactor">Enable Two-Factor Authentication</label>
        </div>
        <button type="submit" disabled={loading} className="save-settings-btn">
          {loading ? 'Saving...' : 'Save Settings'}
        </button>
      </form>
    </div>
  );
};

export default Settings;