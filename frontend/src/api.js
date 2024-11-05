import axios from 'axios';

// Create an Axios instace with url to backend
const API = axios.create({ baseURL: 'http://localhost:9000' });

// API requests to get/add/update and delete users
export const fetchUsers = () => API.get('/users');
export const addUser = (user) => API.post('/users', user);
export const updateUser = (id, user) => API.put(`/users/${id}`, user);
export const deleteUser = (id) => API.delete(`/users/${id}`);
