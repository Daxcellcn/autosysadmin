// frontend/src/pages/Servers.tsx
import React, { useEffect, useState } from 'react';
import { useAlert } from '../context/AlertContext';
import ServerStatus from '../components/ServerStatus';
import { getServers, runServerCommand } from '../services/api';

const Servers: React.FC = () => {
  const { showAlert } = useAlert();
  const [servers, setServers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCommand, setSelectedCommand] = useState('');
  const [selectedServer, setSelectedServer] = useState('');

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

  const handleRunCommand = async () => {
    if (!selectedServer || !selectedCommand) return;

    try {
      await runServerCommand(selectedServer, selectedCommand);
      showAlert('Command executed successfully', 'success');
    } catch (error) {
      showAlert('Failed to execute command', 'error');
    }
  };

  if (loading) {
    return <div>Loading servers...</div>;
  }

  return (
    <div className="servers-page">
      <h1>Server Management</h1>
      <div className="server-actions">
        <select
          value={selectedServer}
          onChange={(e) => setSelectedServer(e.target.value)}
          className="server-select"
        >
          <option value="">Select a server</option>
          {servers.map((server) => (
            <option key={server.id} value={server.id}>
              {server.name}
            </option>
          ))}
        </select>
        <select
          value={selectedCommand}
          onChange={(e) => setSelectedCommand(e.target.value)}
          className="command-select"
        >
          <option value="">Select a command</option>
          <option value="restart">Restart Server</option>
          <option value="update">Update Packages</option>
          <option value="backup">Run Backup</option>
          <option value="status">Check Status</option>
        </select>
        <button
          onClick={handleRunCommand}
          disabled={!selectedServer || !selectedCommand}
          className="run-command-btn"
        >
          Run Command
        </button>
      </div>
      <div className="servers-grid">
        {servers.map((server) => (
          <ServerStatus key={server.id} server={server} />
        ))}
      </div>
    </div>
  );
};

export default Servers;