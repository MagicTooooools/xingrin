/**
 * API client configuration file
 * 
 * Core functionality:
 * 1. Unified HTTP request wrapper
 * 2. Unified error handling
 * 3. Request/response logging
 * 4. JWT token management with auto-refresh
 * 
 * Naming convention explanation:
 * - Frontend (TypeScript/React): camelCase
 *   Example: pageSize, createdAt, organizationId
 * 
 * - Backend (Go): Uses JSON tags for camelCase output
 *   Example: pageSize, createdAt, organizationId
 * 
 * - API JSON format: camelCase
 *   Example: pageSize, createdAt, organizationId
 */

import axios, { AxiosRequestConfig, AxiosError, InternalAxiosRequestConfig } from 'axios';

// Token storage keys
const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

// Track if we're currently refreshing to prevent multiple refresh calls
let isRefreshing = false;
// Queue of failed requests to retry after token refresh
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: unknown) => void;
}> = [];

/**
 * Process the queue of failed requests after token refresh
 */
const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else if (token) {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

/**
 * Token management utilities
 */
export const tokenManager = {
  getAccessToken: (): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  },
  
  getRefreshToken: (): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(REFRESH_TOKEN_KEY);
  },
  
  setTokens: (accessToken: string, refreshToken: string): void => {
    if (typeof window === 'undefined') return;
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  },

  setAccessToken: (accessToken: string): void => {
    if (typeof window === 'undefined') return;
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  },
  
  clearTokens: (): void => {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
  },
  
  hasTokens: (): boolean => {
    return !!tokenManager.getAccessToken();
  }
};

/**
 * Create axios instance
 * Configure base URL, timeout and default headers
 */
const apiClient = axios.create({
  baseURL: '/api',  // API base path
  timeout: 30000,      // 30 second timeout
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * Request interceptor: Handle preparation work before request
 * 
 * Workflow:
 * 1. Add Authorization header with JWT token
 * 2. Log request (for development debugging)
 */
apiClient.interceptors.request.use(
  (config) => {
    // Add JWT token to Authorization header
    const token = tokenManager.getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Only output debug logs in development environment
    if (process.env.NODE_ENV === 'development') {
      console.log('[REQUEST] API Request:', {
        method: config.method?.toUpperCase(),
        url: config.url,
        baseURL: config.baseURL,
        fullURL: `${config.baseURL}${config.url}`,
        data: config.data,
        params: config.params,
        hasToken: !!token
      });
    }

    return config;
  },
  (error) => {
    if (process.env.NODE_ENV === 'development') {
      console.error('[ERROR] Request Error:', error);
    }
    return Promise.reject(error);
  }
);

/**
 * Response interceptor: Handle response data and auto-refresh token
 * 
 * Workflow:
 * 1. Log response (for development debugging)
 * 2. On 401 error, try to refresh token and retry the request
 * 3. Return response data
 */
apiClient.interceptors.response.use(
  (response) => {
    // Only output debug logs in development environment
    if (process.env.NODE_ENV === 'development') {
      console.log('[RESPONSE] API Response:', {
        status: response.status,
        statusText: response.statusText,
        url: response.config.url,
        data: response.data
      });
    }

    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Only output error logs in development environment
    if (process.env.NODE_ENV === 'development') {
      console.error('[ERROR] API Error:', {
        status: error.response?.status,
        statusText: error.response?.statusText,
        url: error.config?.url,
        method: error.config?.method,
        data: error.response?.data,
        message: error.message,
        code: error.code
      });
    }

    // Handle 401 Unauthorized with auto-refresh
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      const url = originalRequest.url || '';
      
      // Don't try to refresh for auth-related APIs
      const isAuthApi = url.includes('/auth/login') || 
                        url.includes('/auth/logout') || 
                        url.includes('/auth/refresh');
      
      if (isAuthApi) {
        return Promise.reject(error);
      }

      // Check if we have a refresh token
      const refreshToken = tokenManager.getRefreshToken();
      if (!refreshToken) {
        // No refresh token, redirect to login
        tokenManager.clearTokens();
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }

      // If already refreshing, queue this request
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return apiClient(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      // Mark as retrying and start refresh
      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Call refresh token API
        const response = await axios.post('/api/auth/refresh/', {
          refreshToken: refreshToken,
        });

        const { accessToken: newAccessToken } = response.data;
        
        // Save new access token
        tokenManager.setAccessToken(newAccessToken);
        
        // Process queued requests with new token
        processQueue(null, newAccessToken);
        
        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        processQueue(refreshError, null);
        tokenManager.clearTokens();
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

// Export default axios instance (generally not used directly)
export default apiClient;

/**
 * Export common HTTP methods
 * 
 * Usage examples:
 * 
 * 1. GET request:
 *    api.get('/organizations', { 
 *      params: { pageSize: 10, sortBy: 'name' }  // Use camelCase
 *    })
 *    Backend receives: page_size=10&sort_by=name (automatically converted)
 * 
 * 2. POST request:
 *    api.post('/organizations/create', {
 *      organizationName: 'test',  // Use camelCase
 *      createdAt: '2024-01-01'
 *    })
 *    Backend receives: organization_name, created_at (automatically converted)
 * 
 * 3. Response data (already camelCase):
 *    const response = await api.get('/organizations')
 *    response.data.pageSize  // [OK] Use camelCase directly
 *    response.data.createdAt // [OK] Use camelCase directly
 * 
 * Type parameters:
 * - T: Response data type (optional)
 * - config: axios configuration object (optional)
 */
export const api = {
  /**
   * GET request
   * @param url - Request path (relative to baseURL)
   * @param config - axios config, recommend using params for query parameters
   * @returns Promise<AxiosResponse<T>>
   */
  get: <T = unknown>(url: string, config?: AxiosRequestConfig) => apiClient.get<T>(url, config),

  /**
   * POST request
   * @param url - Request path (relative to baseURL)
   * @param data - Request body data (will be automatically converted to snake_case)
   * @param config - axios config (optional)
   * @returns Promise<AxiosResponse<T>>
   */
  post: <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) => apiClient.post<T>(url, data, config),

  /**
   * PUT request
   * @param url - Request path (relative to baseURL)
   * @param data - Request body data (will be automatically converted to snake_case)
   * @param config - axios config (optional)
   * @returns Promise<AxiosResponse<T>>
   */
  put: <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) => apiClient.put<T>(url, data, config),

  /**
   * PATCH request (partial update)
   * @param url - Request path (relative to baseURL)
   * @param data - Request body data (will be automatically converted to snake_case)
   * @param config - axios config (optional)
   * @returns Promise<AxiosResponse<T>>
   */
  patch: <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) => apiClient.patch<T>(url, data, config),

  /**
   * DELETE request
   * @param url - Request path (relative to baseURL)
   * @param config - axios config (optional)
   * @returns Promise<AxiosResponse<T>>
   */
  delete: <T = unknown>(url: string, config?: AxiosRequestConfig) => apiClient.delete<T>(url, config),
};

/**
 * Error handling utility function
 * 
 * Function: Extract user-friendly error messages from error objects
 * 
 * Error priority:
 * 1. Request cancelled
 * 2. Request timeout
 * 3. Backend returned error message
 * 4. axios error message
 * 5. Unknown error
 * 
 * Usage example:
 * try {
 *   await api.get('/organizations')
 * } catch (error) {
 *   const message = getErrorMessage(error)
 *   toast.error(message)
 * }
 * 
 * @param error - Error object (can be any type)
 * @returns User-friendly error message string
 */
export const getErrorMessage = (error: unknown): string => {
  // Request was cancelled (user actively cancelled or component unmounted)
  if (axios.isCancel(error)) {
    return 'Request has been cancelled';
  }

  // Type guard: Check if it's an error object
  const err = error as {
    code?: string;
    response?: { data?: { 
      message?: string; 
      error?: string | { code?: string; message?: string; details?: Array<{ field?: string; message?: string }> }; 
      detail?: string 
    } };
    message?: string
  }

  // Request timeout (over 30 seconds)
  if (err.code === 'ECONNABORTED') {
    return 'Request timeout, please try again later';
  }

  // Backend returned error message (supports multiple formats)
  const errorData = err.response?.data?.error;
  if (errorData) {
    // New format: { error: { code, message, details } }
    if (typeof errorData === 'object') {
      // If has validation details, return first detail message
      if (errorData.details && errorData.details.length > 0) {
        const detail = errorData.details[0];
        return detail.message || errorData.message || 'Validation error';
      }
      return errorData.message || 'Unknown error';
    }
    // Old format: { error: "string" }
    return errorData;
  }
  if (err.response?.data?.message) {
    return err.response.data.message;
  }
  if (err.response?.data?.detail) {
    return err.response.data.detail;
  }

  // axios own error message
  if (err.message) {
    return err.message;
  }

  // Fallback error message
  return 'Unknown error occurred';
};
