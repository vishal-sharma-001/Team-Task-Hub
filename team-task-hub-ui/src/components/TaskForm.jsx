import { useState, useEffect } from 'react';
import { useForm } from '../hooks/useAsync';
import { userAPI } from '../api/client';

const TASK_STATUSES = ['OPEN', 'IN_PROGRESS', 'DONE'];
const TASK_PRIORITIES = ['LOW', 'MEDIUM', 'HIGH'];

function TaskForm({ task, onSubmit, onCancel }) {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const data = await userAPI.list();
        setUsers(data || []);
      } catch (err) {
        console.error('Failed to fetch users:', err);
      }
    };
    fetchUsers();
  }, []);

  const { values, errors, touched, handleChange, handleBlur, handleSubmit } =
    useForm(
      {
        title: task?.title || '',
        description: task?.description || '',
        status: task?.status || 'OPEN',
        priority: task?.priority || 'MEDIUM',
        assignee_id: task?.assignee_id ? String(task.assignee_id) : '',
        due_date: task?.due_date || '',
      },
      {
        title: (value) => {
          if (!value) return 'Task title is required';
          if (value.length < 3) return 'Task title must be at least 3 characters';
          return '';
        },
      },
      onSubmit
    );

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-2">
          <label htmlFor="title" className="block text-sm font-medium text-gray-700">
            Task Title *
          </label>
          <input
            id="title"
            type="text"
            name="title"
            value={values.title}
            onChange={handleChange}
            onBlur={handleBlur}
            className={`input w-full ${
              touched.title && errors.title ? 'border-red-500 focus:border-red-500 focus:ring-red-200' : ''
            }`}
            placeholder="Enter task title"
          />
          {touched.title && errors.title && (
            <p className="text-red-600 text-sm mt-1">{errors.title}</p>
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
            placeholder="Enter task description"
          />
        </div>

        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-2">
            <label htmlFor="status" className="block text-sm font-medium text-gray-700">
              Status
            </label>
            <select
              id="status"
              name="status"
              value={values.status}
              onChange={handleChange}
              className="input w-full"
            >
              {TASK_STATUSES.map((status) => (
                <option key={status} value={status}>
                  {status}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-2">
            <label htmlFor="priority" className="block text-sm font-medium text-gray-700">
              Priority
            </label>
            <select
              id="priority"
              name="priority"
              value={values.priority}
              onChange={handleChange}
              className="input w-full"
            >
              {TASK_PRIORITIES.map((priority) => (
                <option key={priority} value={priority}>
                  {priority}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-2">
            <label htmlFor="assignee_id" className="block text-sm font-medium text-gray-700">
              Assignee
            </label>
            <select
              id="assignee_id"
              name="assignee_id"
              value={values.assignee_id}
              onChange={handleChange}
              onBlur={handleBlur}
              className="input w-full"
            >
              <option value="">Select a user...</option>
              {users.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.email}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-2">
            <label htmlFor="due_date" className="block text-sm font-medium text-gray-700">
              Due Date
            </label>
            <input
              id="due_date"
              type="date"
              name="due_date"
              value={values.due_date}
              onChange={handleChange}
              onBlur={handleBlur}
              className="input w-full"
              min={new Date().toISOString().split('T')[0]}
            />
          </div>
        </div>

        <div className="flex gap-4 pt-6 border-t border-gray-200">
          <button type="submit" className="btn-primary">
            {task ? 'Update Task' : 'Create Task'}
          </button>
          <button type="button" onClick={onCancel} className="btn-secondary">
            Cancel
          </button>
        </div>
      </form>
  );
}

export default TaskForm;
