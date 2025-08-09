import { apiClient } from './client';
import type { LoginRequest, LoginResponse, User } from '../types/auth';

export const authApi = {
  login: (credentials: LoginRequest): Promise<LoginResponse> =>
    apiClient.post('/api/auth/login', credentials),

  logout: (): Promise<{ message: string }> =>
    apiClient.post('/api/auth/logout'),

  refresh: (): Promise<{ access_token: string }> =>
    apiClient.post('/api/auth/refresh'),

  me: (): Promise<User> =>
    apiClient.get('/api/auth/me'),
};