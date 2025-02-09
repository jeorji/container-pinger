// src/api/containersApi.ts
import axios from 'axios';

export interface Ping {
  ip: string;
  last_ping_latency: number;
  last_ping_time: string;
}

export interface Container {
  id: string;
  name: string;
  image: string;
  state: 'running' | 'stopped';
  status: string;
  created_at: string;
  pings: Ping[] | null;
}

export const fetchContainers = async (): Promise<Container[]> => {
  try {
    const response = await axios.get<Container[]>('http://localhost:8080/api/containers');
    return response.data;
  } catch (error) {
    console.error('Ошибка при получении контейнеров:', error);
    throw error;
  }
};
