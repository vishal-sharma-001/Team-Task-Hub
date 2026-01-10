const TASK_STATUSES = ['OPEN', 'IN_PROGRESS', 'DONE'];

function TaskCard({ task, onEdit, onDelete, onStatusChange }) {
  const priorityColors = {
    LOW: 'bg-green-100 text-green-800',
    MEDIUM: 'bg-yellow-100 text-yellow-800',
    HIGH: 'bg-red-100 text-red-800',
  };

  const statusColors = {
    OPEN: 'text-gray-600',
    IN_PROGRESS: 'text-blue-600',
    DONE: 'text-green-600',
  };

  const dueDate = task.due_date
    ? new Date(task.due_date).toLocaleDateString()
    : null;

  return (
    <div className="card hover:shadow-lg transition-shadow">
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-bold flex-1">{task.title}</h3>
        <div className="flex gap-2">
          <button
            onClick={onEdit}
            className="text-blue-600 hover:text-blue-800"
          >
            ✏️
          </button>
          <button
            onClick={onDelete}
            className="text-red-600 hover:text-red-800"
          >
            🗑️
          </button>
        </div>
      </div>

      {task.description && (
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">
          {task.description}
        </p>
      )}

      <div className="flex flex-wrap gap-2 mb-4">
        <span className={`px-2 py-1 rounded text-xs font-semibold ${priorityColors[task.priority]}`}>
          {task.priority}
        </span>
        <span className={`text-sm font-semibold ${statusColors[task.status]}`}>
          {task.status}
        </span>
      </div>

      {task.assignee && (
        <p className="text-sm text-gray-600 mb-2">
          <strong>Assigned to:</strong> {task.assignee}
        </p>
      )}

      {dueDate && (
        <p className="text-sm text-gray-600 mb-3">
          <strong>Due:</strong> {dueDate}
        </p>
      )}

      <div className="border-t pt-3">
        <label className="text-xs font-semibold text-gray-600 block mb-2">
          Change Status
        </label>
        <select
          value={task.status}
          onChange={(e) => onStatusChange(e.target.value)}
          className="input w-full text-sm"
        >
          {TASK_STATUSES.map((status) => (
            <option key={status} value={status}>
              {status}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
}

export default TaskCard;
