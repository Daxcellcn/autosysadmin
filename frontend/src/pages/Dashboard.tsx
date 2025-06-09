// frontend/src/pages/Dashboard.tsx
import React, { useEffect } from 'react';
import { useAuthStore } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';
import ServerStatus from '../components/ServerStatus';
import { getServers } from '../services/api';

const Dashboard: React.FC = () => {
  const { user } = useAuthStore();
  const { showAlert } = useAlert();
  const [servers, setServers] = React.useState<any[]>([]);
  const [loading, setLoading] = React.useState(true);

  useEffect(() => {
    const fetchServers = async () => {
      try {
        const data = await getServers();
        setServers(data);
      } catch (error) {
        showAlert('Failed to fetch servers', 'error');
      } finally {
        setLoading(false);
      }
    };

    fetchServers();
  }, [showAlert]);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="dashboard-page">
      <h1>Welcome back, {user?.email}</h1>
      <div className="servers-grid">
        {servers.map((server) => (
          <ServerStatus key={server.id} server={server} />
        ))}
      </div>
    </div>
  );
};

export default Dashboard;