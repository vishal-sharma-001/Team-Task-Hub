import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAsync } from '../hooks/useAsync';
import { taskAPI, projectAPI, commentAPI } from '../api/client';
import Loading from '../components/Loading';
import ErrorMessage from '../components/ErrorMessage';

function Dashboard() {
  const navigate = useNavigate();
  const [allTasks, setAllTasks] = useState([]);
  const [projects, setProjects] = useState({});
  const [recentComments, setRecentComments] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('created_at');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage] = useState(10);
  const [filterStatus, setFilterStatus] = useState('');

  const { execute: fetchData, status, error } = useAsync(
    async () => {
      try {
        // Fetch assigned tasks
        const tasksResponse = await taskAPI.getAssignedToMe();
        const tasks = tasksResponse.data || tasksResponse || [];

        // Fetch projects
        const projectsResponse = await projectAPI.getAll();
        const projectsList = projectsResponse.data || projectsResponse || [];

        // Fetch recent comments
        const commentsResponse = await commentAPI.getRecent(1, 10);
        const comments = commentsResponse.data || commentsResponse || [];

        // Create project lookup
        const projectLookup = {};
        projectsList.forEach((p) => {
          projectLookup[p.id] = p;
        });
        setProjects(projectLookup);
        setRecentComments(comments);
        setAllTasks(Array.isArray(tasks) ? tasks : []);
        return { tasks, projects: projectsList, comments };
      } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
        throw err;
      }
    }
  );

  useEffect(() => {
    fetchData();
  }, []);

  const filteredTasks = allTasks.filter((task) => {
    if (filterStatus && task.status !== filterStatus) return false;
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
      className="px-6 py-3 text-left cursor-pointer hover:bg-gray-100"
      onClick={() => handleSort(field)}
    >
      <div className="flex items-center gap-2">
        {label}
        {sortField === field && (
          <span>{sortOrder === 'asc' ? '↑' : '↓'}</span>
        )}
      </div>
    </th>
  );

  if (status === 'pending') return <Loading />;

  return (
    <div className="min-h-screen bg-white">
      <div className="container pt-8 pb-12">
        {error && <ErrorMessage message={error.message} />}

        {/* Header */}
        <div className="mb-12">
          <h1 className="text-4xl font-bold text-gray-900">Dashboard</h1>
          <p className="text-gray-600 mt-3">{allTasks.length} tasks assigned to you</p>
        </div>

        {/* Tasks Section */}
        {allTasks.length === 0 ? (
          <div className="bg-gray-50 p-12 rounded-lg border border-gray-200 text-center">
            <p className="text-gray-500 text-lg">No tasks assigned yet</p>
          </div>
        ) : (
        <>
          {/* Tasks Grouped by Project */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-12">
            {Object.entries(
              allTasks.reduce((acc, task) => {
                const projectId = task.project_id;
                const projectName = projects[projectId]?.name || `Project ${projectId}`;
                if (!acc[projectId]) {
                  acc[projectId] = { name: projectName, tasks: [] };
                }
                acc[projectId].tasks.push(task);
                return acc;
              }, {})
            ).map(([projectId, projectGroup]) => (
              <div key={projectId} className="border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">{projectGroup.name}</h3>
                <div className="space-y-3">
                  {projectGroup.tasks.map((task) => (
                    <div
                      key={task.id}
                      onClick={() => navigate(`/tasks/${task.id}`)}
                      className="p-4 border border-gray-100 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-all cursor-pointer"
                    >
                      <div className="flex justify-between items-start gap-2 mb-2">
                        <h4 className="font-medium text-gray-900 flex-1 text-sm">{task.title}</h4>
                        <span className={`inline-flex px-2 py-1 rounded text-xs font-semibold border ${
                          task.status === 'DONE' ? 'bg-green-50 text-green-700 border-green-200' :
                          task.status === 'IN_PROGRESS' ? 'bg-yellow-50 text-yellow-700 border-yellow-200' :
                          'bg-gray-50 text-gray-700 border-gray-200'
                        }`}>
                          {task.status}
                        </span>
                      </div>
                      <div className="flex gap-2 items-center text-xs text-gray-500">
                        <span className={`px-2 py-0.5 rounded ${
                          task.priority === 'HIGH' ? 'bg-red-100 text-red-700' :
                          task.priority === 'MEDIUM' ? 'bg-yellow-100 text-yellow-700' :
                          'bg-blue-100 text-blue-700'
                        }`}>
                          {task.priority}
                        </span>
                        {task.due_date && (
                          <span className="text-gray-400">
                            Due: {new Date(task.due_date).toLocaleDateString()}
                          </span>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>

          {/* Tasks Table (Alternative view) */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">All Tasks View</h3>
            <div className="overflow-x-auto">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-gray-50 border-b">
                    <SortHeader field="title" label="Title" />
                    <SortHeader field="project_id" label="Project" />
                    <SortHeader field="status" label="Status" />
                    <SortHeader field="priority" label="Priority" />
                    <SortHeader field="due_date" label="Due Date" />
                  </tr>
                </thead>
                <tbody>
                  {paginatedTasks.map((task) => (
                    <tr key={task.id} className="border-b hover:bg-gray-50 transition-colors">
                      <td 
                        className="px-6 py-4 font-medium text-blue-600 hover:text-blue-700 hover:underline cursor-pointer"
                        onClick={() => navigate(`/tasks/${task.id}`)}
                      >
                        {task.title}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        {projects[task.project_id]?.name || `Project ${task.project_id}`}
                      </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold border ${
                        task.status === 'DONE' ? 'bg-green-50 text-green-700 border-green-200' :
                        task.status === 'IN_PROGRESS' ? 'bg-yellow-50 text-yellow-700 border-yellow-200' :
                        'bg-gray-50 text-gray-700 border-gray-200'
                      }`}>
                        {task.status}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold border ${
                        task.priority === 'HIGH' ? 'bg-red-50 text-red-700 border-red-200' :
                        task.priority === 'MEDIUM' ? 'bg-yellow-50 text-yellow-700 border-yellow-200' :
                        'bg-blue-50 text-blue-700 border-blue-200'
                      }`}>
                        {task.priority}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {task.due_date ? new Date(task.due_date).toLocaleDateString() : '-'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            </div>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-between items-center mb-8">
              <div className="text-sm text-gray-600">
                Showing {(currentPage - 1) * itemsPerPage + 1} to{' '}
                {Math.min(currentPage * itemsPerPage, sortedTasks.length)} of{' '}
                {sortedTasks.length} tasks
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                  className="btn-sm bg-gray-100 text-gray-700 hover:bg-gray-200 border border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  ← Previous
                </button>
                <span className="px-4 py-2 text-sm text-gray-600">
                  Page {currentPage} of {totalPages}
                </span>
                <button
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                  className="btn-sm bg-gray-100 text-gray-700 hover:bg-gray-200 border border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next →
                </button>
              </div>
            </div>
          )}
        </>
      )}

      {/* Recent Comments Section */}
      <div className="mt-12">
        <h2 className="text-2xl font-bold mb-6">Recent Comments</h2>
        {recentComments.length === 0 ? (
          <div className="bg-gray-50 p-8 rounded-lg text-center">
            <p className="text-gray-500">No comments yet</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr className="bg-gray-50 border-b">
                  <th className="px-6 py-3 text-left">ID</th>
                  <th className="px-6 py-3 text-left">Content</th>
                  <th className="px-6 py-3 text-left">Task</th>
                  <th className="px-6 py-3 text-left">Created</th>
                </tr>
              </thead>
              <tbody>
                {recentComments.map((comment) => (
                  <tr key={comment.id} className="border-b hover:bg-gray-50">
                    <td className="px-6 py-4">{comment.id}</td>
                    <td className="px-6 py-4 max-w-md truncate">{comment.content}</td>
                    <td 
                      className="px-6 py-4 text-sm text-blue-600 hover:text-blue-700 hover:underline cursor-pointer"
                      onClick={() => navigate(`/tasks/${comment.task_id}`)}
                    >
                      Task #{comment.task_id}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {new Date(comment.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
      </div>
    </div>
  );
}

export default Dashboard;
