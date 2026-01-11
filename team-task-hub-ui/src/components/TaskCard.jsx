const TASK_STATUSES = ['OPEN', 'IN_PROGRESS', 'DONE'];

function TaskCard({ task, onEdit, onDelete, onStatusChange }) {

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
            className="text-gray-400 hover:text-gray-600"
          >
            ✏️
          </button>
          <button
            onClick={onDelete}
            className="text-gray-400 hover:text-red-600"
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

      <div className="flex flex-wrap gap-3 mb-4">
        <span className="text-xs text-gray-600">{task.priority}</span>
        <span className="text-xs text-gray-600">{task.status}</span>
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
