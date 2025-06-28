window.TodoList = function TodoList({ todos, loading, onUpdate, onDelete }) {
    if (loading) {
        return <div className="loading">Loading todos...</div>;
    }

    if (todos.length === 0) {
        return (
            <div className="empty-state">
                <h3>No todos yet!</h3>
                <p>Add your first todo above to get started.</p>
            </div>
        );
    }

    const completedCount = todos.filter(todo => todo.completed).length;
    const totalCount = todos.length;

    return (
        <div>
            <div style={{ 
                marginBottom: '20px', 
                textAlign: 'center', 
                color: '#6c757d',
                fontSize: '1.1rem'
            }}>
                {completedCount} of {totalCount} tasks completed
            </div>
            <div className="todo-list">
                {todos.map(todo => (
                    <TodoItem
                        key={todo.id}
                        todo={todo}
                        onUpdate={onUpdate}
                        onDelete={onDelete}
                    />
                ))}
            </div>
        </div>
    );
};
