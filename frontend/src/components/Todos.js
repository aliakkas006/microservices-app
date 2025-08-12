import React, { useEffect, useState } from 'react';
import { fetchTodos, createTodo, updateTodo, deleteTodo } from '../api/todoApi';

export default function Todos() {
  const [todos, setTodos] = useState([]);
  const [form, setForm] = useState({ title: '', description: '' });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadTodos();
  }, []);

  const loadTodos = async () => {
    setLoading(true);
    try {
      const res = await fetchTodos();
      setTodos(res.data);
    } catch (err) {
      alert('Failed to load todos');
    }
    setLoading(false);
  };

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!form.title) {
      alert('Title is required');
      return;
    }
    try {
      await createTodo(form);
      setForm({ title: '', description: '' });
      loadTodos();
    } catch {
      alert('Failed to create todo');
    }
  };

  const toggleComplete = async (todo) => {
    try {
      await updateTodo(todo.id, { ...todo, completed: !todo.completed });
      loadTodos();
    } catch {
      alert('Failed to update todo');
    }
  };

  const handleDelete = async (id) => {
    try {
      await deleteTodo(id);
      loadTodos();
    } catch {
      alert('Failed to delete todo');
    }
  };

  return (
    <div>
      <h2>To-Dos</h2>
      {loading ? (
        <p>Loading...</p>
      ) : (
        <ul>
          {todos.map((todo) => (
            <li
              key={todo.id}
              style={{
                textDecoration: todo.completed ? 'line-through' : 'none',
              }}
            >
              <input
                type="checkbox"
                checked={todo.completed}
                onChange={() => toggleComplete(todo)}
              />
              <strong>{todo.title}</strong> - {todo.description}
              <button
                onClick={() => handleDelete(todo.id)}
                style={{ marginLeft: 10 }}
              >
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
      <form onSubmit={handleSubmit}>
        <input
          name="title"
          value={form.title}
          onChange={handleChange}
          placeholder="Title"
          required
        />
        <input
          name="description"
          value={form.description}
          onChange={handleChange}
          placeholder="Description"
        />
        <button type="submit">Add To-Do</button>
      </form>
    </div>
  );
}
