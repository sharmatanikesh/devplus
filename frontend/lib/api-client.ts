import axios, { AxiosInstance, AxiosRequestConfig, AxiosError, AxiosResponse } from 'axios';
import { API_BASE_URL, API_ENDPOINTS } from './constants';
import type { ApiResponse, User } from './types';

// Create Axios instance with default config
const axiosInstance: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Critical for handling cookies
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30s timeout
});

// Request interceptor (optional: for logging or attaching headers if needed later)
axiosInstance.interceptors.request.use(
  (config) => {
    // You can attach tokens here if you were using Bearer tokens instead of cookies
    // const token = localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN);
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for global error handling
axiosInstance.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response) {
      // Server responded with non-2xx code
      if (error.response.status === 401) {
        // Handle Unauthorized (e.g., redirect to login)
        // window.location.href = '/login'; // Careful with server-side rendering
        console.warn('Unauthorized access. Session might have expired.');
      }
    }
    return Promise.reject(error);
  }
);

// Helper to standardise responses
const handleResponse = <T>(response: AxiosResponse<T>): ApiResponse<T> => {
  return {
    success: true,
    data: response.data,
  };
};

// Helper to handle errors
const handleError = <T>(error: any): ApiResponse<T> => {
  const isAxiosError = axios.isAxiosError(error);
  return {
    success: false,
    error: {
      code: isAxiosError ? `HTTP_${error.response?.status}` : 'UNKNOWN_ERROR',
      message: isAxiosError ? error.response?.statusText || error.message : 'An unexpected error occurred',
      details: isAxiosError ? error.response?.data : error,
    },
  } as ApiResponse<T>;
};

class ApiClient {
  // Generic Request Methods
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response = await axiosInstance.get<T>(url, config);
      return handleResponse(response);
    } catch (error) {
      return handleError<T>(error);
    }
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response = await axiosInstance.post<T>(url, data, config);
      return handleResponse(response);
    } catch (error) {
      return handleError<T>(error);
    }
  }

  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response = await axiosInstance.put<T>(url, data, config);
      return handleResponse(response);
    } catch (error) {
      return handleError<T>(error);
    }
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response = await axiosInstance.delete<T>(url, config);
      return handleResponse(response);
    } catch (error) {
      return handleError<T>(error);
    }
  }

  // Domain Specific Methods

  // Auth
  auth = {
    login: () => {
      // Redirect to backend login endpoint
      // Note: We use window.location because this starts the OAuth flow, it's not an AJAX call
      if (typeof window !== 'undefined') {
        window.location.href = `${API_BASE_URL}${API_ENDPOINTS.AUTH_GITHUB_CONNECT}`;
      }
    },
    logout: () => this.post('/auth/logout'),
    me: () => this.get<User>(API_ENDPOINTS.AUTH_ME),
  };

  // Repositories
  repos = {
    list: () => this.get<any[]>(API_ENDPOINTS.REPOS_LIST),
    syncAll: () => this.post<any[]>('/v1/repos/sync'),
    syncOne: (id: string) => this.post<any[]>(`/v1/repos/${id}/sync`),
    analyze: (id: string) => this.post<{ status: string }>(API_ENDPOINTS.REPOS_ANALYZE(id)),
    get: (id: string) => this.get(API_ENDPOINTS.REPOS_DETAIL(id)),
    getPullRequests: (owner: string, repo: string) => this.get(`/v1/repos/${owner}/${repo}/pulls`),
  };

  // Pull Requests
  pullRequests = {
    get: (repoId: string, prNumber: number) => this.get(API_ENDPOINTS.PR_DETAIL(repoId, prNumber)),
    analyze: (repoId: string, prNumber: number) => this.post<{ status: string }>(API_ENDPOINTS.PR_ANALYZE(repoId, prNumber)),
  };

  // Metrics
  metrics = {
    list: () => this.get(API_ENDPOINTS.METRICS),
  };

  // Dashboard
  dashboard = {
    stats: () => this.get<any>('/v1/dashboard/stats'),
    recentActivity: () => this.get<any[]>('/v1/dashboard/recent-prs'),
  };
}

export const apiClient = new ApiClient();
export default apiClient;
