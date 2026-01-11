function ProjectCard({ project, onEdit, onDelete, onViewTasks }) {
  const taskCount = project.tasks_count || 0;

  return (
    <div className="card hover:shadow-lg transition-shadow">
      <div className="flex justify-between items-start mb-4">
        <h3 className="text-xl font-bold flex-1">{project.name}</h3>
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

      {project.description && (
        <p className="text-gray-600 mb-4 text-sm">{project.description}</p>
      )}

      <div className="mb-4 p-3 bg-gray-50 rounded">
        <p className="text-sm text-gray-600">
          <strong>{taskCount}</strong> {taskCount === 1 ? 'task' : 'tasks'}
        </p>
      </div>

      <button onClick={onViewTasks} className="btn-primary w-full">
        View Tasks →
      </button>
    </div>
  );
}

export default ProjectCard;
