import axios from 'axios';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE,
  timeout: 10000,
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle errors
api.interceptors.response.use(
  (response) => {
    const body = response.data;
    console.log('API Raw Response:', body);
    
    // Check if response is wrapped with status/message
    if (body && typeof body === 'object') {
      // If it has a 'data' field, extract it (could be array or single object)
      if (body.data !== undefined) {
        console.log('API Unwrapped Data:', body.data);
        return body.data;
      }
      // If it has 'status' field but no 'data', return the whole thing
      if (body.status) {
        return body;
      }
    }
    
    // Return as-is if not wrapped
    console.log('API Final Response:', body);
    return body;
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth APIs
export const authAPI = {
  signup: (data) =>
    api.post('/auth/signup', data),
  login: (data) =>
    api.post('/auth/login', data),
  getProfile: () =>
    api.get('/auth/me'),
  updateProfile: (data) =>
    api.put('/auth/me', data),
};

// User APIs
export const userAPI = {
  list: () =>
    api.get('/users'),
};

// Project APIs
export const projectAPI = {
  getAll: () =>
    api.get('/projects'),
  get: (id) =>
    api.get(`/projects/${id}`),
  create: (data) =>
    api.post('/projects', data),
  update: (id, data) =>
    api.put(`/projects/${id}`, data),
  delete: (id) =>
    api.delete(`/projects/${id}`),
};

// Task APIs
export const taskAPI = {
  getByID: (taskId) =>
    api.get(`/tasks/${taskId}`),
  getByProject: (projectId, status = '', priority = '') =>
    api.get(`/projects/${projectId}/tasks`, {
      params: { status, priority },
    }),
  getAssignedToMe: (status = '', priority = '') =>
    api.get('/tasks/assigned', {
      params: { status, priority },
    }),
  create: (projectId, data) =>
    api.post(`/projects/${projectId}/tasks`, data),
  update: (taskId, data) =>
    api.put(`/tasks/${taskId}`, data),
  updateStatus: (taskId, status) =>
    api.patch(`/tasks/${taskId}/status`, { status }),
  updatePriority: (taskId, priority) =>
    api.patch(`/tasks/${taskId}/priority`, { priority }),
  updateAssignee: (taskId, assignee_id) =>
    api.patch(`/tasks/${taskId}/assignee`, { assignee_id }),
  assign: (taskId, data) =>
    api.post(`/tasks/${taskId}/assign`, data),
  delete: (taskId) =>
    api.delete(`/tasks/${taskId}`),
};

// Comment APIs
export const commentAPI = {
  getByTask: (taskId) =>
    api.get(`/tasks/${taskId}/comments`),
  getRecent: (page = 1, pageSize = 10) =>
    api.get('/comments/recent', {
      params: { page, page_size: pageSize },
    }),
  create: (taskId, data) =>
    api.post(`/tasks/${taskId}/comments`, data),
  update: (commentId, data) =>
    api.put(`/comments/${commentId}`, data),
  delete: (commentId) =>
    api.delete(`/comments/${commentId}`),
};

export default api;
