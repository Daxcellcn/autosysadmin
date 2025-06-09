// frontend/src/components/ServerStatus.tsx
import React from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface ServerStatusProps {
  server: {
    id: string;
    name: string;
    status: 'online' | 'offline' | 'degraded';
    cpuUsage: number[];
    memoryUsage: number[];
    responseTimes: number[];
    lastUpdated: string;
  };
}

const ServerStatus: React.FC<ServerStatusProps> = ({ server }) => {
  const labels = Array.from({ length: server.cpuUsage.length }, (_, i) => `${i * 5} min ago`).reverse();

  const data = {
    labels,
    datasets: [
      {
        label: 'CPU Usage %',
        data: server.cpuUsage,
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.5)',
        tension: 0.1,
      },
      {
        label: 'Memory Usage %',
        data: server.memoryUsage,
        borderColor: 'rgb(53, 162, 235)',
        backgroundColor: 'rgba(53, 162, 235, 0.5)',
        tension: 0.1,
      },
    ],
  };

  const options = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: true,
        text: `${server.name} - ${server.status.toUpperCase()}`,
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
      },
    },
  };

  return (
    <div className="server-status-card">
      <div className="status-indicator" data-status={server.status}></div>
      <div className="server-stats">
        <Line options={options} data={data} />
        <div className="server-meta">
          <span>Last updated: {new Date(server.lastUpdated).toLocaleString()}</span>
        </div>
      </div>
    </div>
  );
};

export default ServerStatus;