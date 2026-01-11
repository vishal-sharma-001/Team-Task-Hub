import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAsync } from '../hooks/useAsync';
import { taskAPI, projectAPI, commentAPI } from '../api/client';
import Loading from '../components/Loading';
import ErrorMessage from '../components/ErrorMessage';

function Dashboard() {
  const navigate = useNavigate();
  const [assignedTasks, setAssignedTasks] = useState([]);
  const [projects, setProjects] = useState({});
  const [recentComments, setRecentComments] = useState([]);
  const [tasks, setTasks] = useState({});

  const { execute: fetchData, status, error } = useAsync(
    async () => {
      try {
        const tasksResponse = await taskAPI.getAssignedToMe();
        const tasks = tasksResponse.data || tasksResponse || [];

        const projectsResponse = await projectAPI.getAll();
        const projectsList = projectsResponse.data || projectsResponse || [];

        const commentsResponse = await commentAPI.getRecent(1, 100);
        const allComments = commentsResponse.data || commentsResponse || [];

        const projectLookup = {};
        projectsList.forEach((p) => {
          projectLookup[p.id] = p;
        });

        const taskLookup = {};
        tasks.forEach((t) => {
          taskLookup[t.id] = t;
        });

        const assignedTaskIds = new Set(tasks.map(t => t.id));
        const filteredComments = allComments.filter(c => assignedTaskIds.has(c.task_id));

        setProjects(projectLookup);
        setTasks(taskLookup);
        setRecentComments(filteredComments);
        setAssignedTasks(Array.isArray(tasks) ? tasks : []);
        return { tasks, projects: projectsList, comments: filteredComments };
      } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
        throw err;
      }
    }
  );

  useEffect(() => {
    fetchData();
  }, []);

  const tasksByProject = assignedTasks.reduce((acc, task) => {
    const projectId = task.project_id;
    const projectName = projects[projectId]?.name || `Project ${projectId}`;
    if (!acc[projectId]) {
      acc[projectId] = { name: projectName, tasks: [] };
    }
    acc[projectId].tasks.push(task);
    return acc;
  }, {});

  if (status === 'pending') return <Loading />;

  const getStatusColor = (status) => {
    switch (status) {
      case 'DONE':
        return 'bg-green-50 text-green-700';
      case 'IN_PROGRESS':
        return 'bg-blue-50 text-blue-700';
      default:
        return 'bg-gray-100 text-gray-700';
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'HIGH':
        return 'bg-red-50 text-red-700';
      case 'MEDIUM':
        return 'bg-yellow-50 text-yellow-700';
      default:
        return 'bg-green-50 text-green-700';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container pt-12 pb-16 px-2 sm:px-4">
        {error && <ErrorMessage message={error.message} />}

        {/* Header Section */}
        <div className="mb-12">
          <div className="flex items-center justify-between mb-2">
            <h1 className="text-4xl font-bold text-gray-900">Dashboard</h1>
            <span className="inline-flex items-center px-4 py-2 rounded-full bg-blue-50 text-blue-700 text-sm font-semibold">
              {assignedTasks.length} {assignedTasks.length === 1 ? 'task' : 'tasks'}
            </span>
          </div>
          <p className="text-gray-600 text-base">Manage your assigned tasks and projects</p>
        </div>

        {/* Tasks and Comments Section */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Tasks Section - 3 columns on large screens */}
          <div className="lg:col-span-3">
            <h2 className="text-2xl font-bold text-gray-900 mb-8">Assigned Tasks</h2>
            
            {assignedTasks.length === 0 ? (
              <div className="bg-white p-12 text-center rounded-lg shadow-sm border border-gray-200">
                <p className="text-gray-500 text-base">No tasks assigned yet</p>
              </div>
            ) : (
              <div className="space-y-4">
                {Object.entries(tasksByProject).map(([projectId, projectGroup]) => (
                  <div key={projectId} className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                    <div className="bg-gradient-to-r from-blue-50 to-blue-100 px-6 py-3 border-b border-gray-200">
                      <h3 className="font-bold text-gray-900 text-sm">{projectGroup.name}</h3>
                    </div>
                    <div className="space-y-0">
                      {projectGroup.tasks.map((task) => (
                        <div
                          key={task.id}
                          onClick={() => navigate(`/tasks/${task.id}`)}
                          className="px-6 py-4 cursor-pointer hover:bg-blue-50 transition-colors duration-150 border-b border-gray-100 last:border-b-0 hover:border-l-4 hover:border-l-blue-500"
                        >
                          <div className="flex justify-between items-center gap-3 mb-2">
                            <h4 className="font-medium text-gray-900 flex-1 text-sm">{task.title}</h4>
                            <span className={`px-2.5 py-1 rounded text-xs font-semibold whitespace-nowrap ${getStatusColor(task.status)}`}>
                              {task.status === 'DONE' ? 'Done' :
                               task.status === 'IN_PROGRESS' ? 'In Progress' :
                               'Open'}
                            </span>
                          </div>
                          <div className="flex flex-wrap gap-3 items-center text-xs text-gray-600">
                            <span className={`px-2 py-0.5 rounded text-xs font-semibold ${getPriorityColor(task.priority)}`}>
                              {task.priority === 'HIGH' ? 'High' :
                               task.priority === 'MEDIUM' ? 'Medium' :
                               'Low'} priority
                            </span>
                            {task.due_date && (
                              <span>
                                📅 {new Date(task.due_date).toLocaleDateString()}
                              </span>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Recent Comments Section - 1 column on large screens */}
          <div className="lg:col-span-1">
            <h2 className="text-2xl font-bold text-gray-900 mb-8">Recent Comments</h2>
            {recentComments.length === 0 ? (
              <div className="bg-white p-8 text-center rounded-lg shadow-sm border border-gray-200">
                <p className="text-gray-500 text-base">No comments yet</p>
              </div>
            ) : (
              <div className="space-y-3">
                {recentComments.map((comment) => (
                  <div key={comment.id} className="bg-white p-4 rounded-lg border border-gray-200 hover:shadow-md hover:border-blue-300 transition-all duration-150">
                    <div className="flex items-center justify-between gap-2 mb-2">
                      <p className="text-xs font-bold text-gray-900">{comment.author_name || 'Unknown'}</p>
                      <span className="text-xs text-gray-500 whitespace-nowrap">
                        {new Date(comment.created_at).toLocaleDateString()}
                      </span>
                    </div>
                    <p className="text-xs text-gray-700 mb-3 line-clamp-2">{comment.content}</p>
                    <button 
                      onClick={() => navigate(`/tasks/${comment.task_id}`)}
                      className="text-xs text-blue-600 hover:text-blue-700 hover:underline font-semibold"
                    >
                      {tasks[comment.task_id]?.title || `Task #${comment.task_id}`} →
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
