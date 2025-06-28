function App() {
    const [todos, setTodos] = React.useState([]);
    const [loading, setLoading] = React.useState(true);
    const [error, setError] = React.useState(null);

    const loadTodos = async () => {
        try {
            setError(null);
            const data = await API.getTodos();
            setTodos(data || []);
        } catch (error) {
            setError('Failed to load todos. Please check if the backend is running.');
            console.error('Failed to load todos:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleAddTodo = async (todoData) => {
        try {
            const newTodo = await API.createTodo(todoData);
            setTodos(prev => [newTodo, ...prev]);
            setError(null);
        } catch (error) {
            setError('Failed to add todo. Please try again.');
            throw error;
        }
    };

    const handleUpdateTodo = async (id, updates) => {
        try {
            const updatedTodo = await API.updateTodo(id, updates);
            setTodos(prev => prev.map(todo => 
                todo.id === id ? updatedTodo : todo
            ));
            setError(null);
        } catch (error) {
            setError('Failed to update todo. Please try again.');
            throw error;
        }
    };

    const handleDeleteTodo = async (id) => {
        try {
            await API.deleteTodo(id);
            setTodos(prev => prev.filter(todo => todo.id !== id));
            setError(null);
        } catch (error) {
            setError('Failed to delete todo. Please try again.');
            throw error;
        }
    };

    React.useEffect(() => {
        loadTodos();
    }, []);

    return (
        <div className="container">
            <div className="header">
                <h1>Daily Todo</h1>
                <p>Organize your daily tasks efficiently</p>
            </div>
            
            <div className="content">
                {error && (
                    <div className="error">
                        {error}
                    </div>
                )}
                
                <AddTodo onAdd={handleAddTodo} />
                
                <TodoList 
                    todos={todos}
                    loading={loading}
                    onUpdate={handleUpdateTodo}
                    onDelete={handleDeleteTodo}
                />
            </div>
        </div>
    );
}

ReactDOM.render(<App />, document.getElementById('root'));
