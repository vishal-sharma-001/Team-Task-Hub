import { useForm } from '../hooks/useAsync';

function ProjectForm({ project, onSubmit, onCancel }) {
  const { values, errors, touched, handleChange, handleBlur, handleSubmit } =
    useForm(
      {
        name: project?.name || '',
        description: project?.description || '',
      },
      {
        name: (value) => {
          if (!value) return 'Project name is required';
          if (value.length < 3) return 'Project name must be at least 3 characters';
          return '';
        },
      },
      onSubmit
    );

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-2">
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">
            Project Name *
          </label>
          <input
            id="name"
            type="text"
            name="name"
            value={values.name}
            onChange={handleChange}
            onBlur={handleBlur}
            className={`input w-full ${
              touched.name && errors.name ? 'border-red-500 focus:border-red-500 focus:ring-red-200' : ''
            }`}
            placeholder="Enter project name"
          />
          {touched.name && errors.name && (
            <p className="text-red-600 text-sm mt-1">{errors.name}</p>
          )}
        </div>

        <div className="space-y-2">
          <label htmlFor="description" className="block text-sm font-medium text-gray-700">
            Description
          </label>
          <textarea
            id="description"
            name="description"
            value={values.description}
            onChange={handleChange}
            onBlur={handleBlur}
            className="input w-full resize-none min-h-24"
            rows="4"
            placeholder="Enter project description"
          />
        </div>

        <div className="flex gap-4 pt-6 border-t border-gray-200">
          <button type="submit" className="btn-primary">
            {project ? 'Update Project' : 'Create Project'}
          </button>
          <button type="button" onClick={onCancel} className="btn-secondary">
            Cancel
          </button>
        </div>
      </form>
  );
}

export default ProjectForm;
