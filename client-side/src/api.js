import axios from "axios";
import { store } from "./redux/store.js";
import { baseUrl } from "./const";
const api = axios.create({ baseUrl });

// Add the token to every request if available
api.interceptors.request.use(
  (config) => {
    const state = store.getState();
    const token = state.user.userToken; // Get token from Redux
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

export default api;
