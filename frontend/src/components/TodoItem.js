window.TodoItem = function TodoItem({ todo, onUpdate, onDelete }) {
    const [loading, setLoading] = React.useState(false);

    const handleToggle = async () => {
        setLoading(true);
        try {
            await onUpdate(todo.id, { completed: !todo.completed });
        } catch (error) {
            console.error('Failed to toggle todo:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async () => {
        if (window.confirm('Are you sure you want to delete this todo?')) {
            setLoading(true);
            try {
                await onDelete(todo.id);
            } catch (error) {
                console.error('Failed to delete todo:', error);
                setLoading(false);
            }
        }
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { 
            hour: '2-digit', 
            minute: '2-digit' 
        });
    };

    return (
        <div className={`todo-item ${todo.completed ? 'completed' : ''}`}>
            <div className="todo-header">
                <div 
                    className={`todo-checkbox ${todo.completed ? 'checked' : ''}`}
                    onClick={handleToggle}
                    style={{ opacity: loading ? 0.5 : 1 }}
                >
                    {todo.completed && '✓'}
                </div>
                <div className={`todo-title ${todo.completed ? 'completed' : ''}`}>
                    {todo.title}
                </div>
            </div>
            
            {todo.description && (
                <div className="todo-description">
                    {todo.description}
                </div>
            )}
            
            <div className="todo-meta">
                <div className="todo-date">
                    Created: {formatDate(todo.created_at)}
                </div>
                <button 
                    className="btn btn-danger"
                    onClick={handleDelete}
                    disabled={loading}
                >
                    {loading ? 'Deleting...' : 'Delete'}
                </button>
            </div>
        </div>
    );
};
