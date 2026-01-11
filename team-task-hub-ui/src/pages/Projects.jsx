import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAsync } from '../hooks/useAsync';
import { projectAPI } from '../api/client';
import Loading from '../components/Loading';
import ErrorMessage from '../components/ErrorMessage';
import ProjectForm from '../components/ProjectForm';
import Modal from '../components/Modal';
import ConfirmDialog from '../components/ConfirmDialog';

function Projects() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [editingProject, setEditingProject] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('id');
  const [sortOrder, setSortOrder] = useState('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage] = useState(10);
  const [deleteConfirm, setDeleteConfirm] = useState({ isOpen: false, projectId: null, projectName: '' });

  const { execute: fetchProjects, status, error } = useAsync(
    async () => {
      const data = await projectAPI.getAll();
      return data;
    }
  );

  useEffect(() => {
    fetchProjects().then((data) => {
      if (data && Array.isArray(data)) {
        setProjects(data);
      }
    }).catch((err) => {
      console.error('Failed to fetch projects:', err);
    });
  }, []);

  const handleCreate = async (formData) => {
    try {
      const newProject = await projectAPI.create(formData);
      setProjects([...projects, newProject]);
      setShowForm(false);
    } catch (err) {
      console.error('Failed to create project:', err);
    }
  };

  const handleUpdate = async (projectId, formData) => {
    try {
      const updated = await projectAPI.update(projectId, formData);
      setProjects(projects.map((p) => (p.id === projectId ? updated : p)));
      setEditingProject(null);
    } catch (err) {
      console.error('Failed to update project:', err);
    }
  };

  const handleDelete = async (projectId) => {
    const project = projects.find(p => p.id === projectId);
    setDeleteConfirm({
      isOpen: true,
      projectId: projectId,
      projectName: project?.name || 'Project'
    });
  };

  const handleConfirmDelete = async () => {
    try {
      await projectAPI.delete(deleteConfirm.projectId);
      setProjects(projects.filter((p) => p.id !== deleteConfirm.projectId));
      setDeleteConfirm({ isOpen: false, projectId: null, projectName: '' });
    } catch (err) {
      console.error('Failed to delete project:', err);
    }
  };

  const handleCancelDelete = () => {
    setDeleteConfirm({ isOpen: false, projectId: null, projectName: '' });
  };

  // Filter projects
  const filteredProjects = projects.filter((p) =>
    p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (p.description && p.description.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  // Sort projects
  const sortedProjects = [...filteredProjects].sort((a, b) => {
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

  // Paginate
  const totalPages = Math.ceil(sortedProjects.length / itemsPerPage);
  const paginatedProjects = sortedProjects.slice(
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
      className="px-6 py-2 text-left text-xs font-semibold text-gray-600 uppercase tracking-wide cursor-pointer hover:bg-gray-100"
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
    <div className="min-h-screen bg-gray-50">
      <div className="container pt-12 pb-16 px-2 sm:px-4">
        {/* Header Section */}
        <div className="mb-12">
          <div className="flex items-center justify-between gap-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Projects</h1>
              <p className="text-gray-600 text-sm">Manage all your projects and tasks</p>
            </div>
            <button
              onClick={() => setShowForm(!showForm)}
              className="btn-primary flex-shrink-0"
            >
              {showForm ? 'Cancel' : '+ New Project'}
            </button>
          </div>
        </div>

        {error && <ErrorMessage message={error.message} />}

        <Modal
          isOpen={showForm}
          onClose={() => setShowForm(false)}
          title="New Project"
        >
          <ProjectForm
            onSubmit={handleCreate}
            onCancel={() => setShowForm(false)}
          />
        </Modal>

        <Modal
          isOpen={!!editingProject}
          onClose={() => setEditingProject(null)}
          title="Edit Project"
        >
          {editingProject && (
            <ProjectForm
              project={editingProject}
              onSubmit={(data) => handleUpdate(editingProject.id, data)}
              onCancel={() => setEditingProject(null)}
            />
          )}
        </Modal>

        {/* Search Bar */}
        {projects.length > 0 && (
          <div className="mb-8">
            <input
              type="text"
              placeholder="Search projects..."
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value);
                setCurrentPage(1);
              }}
              className="input w-full max-w-md"
            />
          </div>
        )}

        {/* Projects Table */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="bg-gray-100 border-b border-gray-200">
                <SortHeader field="id" label="ID" />
                <SortHeader field="name" label="Project Name" />
                <SortHeader field="description" label="Description" />
                <SortHeader field="created_at" label="Created" />
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {paginatedProjects.map((project) => (
                <tr 
                  key={project.id} 
                  onClick={() => navigate(`/projects/${project.id}/tasks`)}
                  className="cursor-pointer hover:bg-gray-50 transition-colors"
                >
                  <td className="px-6 py-3 text-sm text-gray-900">{project.id}</td>
                  <td 
                    className="px-6 py-3 text-sm font-medium text-gray-900"
                  >
                    {project.name}
                  </td>
                  <td className="px-6 py-3 text-sm text-gray-600">{project.description || '-'}</td>
                  <td className="px-6 py-3 text-sm text-gray-500">
                    {new Date(project.created_at).toLocaleDateString()}
                  </td>
                </tr>
              ))}
              {paginatedProjects.length === 0 && (
                <tr>
                  <td colSpan="4" className="px-6 py-8 text-center text-sm text-gray-500">
                    No projects found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {projects.length > 0 && totalPages > 1 && (
          <div className="mt-8 flex justify-between items-center">
            <div className="text-sm text-gray-600">
              Showing {(currentPage - 1) * itemsPerPage + 1} to{' '}
              {Math.min(currentPage * itemsPerPage, sortedProjects.length)} of{' '}
              {sortedProjects.length} projects
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                disabled={currentPage === 1}
                className="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                ← Previous
              </button>
              <span className="px-4 py-2 text-sm text-gray-600">
                Page {currentPage} of {totalPages}
              </span>
              <button
                onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                disabled={currentPage === totalPages}
                className="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next →
              </button>
            </div>
          </div>
        )}

        <ConfirmDialog
          isOpen={deleteConfirm.isOpen}
          title="Delete Project"
          message={`Are you sure you want to delete "${deleteConfirm.projectName}"? All tasks in this project will also be deleted.`}
          confirmText="Delete"
          cancelText="Cancel"
          onConfirm={handleConfirmDelete}
          onCancel={() => setDeleteConfirm({ isOpen: false, projectId: null, projectName: '' })}
          isDanger={true}
        />
      </div>
    </div>
  );}

export default Projects;