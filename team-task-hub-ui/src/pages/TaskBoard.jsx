import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAsync } from '../hooks/useAsync';
import { taskAPI, projectAPI, userAPI } from '../api/client';
import Loading from '../components/Loading';
import ErrorMessage from '../components/ErrorMessage';
import TaskForm from '../components/TaskForm';
import Modal from '../components/Modal';
import ConfirmDialog from '../components/ConfirmDialog';

const TASK_STATUSES = ['OPEN', 'IN_PROGRESS', 'DONE'];
const TASK_PRIORITIES = ['LOW', 'MEDIUM', 'HIGH'];

function TaskBoard() {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const [tasks, setTasks] = useState([]);
  const [project, setProject] = useState(null);
  const [users, setUsers] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [editingTask, setEditingTask] = useState(null);
  const [editingProjectName, setEditingProjectName] = useState(false);
  const [editingProjectDesc, setEditingProjectDesc] = useState(false);
  const [projectName, setProjectName] = useState('');
  const [projectDesc, setProjectDesc] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('created_at');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage] = useState(10);
  const [filters, setFilters] = useState({
    status: '',
    priority: '',
  });
  const [deleteConfirm, setDeleteConfirm] = useState({ isOpen: false, taskId: null, taskName: '' });
  const [deleteProjectConfirm, setDeleteProjectConfirm] = useState({ isOpen: false, projectName: '' });

  const { execute: fetchProject, status: projectStatus, error: projectError } =
    useAsync(async () => {
      const data = await projectAPI.get(projectId);
      return data;
    });

  const { execute: fetchTasks, status: tasksStatus, error: tasksError } =
    useAsync(async () => {
      const data = await taskAPI.getByProject(projectId);
      return data || [];
    });

  const { execute: fetchUsers } = useAsync(async () => {
    const data = await userAPI.list();
    return Array.isArray(data) ? data : data.data || [];
  });

  useEffect(() => {
    fetchProject().then((data) => {
      if (data) {
        setProject(data);
        setProjectName(data.name);
        setProjectDesc(data.description || '');
      }
    });
    fetchTasks().then((data) => {
      if (data) setTasks(data);
    });
    fetchUsers().then((data) => {
      if (data) setUsers(data);
    });
  }, [projectId]);

  const filteredTasks = tasks.filter((task) => {
    if (filters.status && task.status !== filters.status) return false;
    if (filters.priority && task.priority !== filters.priority) return false;
    if (searchTerm && !task.title.toLowerCase().includes(searchTerm.toLowerCase())) {
      return false;
    }
    return true;
  });

  const sortedTasks = [...filteredTasks].sort((a, b) => {
    let aVal = a[sortField];
    let bVal = b[sortField];

    if (typeof aVal === 'string') aVal = aVal.toLowerCase();
    if (typeof bVal === 'string') bVal = bVal.toLowerCase();

    if (sortOrder === 'asc') {
      return aVal > bVal ? 1 : -1;
    } else {
      return aVal < bVal ? 1 : -1;
    }
  });

  const totalPages = Math.ceil(sortedTasks.length / itemsPerPage);
  const paginatedTasks = sortedTasks.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const handleCreateTask = async (formData) => {
    try {
      const newTask = await taskAPI.create(projectId, formData);
      setTasks([...tasks, newTask]);
      setShowForm(false);
    } catch (err) {
      console.error('Failed to create task:', err);
    }
  };

  const handleUpdateTask = async (taskId, formData) => {
    try {
      const updated = await taskAPI.update(taskId, formData);
      setTasks(tasks.map((t) => (t.id === taskId ? updated : t)));
      setEditingTask(null);
    } catch (err) {
      console.error('Failed to update task:', err);
    }
  };

  const handleDeleteTask = async (taskId) => {
    const task = tasks.find(t => t.id === taskId);
    setDeleteConfirm({
      isOpen: true,
      taskId: taskId,
      taskName: task?.title || 'Task'
    });
  };

  const handleConfirmDeleteTask = async () => {
    try {
      await taskAPI.delete(deleteConfirm.taskId);
      setTasks(tasks.filter((t) => t.id !== deleteConfirm.taskId));
      setDeleteConfirm({ isOpen: false, taskId: null, taskName: '' });
    } catch (err) {
      console.error('Failed to delete task:', err);
    }
  };

  const handleUpdateStatus = async (taskId, newStatus) => {
    try {
      await taskAPI.updateStatus(taskId, newStatus);
      setTasks(
        tasks.map((t) =>
          t.id === taskId ? { ...t, status: newStatus } : t
        )
      );
    } catch (err) {
      console.error('Failed to update task status:', err);
    }
  };

  const handleUpdatePriority = async (taskId, newPriority) => {
    try {
      await taskAPI.updatePriority(taskId, newPriority);
      setTasks(
        tasks.map((t) =>
          t.id === taskId ? { ...t, priority: newPriority } : t
        )
      );
    } catch (err) {
      console.error('Failed to update task priority:', err);
    }
  };

  const handleUpdateAssignee = async (taskId, newAssigneeId) => {
    try {
      const assigneeId = newAssigneeId ? parseInt(newAssigneeId) : null;
      await taskAPI.updateAssignee(taskId, assigneeId);
      
      const task = tasks.find(t => t.id === taskId);
      const assignee = assigneeId ? users.find(u => u.id === assigneeId) : null;
      
      setTasks(
        tasks.map((t) =>
          t.id === taskId ? { ...t, assignee_id: assigneeId, assignee } : t
        )
      );
    } catch (err) {
      console.error('Failed to update task assignee:', err);
    }
  };

  const handleUpdateProject = async (formData) => {
    try {
      const updated = await projectAPI.update(projectId, formData);
      setProject(updated);
    } catch (err) {
      console.error('Failed to update project:', err);
    }
  };

  const handleSaveProjectName = async () => {
    if (projectName.trim() && projectName !== project?.name) {
      await handleUpdateProject({ 
        name: projectName,
        description: projectDesc 
      });
    } else {
      setProjectName(project?.name || '');
    }
    setEditingProjectName(false);
  };

  const handleSaveProjectDesc = async () => {
    if (projectDesc !== project?.description) {
      await handleUpdateProject({ 
        name: projectName,
        description: projectDesc 
      });
    }
    setEditingProjectDesc(false);
  };

  const handleDeleteProject = () => {
    setDeleteProjectConfirm({
      isOpen: true,
      projectName: project?.name || 'Project'
    });
  };

  const handleConfirmDeleteProject = async () => {
    try {
      await projectAPI.delete(projectId);
      setDeleteProjectConfirm({ isOpen: false, projectName: '' });
      navigate('/projects');
    } catch (err) {
      console.error('Failed to delete project:', err);
    }
  };

  const handleSort = (field) => {
    if (sortField === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortOrder('asc');
    }
    setCurrentPage(1);
  };

  const SortHeader = ({ field, label }) => (
    <th
      className="px-6 py-3 text-left text-sm font-medium text-slate-700 cursor-pointer hover:bg-slate-100 transition-colors"
      onClick={() => handleSort(field)}
    >
      <div className="flex items-center gap-2">
        {label}
        {sortField === field && (
          <span className="text-slate-500">{sortOrder === 'asc' ? '↑' : '↓'}</span>
        )}
      </div>
    </th>
  );

  if (projectStatus === 'pending' || tasksStatus === 'pending') {
    return <Loading />;
  }

  return (
    <div className="min-h-screen bg-white">
      <div className="container pt-8 pb-12">
        <button
          onClick={() => navigate('/projects')}
          className="text-sm font-medium text-slate-600 hover:text-slate-900 mb-8 transition-colors"
        >
          ← Back to Projects
        </button>

        {/* Header Section */}
        <div className="mb-10">
          {/* Title Row */}
          <div className="flex justify-between items-center gap-8 mb-6">
            <div className="flex-1">
              {editingProjectName ? (
                <input
                  type="text"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  onBlur={handleSaveProjectName}
                  onKeyDown={(e) => e.key === 'Enter' && handleSaveProjectName()}
                  autoFocus
                  className="text-3xl font-semibold text-slate-900 bg-slate-50 px-3 py-2 rounded border border-slate-300 focus:outline-none w-full"
                />
              ) : (
                <h1
                  onClick={() => setEditingProjectName(true)}
                  className="text-3xl font-semibold text-slate-900 cursor-text"
                  title="Click to edit"
                >
                  {project?.name || 'Tasks'}
                </h1>
              )}
            </div>

            {/* Action Buttons */}
            <div className="flex gap-3 items-center flex-shrink-0">
              <button
                onClick={() => setShowForm(!showForm)}
                className="btn-primary whitespace-nowrap"
              >
                {showForm ? 'Cancel' : '+ New Task'}
              </button>
              <button
                onClick={handleDeleteProject}
                className="btn-secondary text-sm border-slate-300 text-slate-700 hover:bg-slate-100"
              >
                Delete
              </button>
            </div>
          </div>

          {/* Description Box */}
          <div className="mb-6">
            {editingProjectDesc ? (
              <textarea
                value={projectDesc}
                onChange={(e) => setProjectDesc(e.target.value)}
                onBlur={handleSaveProjectDesc}
                onKeyDown={(e) => e.key === 'Escape' && setEditingProjectDesc(false)}
                autoFocus
                className="text-slate-600 bg-slate-50 px-3 py-2 rounded border border-slate-300 focus:outline-none w-full resize-none"
                rows="2"
              />
            ) : (
              <div className="border border-slate-200 rounded p-3 bg-slate-50">
                <p
                  onClick={() => setEditingProjectDesc(true)}
                  className={`cursor-text ${project?.description ? 'text-slate-600' : 'text-slate-400 italic'}`}
                  title="Click to add description"
                >
                  {project?.description || 'Click to add description'}
                </p>
              </div>
            )}
          </div>
        </div>

        {projectError && <ErrorMessage message={projectError.message} />}
        {tasksError && <ErrorMessage message={tasksError.message} />}

        <Modal
          isOpen={showForm}
          onClose={() => setShowForm(false)}
          title="New Task"
        >
          <TaskForm
            onSubmit={handleCreateTask}
            onCancel={() => setShowForm(false)}
          />
        </Modal>

        <Modal
          isOpen={!!editingTask}
          onClose={() => setEditingTask(null)}
          title="Edit Task"
        >
          {editingTask && (
            <TaskForm
              task={editingTask}
              onSubmit={(data) => handleUpdateTask(editingTask.id, data)}
              onCancel={() => setEditingTask(null)}
            />
          )}
        </Modal>


        {/* Tasks Table - Always show structure */}
        <>
          {/* Search and Filters - Only show when there are tasks */}
          {tasks.length > 0 && (
            <div className="mb-6 flex gap-4 items-center justify-between">
              <div className="flex gap-3 items-center">
                <span className="text-sm font-medium text-slate-700">Filters:</span>
                <select
                  value={filters.status}
                  onChange={(e) => {
                    setFilters({ ...filters, status: e.target.value });
                    setCurrentPage(1);
                  }}
                  className="input text-sm border-slate-300"
                >
                  <option value="">All Status</option>
                  {TASK_STATUSES.map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
                <select
                  value={filters.priority}
                  onChange={(e) => {
                    setFilters({ ...filters, priority: e.target.value });
                    setCurrentPage(1);
                  }}
                  className="input text-sm border-slate-300"
                >
                  <option value="">All Priorities</option>
                  {TASK_PRIORITIES.map((priority) => (
                    <option key={priority} value={priority}>
                      {priority}
                    </option>
                  ))}
                </select>
              </div>
              <input
                type="text"
                placeholder="Search tasks..."
                value={searchTerm}
                onChange={(e) => {
                  setSearchTerm(e.target.value);
                  setCurrentPage(1);
                }}
                className="input search-input"
              />
            </div>
          )}

          {/* Tasks Table */}
          <div className="overflow-x-auto mb-6 border border-slate-200 rounded-lg">
            <table className="w-full">
                <thead>
                  <tr className="bg-slate-50 border-b border-slate-200">
                    <SortHeader field="title" label="Title" />
                    <SortHeader field="status" label="Status" />
                    <SortHeader field="priority" label="Priority" />
                    <SortHeader field="due_date" label="Due Date" />
                    <SortHeader field="assignee_id" label="Assignee" />
                  </tr>
                </thead>
                <tbody>
                  {paginatedTasks.map((task) => (
                    <tr 
                      key={task.id} 
                      className="border-b border-slate-200 hover:bg-slate-50 transition-colors"
                      onClick={(e) => {
                        if (e.target.tagName !== 'SELECT') {
                          navigate(`/tasks/${task.id}`);
                        }
                      }}
                    >
                      <td className="px-6 py-4 font-medium text-slate-900 cursor-pointer">
                        {task.title}
                      </td>
                      <td className="px-6 py-4">
                        <select
                          value={task.status}
                          onChange={(e) => {
                            handleUpdateStatus(task.id, e.target.value);
                          }}
                          className={`text-xs font-semibold px-2 py-1 rounded border appearance-none cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                            task.status === 'DONE' ? 'bg-green-50 border-green-200 text-green-700' :
                            task.status === 'IN_PROGRESS' ? 'bg-yellow-50 border-yellow-200 text-yellow-700' :
                            'bg-gray-50 border-gray-200 text-gray-700'
                          }`}
                        >
                          {TASK_STATUSES.map(status => (
                            <option key={status} value={status}>{status}</option>
                          ))}
                        </select>
                      </td>
                      <td className="px-6 py-4">
                        <select
                          value={task.priority}
                          onChange={(e) => {
                            handleUpdatePriority(task.id, e.target.value);
                          }}
                          className={`text-xs font-semibold px-2 py-1 rounded border appearance-none cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                            task.priority === 'HIGH' ? 'bg-red-50 text-red-700 border-red-200' :
                            task.priority === 'MEDIUM' ? 'bg-yellow-50 text-yellow-700 border-yellow-200' :
                            'bg-blue-50 text-blue-700 border-blue-200'
                          }`}
                        >
                          {TASK_PRIORITIES.map(priority => (
                            <option key={priority} value={priority}>{priority}</option>
                          ))}
                        </select>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        {task.due_date ? new Date(task.due_date).toLocaleDateString() : '-'}
                      </td>
                      <td className="px-6 py-4 text-sm">
                        {task.assignee?.email ? (
                          <div className="flex items-center gap-2">
                            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-slate-400 to-slate-600 flex items-center justify-center text-white font-semibold text-xs">
                              {task.assignee.email.charAt(0).toUpperCase()}
                            </div>
                            <div className="text-left">
                              <p className="text-xs font-medium text-slate-900">{task.assignee.email.split('@')[0]}</p>
                              <p className="text-xs text-slate-500">{task.assignee.email.split('@')[1]}</p>
                            </div>
                          </div>
                        ) : (
                          <span className="text-slate-400">-</span>
                        )}
                      </td>
                    </tr>
                  ))}
                  {paginatedTasks.length === 0 && (
                    <tr>
                      <td colSpan="5" className="px-6 py-12 text-center text-slate-500">
                        No tasks found
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-between items-center">
                <div className="text-sm text-slate-600">
                  Showing {(currentPage - 1) * itemsPerPage + 1} to{' '}
                  {Math.min(currentPage * itemsPerPage, sortedTasks.length)} of{' '}
                  {sortedTasks.length} tasks
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                    disabled={currentPage === 1}
                    className="btn-sm bg-slate-100 text-slate-700 hover:bg-slate-200 border border-slate-200 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    ← Previous
                  </button>
                  <span className="px-4 py-2 text-sm text-slate-600">
                    Page {currentPage} of {totalPages}
                  </span>
                  <button
                    onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                    disabled={currentPage === totalPages}
                    className="btn-sm bg-slate-100 text-slate-700 hover:bg-slate-200 border border-slate-200 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Next →
                  </button>
                </div>
              </div>
            )}
          </>

      <ConfirmDialog
        isOpen={deleteConfirm.isOpen}
        title="Delete Task"
        message={`Are you sure you want to delete "${deleteConfirm.taskName}"? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleConfirmDeleteTask}
        onCancel={() => setDeleteConfirm({ isOpen: false, taskId: null, taskName: '' })}
        isDanger={true}
      />

      {/* Delete Project Confirmation */}
      <ConfirmDialog
        isOpen={deleteProjectConfirm.isOpen}
        title="Delete Project"
        message={`Are you sure you want to delete "${deleteProjectConfirm.projectName}"? All tasks in this project will also be deleted. This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleConfirmDeleteProject}
        onCancel={() => setDeleteProjectConfirm({ isOpen: false, projectName: '' })}
        isDanger={true}
      />
      </div>
    </div>
  );
}

export default TaskBoard;
