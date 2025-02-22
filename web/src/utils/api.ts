import axios from "axios";
import { BookData } from "./schema";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  login: (data: { email: string }) => api.post("/auth/login", data),
  register: (data: { name: string; email: string; contact: string }) =>
    api.post("/auth/register", data),
  getCurrentUser: () => api.get("/me"),
  logout: () => api.get("/auth/logout"),
};

export const bookAPI = {
  getAllBooks: () => api.get("/books"),
  addBook: (data: BookData) => api.post("/books", data),
  updateBook: (id: string, data: Partial<BookData>) =>
    api.put(`/books/${id}`, data),
  deleteBook: (id: string) => api.delete(`/books/${id}`),
  searchBooks: (params: { searchString: string; field: string }) =>
    api.post("/books/search", params),
};

export const issueAPI = {
  createRequest: (bookId: string) => api.post(`/issues/request/${bookId}`),
  approveRequest: (issueId: string) => api.put(`/issues/approve/${issueId}`),
  getRequests: () => api.get("/issues/requests"),
};

export default api;
