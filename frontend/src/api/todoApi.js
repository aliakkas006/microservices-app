import axios from 'axios';

const GO_API_BASE = process.env.REACT_APP_GO_API_BASE;

export const fetchTodos = () => axios.get(`${GO_API_BASE}/todos`);

console.log(fetchTodos);

export const fetchTodoById = (id) => axios.get(`${GO_API_BASE}/todos/${id}`);

export const createTodo = (todo) => axios.post(`${GO_API_BASE}/todos`, todo);

export const updateTodo = (id, todo) =>
  axios.put(`${GO_API_BASE}/todos/${id}`, todo);

export const deleteTodo = (id) => axios.delete(`${GO_API_BASE}/todos/${id}`);
