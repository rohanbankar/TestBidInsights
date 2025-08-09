import { apiClient } from './client';
import type { 
  PlatformStats, 
  ContentHealth, 
  VideoHealth, 
  DashboardSummary,
  ReportsResponse 
} from '../types/reports';

export const reportsApi = {
  getDashboard: (): Promise<DashboardSummary> =>
    apiClient.get('/api/reports/dashboard'),

  getPlatformStats: (params: { start: string; end: string }): Promise<ReportsResponse<PlatformStats>> =>
    apiClient.get('/api/reports/platform', params),

  getContentHealth: (params: { 
    platform: string; 
    start: string; 
    end: string 
  }): Promise<ReportsResponse<ContentHealth>> =>
    apiClient.get('/api/reports/content', params),

  getVideoHealth: (params: { 
    platform: string; 
    start: string; 
    end: string 
  }): Promise<ReportsResponse<VideoHealth>> =>
    apiClient.get('/api/reports/video', params),
};