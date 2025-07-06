import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import './App.css';
import { API_URL } from './utils/env';

const EditTodoForm = ({ todo, onSave, onCancel }) => {
  const [title, setTitle] = useState(todo.title);
  const [description, setDescription] = useState(todo.description || '');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({ ...todo, title, description });
  };

  return (
    <form onSubmit={handleSubmit} className="edit-form">
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="edit-input"
        required
      />
      <textarea
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        className="edit-textarea"
        placeholder="Add description (optional)"
      />
      <div className="edit-actions">
        <button type="submit" className="save-btn">
          Save
        </button>
        <button type="button" onClick={onCancel} className="cancel-btn">
          Cancel
        </button>
      </div>
    </form>
  );
};

// PropTypes validation
EditTodoForm.propTypes = {
  todo: PropTypes.shape({
    id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
    title: PropTypes.string.isRequired,
    description: PropTypes.string,
    completed: PropTypes.bool.isRequired,
    created_at: PropTypes.string,
    updated_at: PropTypes.string
  }).isRequired,
  onSave: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired
};

function App() {
  const [todos, setTodos] = useState([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editingId, setEditingId] = useState(null);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  useEffect(() => {
    loadTodos();
  }, []);

  // Helper function to render todo list content
  const renderTodoListContent = () => {
    if (loading) {
      return (
        <div className="loading">
          <div className="spinner"></div>
          Loading your todos...
        </div>
      );
    }

    if (error) {
      return (
        <div className="empty-state">
          <div style={{ fontSize: '4em', marginBottom: '20px' }}>⚠️</div>
          <h3>Error</h3>
          <p>{error}</p>
          <button onClick={loadTodos}>Retry</button>
        </div>
      );
    }

    if (total === 0) {
      return (
        <div className="empty-state">
          <div style={{ fontSize: '4em', marginBottom: '20px' }}>📝</div>
          <h3>No tasks yet</h3>
          <p>Add your first task using the form on the right!</p>
        </div>
      );
    }

    return todos.map(todo => (
      <div key={todo.id} className={`todo-item ${todo.completed ? 'completed' : ''}`}>
        {editingId === todo.id ? (
          <EditTodoForm
            todo={todo}
            onSave={updateTodo}
            onCancel={() => setEditingId(null)}
          />
        ) : (
          <>
            <div className="todo-header">
              <input
                type="checkbox"
                className="todo-checkbox"
                checked={todo.completed}
                onChange={(e) => toggleTodo(todo.id, e.target.checked)}
              />
              <div className="todo-title">{todo.title}</div>
            </div>
            {todo.description && <div className="todo-description">{todo.description}</div>}
            <div className="todo-actions">
              <button onClick={() => setEditingId(todo.id)} className="edit-btn">
                <span>✏️</span> Edit
              </button>
              <button className="danger" onClick={() => deleteTodo(todo.id)}>
                <span>🗑️</span> Delete
              </button>
            </div>
            <div className="todo-meta">
              <span>Created: {new Date(todo.created_at).toLocaleDateString()}</span>
              <span>Updated: {new Date(todo.updated_at).toLocaleDateString()}</span>
            </div>
          </>
        )}
      </div>
    ));
  };

  const loadTodos = async () => {
    try {
      const response = await fetch(`${API_URL}/todos`);
      if (!response.ok) throw new Error('Failed to fetch todos');
      const data = await response.json();
      setTodos(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err.message);
      setTodos([]); // Ensure todos is never null
    } finally {
      setLoading(false);
    }
  };

  const addTodo = async (e) => {
    e.preventDefault();
    if (!title.trim()) return;
    try {
      const response = await fetch(`${API_URL}/todos`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, completed: false })
      });
      if (!response.ok) throw new Error('Failed to create todo');
      const newTodo = await response.json();
      setTodos(prev => [newTodo, ...prev]);
      setTitle('');
      setDescription('');
    } catch (err) {
      setError(err.message);
    }
  };

  const updateTodo = async (updatedTodo) => {
    try {
      const response = await fetch(`${API_URL}/todos/${updatedTodo.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedTodo)
      });
      if (!response.ok) throw new Error('Failed to update todo');
      const updated = await response.json();
      setTodos(prev => prev.map(t => (t.id === updatedTodo.id ? updated : t)));
      setEditingId(null);
    } catch (err) {
      setError(err.message);
    }
  };

  const toggleTodo = (id, completed) => {
    const todo = todos.find(t => t.id === id);
    if (todo) {
      updateTodo({ ...todo, completed });
    }
  };

  const deleteTodo = async (id) => {
    if (!window.confirm('Are you sure you want to delete this task?')) return;
    try {
      await fetch(`${API_URL}/todos/${id}`, { method: 'DELETE' });
      setTodos(prev => prev.filter(t => t.id !== id));
    } catch (err) {
      setError('Failed to delete todo');
    }
  };

  const total = Array.isArray(todos) ? todos.length : 0;
  const completed = Array.isArray(todos) ? todos.filter(t => t.completed).length : 0;
  const pending = total - completed;

  return (
    <div className="container">
      <div className={`main-content ${sidebarCollapsed ? 'expanded' : ''}`}>
        <h1>📝 MinimalDo</h1>

        <div className="stats">
          <div className="stat-card">
            <div className="stat-number">{total}</div>
            <div className="stat-label">Total Tasks</div>
          </div>
          <div className="stat-card">
            <div className="stat-number">{completed}</div>
            <div className="stat-label">Completed</div>
          </div>
          <div className="stat-card">
            <div className="stat-number">{pending}</div>
            <div className="stat-label">Pending</div>
          </div>
        </div>

        <div id="todoList" className="todo-list">
          {renderTodoListContent()}
        </div>
      </div>

      <div className={`sidebar ${sidebarCollapsed ? 'collapsed' : ''}`}> 
        <h2>✨ Add New Task</h2>
        <form onSubmit={addTodo} className="todo-form">
          <div className="form-group">
            <label htmlFor="title">Task Title *</label>
            <input id="title" type="text" value={title} onChange={(e) => setTitle(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea id="description" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <button type="submit">
            <span>➕</span> Add Task
          </button>
        </form>
      </div>

      <button className={`toggle-btn ${sidebarCollapsed ? 'sidebar-closed' : 'sidebar-open'}`} onClick={() => setSidebarCollapsed(!sidebarCollapsed)}>
        {sidebarCollapsed ? '➕' : '✕'}
      </button>
    </div>
  );
}

export default App;
