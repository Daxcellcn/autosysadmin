// frontend/src/types/server.ts
export interface Server {
  id: string;
  name: string;
  hostname: string;
  ipAddress: string;
  os: string;
  architecture: string;
  status: 'online' | 'offline' | 'degraded';
  lastHeartbeat: string;
  tags: string[];
  cpuUsage: number[];
  memoryUsage: number[];
  responseTimes: number[];
}

export interface ServerCommand {
  command: string;
  args?: string[];
  timeout?: number;
}

export interface ServerStats {
  cpu: number;
  memory: number;
  disk: number;
  networkIn: number;
  networkOut: number;
  processes: ProcessStats[];
}

export interface ProcessStats {
  pid: number;
  name: string;
  cpuUsage: number;
  memoryUsage: number;
}