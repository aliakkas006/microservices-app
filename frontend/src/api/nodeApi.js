import axios from 'axios';

const NODE_API_BASE = process.env.REACT_APP_NODE_API_BASE;

export const fetchUsers = () => axios.get(`${NODE_API_BASE}/users`);

export const fetchUserById = (id) => axios.get(`${NODE_API_BASE}/users/${id}`);

export const createUser = (user) => axios.post(`${NODE_API_BASE}/users`, user);
