import axios, { type AxiosInstance, type InternalAxiosRequestConfig } from 'axios';

const DEFAULT_API_PATH = '/api/v1';

const resolveApiBaseUrl = (): string => {
  const envUrl = import.meta.env.VITE_API_URL?.trim();
  if (envUrl) {
    return envUrl.endsWith('/') ? envUrl.slice(0, -1) : envUrl;
  }

  if (typeof window === 'undefined') {
    return DEFAULT_API_PATH;
  }

  const { protocol, hostname, port } = window.location;
  const normalizedPort = port || (protocol === 'https:' ? '443' : '80');

  // If we're already on the API port (e.g. reverse proxy), reuse the current origin + path.
  if (normalizedPort === '8080' || normalizedPort === '80' || normalizedPort === '443' || normalizedPort === '') {
    const portSegment =
      normalizedPort === '80' || normalizedPort === '443' || normalizedPort === '' ? '' : `:${normalizedPort}`;
    return `${protocol}//${hostname}${portSegment}${DEFAULT_API_PATH}`;
  }

  // Fall back to the common local setup where the backend listens on 8080.
  return `${protocol}//${hostname}:8080${DEFAULT_API_PATH}`;
};

const API_BASE_URL = resolveApiBaseUrl();

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('auth_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export interface Movie {
  id: string;
  title: string;
  release_year: number;
  director: string;
  duration: number;
  duration_minutes?: number;
  description: string;
  average_rating: number;
  genres: Genre[];
}

export interface Genre {
  id: string;
  name: string;
}

export interface Review {
  id: string;
  movie_id: string;
  user_id: string;
  rating: number;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
  movie?: {
    id: string;
    title: string;
  };
  user?: {
    id: string;
    username: string;
    email: string;
  };
}

export interface User {
  id: string;
  email: string;
  username: string;
  role: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface CreateMovieRequest {
  title: string;
  release_year: number;
  director: string;
  duration: number;
  duration_minutes?: number;
  description: string;
  genre_ids: string[];
}

export interface CreateReviewRequest {
  rating: number;
  title: string;
  content: string;
}

export interface CreateGenreRequest {
  name: string;
}

export const authAPI = {
  login: (data: LoginRequest) => apiClient.post('/auth/login', data),
  register: (data: RegisterRequest) => apiClient.post('/auth/register', data),
  getProfile: () => apiClient.get('/users/me'),
};

export const moviesAPI = {
  getAll: (params?: {
    genre_id?: string;
    year?: number;
    min_rating?: number;
    search?: string;
    page?: number;
    limit?: number;
  }) => apiClient.get('/movies', { params }),
  getById: (id: string) => apiClient.get(`/movies/${id}`),
  create: (data: CreateMovieRequest) => apiClient.post('/movies', data),
  update: (id: string, data: Partial<CreateMovieRequest>) =>
    apiClient.put(`/movies/${id}`, data),
  delete: (id: string) => apiClient.delete(`/movies/${id}`),
};

export const reviewsAPI = {
  getAll: (params?: { page?: number; limit?: number }) =>
    apiClient.get('/reviews', { params }),
  getByMovie: (movieId: string, params?: { page?: number; limit?: number }) =>
    apiClient.get(`/movies/${movieId}/reviews`, { params }),
  getMyReviews: (params?: { page?: number; limit?: number }) =>
    apiClient.get('/users/me/reviews', { params }),
  create: (movieId: string, data: CreateReviewRequest) =>
    apiClient.post(`/movies/${movieId}/reviews`, data),
  update: (id: string, data: Partial<CreateReviewRequest>) =>
    apiClient.put(`/reviews/${id}`, data),
  delete: (id: string) => apiClient.delete(`/reviews/${id}`),
};

export const genresAPI = {
  getAll: () => apiClient.get('/genres'),
  create: (data: CreateGenreRequest) => apiClient.post('/genres', data),
  update: (id: string, data: Partial<CreateGenreRequest>) =>
    apiClient.put(`/genres/${id}`, data),
  delete: (id: string) => apiClient.delete(`/genres/${id}`),
};

export default apiClient;
