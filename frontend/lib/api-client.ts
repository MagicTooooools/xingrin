/**
 * API client configuration file
 * 
 * Core functionality:
 * 1. Unified HTTP request wrapper
 * 2. Unified error handling
 * 3. Request/response logging
 * 4. JWT token management
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

import axios, { AxiosRequestConfig } from 'axios';

// Token storage keys
const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

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
 * Response interceptor: Handle response data
 * 
 * Workflow:
 * 1. Log response (for development debugging)
 * 2. Return response data
 * 
 * Notes:
 * - Backend returns camelCase JSON
 * - Frontend can use directly, no additional conversion needed
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
  (error) => {
    // Only output error logs in development environment
    if (process.env.NODE_ENV === 'development') {
      // Check if it's an Axios error
      if (axios.isAxiosError(error)) {
        console.error('[ERROR] API Error:', {
          status: error.response?.status,
          statusText: error.response?.statusText,
          url: error.config?.url,
          method: error.config?.method,
          data: error.response?.data,
          message: error.message,
          code: error.code
        });
      } else if (error instanceof Error) {
        // Regular Error object
        console.error('[ERROR] API Error:', error.message, error.stack);
      } else {
        // Unknown error type
        console.error('[ERROR] API Error: Unknown error', String(error));
      }
    }

    // Handle 401 Unauthorized: clear tokens and redirect to login page
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      const url = error.config?.url || '';
      // Exclude auth-related APIs to avoid redirect loops
      const isAuthApi = url.includes('/auth/login') || 
                        url.includes('/auth/logout') || 
                        url.includes('/auth/me') ||
                        url.includes('/auth/refresh');
      
      if (!isAuthApi && typeof window !== 'undefined') {
        // Clear tokens and redirect to login
        tokenManager.clearTokens();
        window.location.href = '/login';
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
    response?: { data?: { message?: string; error?: string; detail?: string } };
    message?: string
  }

  // Request timeout (over 30 seconds)
  if (err.code === 'ECONNABORTED') {
    return 'Request timeout, please try again later';
  }

  // Backend returned error message (supports multiple formats)
  if (err.response?.data?.error) {
    return err.response.data.error;
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
