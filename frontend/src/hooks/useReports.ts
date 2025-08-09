import { useQuery } from '@tanstack/react-query';
import { reportsApi } from '../api/reports';

export function useDashboard() {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: reportsApi.getDashboard,
    refetchInterval: 30000, // Refresh every 30 seconds
  });
}

export function usePlatformStats(startDate: string, endDate: string) {
  return useQuery({
    queryKey: ['platformStats', startDate, endDate],
    queryFn: () => reportsApi.getPlatformStats({ start: startDate, end: endDate }),
    enabled: !!startDate && !!endDate,
    refetchInterval: 30000,
  });
}

export function useContentHealth(platform: string, startDate: string, endDate: string) {
  return useQuery({
    queryKey: ['contentHealth', platform, startDate, endDate],
    queryFn: () => reportsApi.getContentHealth({ platform, start: startDate, end: endDate }),
    enabled: !!platform && !!startDate && !!endDate,
    refetchInterval: 30000,
  });
}

export function useVideoHealth(platform: string, startDate: string, endDate: string) {
  return useQuery({
    queryKey: ['videoHealth', platform, startDate, endDate],
    queryFn: () => reportsApi.getVideoHealth({ platform, start: startDate, end: endDate }),
    enabled: !!platform && !!startDate && !!endDate,
    refetchInterval: 30000,
  });
}