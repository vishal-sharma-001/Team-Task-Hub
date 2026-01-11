import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAsync } from '../hooks/useAsync';
import { taskAPI, commentAPI, userAPI } from '../api/client';
import Loading from '../components/Loading';
import ErrorMessage from '../components/ErrorMessage';
import CommentForm from '../components/CommentForm';
import CommentList from '../components/CommentList';
import TaskForm from '../components/TaskForm';
import ConfirmDialog from '../components/ConfirmDialog';
import Modal from '../components/Modal';
import UserProfileModal from '../components/UserProfileModal';

const TASK_STATUSES = ['OPEN', 'IN_PROGRESS', 'DONE'];
const TASK_PRIORITIES = ['LOW', 'MEDIUM', 'HIGH'];

function TaskDetail() {
  const { taskId } = useParams();
  const navigate = useNavigate();
  const [task, setTask] = useState(null);
  const [comments, setComments] = useState([]);
  const [users, setUsers] = useState({});
  const [showEditForm, setShowEditForm] = useState(false);
  const [showUserProfile, setShowUserProfile] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState({ isOpen: false, isComment: false, commentId: null, commentText: '' });
  const [editingField, setEditingField] = useState(null);
  const [editValues, setEditValues] = useState({});
  const [isSaving, setIsSaving] = useState(false);
  const [currentUser, setCurrentUser] = useState(null);

  const { execute: fetchTask, status: taskStatus, error: taskError } = useAsync(
    async () => {
      const data = await taskAPI.getByID(taskId);
      return data;
    }
  );

  const { execute: fetchComments, status: commentsStatus } = useAsync(
    async () => {
      const data = await commentAPI.getByTask(taskId);
      return Array.isArray(data) ? data : data.data || [];
    }
  );

  const { execute: fetchUsers } = useAsync(
    async () => {
      const data = await userAPI.list();
      const lookup = {};
      (Array.isArray(data) ? data : data.data || []).forEach(u => {
        lookup[u.id] = u;
      });
      return lookup;
    }
  );

  useEffect(() => {
    fetchTask().then(setTask);
    fetchComments().then(setComments);
    fetchUsers().then(setUsers);
    
    // Get current user from localStorage
    try {
      const token = localStorage.getItem('authToken');
      if (token) {
        const payload = JSON.parse(atob(token.split('.')[1]));
        setCurrentUser({ id: payload.user_id, email: payload.email });
      }
    } catch (err) {
      console.error('Failed to parse current user:', err);
    }
  }, [taskId]);

  const handleUpdateTask = async (formData) => {
    try {
      const updated = await taskAPI.update(taskId, formData);
      setTask(updated);
      setShowEditForm(false);
    } catch (err) {
      console.error('Failed to update task:', err);
    }
  };

  const handleDeleteTask = async () => {
    setDeleteConfirm({
      isOpen: true,
      isComment: false,
      commentId: null,
      commentText: task.title
    });
  };

  const handleConfirmDeleteTask = async () => {
    try {
      await taskAPI.delete(taskId);
      setDeleteConfirm({ isOpen: false, isComment: false, commentId: null, commentText: '' });
      navigate(`/projects/${task.project_id}/tasks`);
    } catch (err) {
      console.error('Failed to delete task:', err);
    }
  };

  const handleCommentAdded = (newComment) => {
    setComments([...comments, newComment]);
  };

  const handleCommentUpdated = (updatedComment) => {
    setComments(comments.map(c => c.id === updatedComment.id ? updatedComment : c));
  };

  const handleCommentDeleted = async (commentId) => {
    const comment = comments.find(c => c.id === commentId);
    setDeleteConfirm({
      isOpen: true,
      isComment: true,
      commentId: commentId,
      commentText: comment?.text?.substring(0, 50) || 'Comment'
    });
  };

  const handleEditComment = (comment) => {
    // CommentList handles the edit action internally
    // This callback is just for notification purposes
  };

  const handleConfirmDeleteComment = async () => {
    try {
      await commentAPI.delete(deleteConfirm.commentId);
      setComments(comments.filter(c => c.id !== deleteConfirm.commentId));
      setDeleteConfirm({ isOpen: false, isComment: false, commentId: null, commentText: '' });
    } catch (err) {
      console.error('Failed to delete comment:', err);
    }
  };

  const handleInlineEdit = async (field, value = null) => {
    try {
      const newValue = value !== null ? value : (editValues[field] || task[field]);
      
      // Validate required fields
      if (field === 'title' && (!newValue || !newValue.trim())) {
        console.error('Title cannot be empty');
        setEditingField(null);
        setEditValues({});
        return;
      }
      
      setIsSaving(true);
      
      // Use dedicated PATCH endpoints for specific fields
      if (field === 'status') {
        await taskAPI.updateStatus(taskId, newValue);
      } else if (field === 'priority') {
        await taskAPI.updatePriority(taskId, newValue);
      } else if (field === 'assignee_id') {
        const assigneeId = newValue || null; // keep UUID string
        await taskAPI.updateAssignee(taskId, assigneeId);
      } else {
        // For other fields (title, description, due_date), use PUT
        await taskAPI.update(taskId, { [field]: newValue });
      }
      
      setTask({ ...task, [field]: newValue });
      setEditingField(null);
      setEditValues({});
    } catch (err) {
      console.error('Failed to update field:', err);
    } finally {
      setIsSaving(false);
    }
  };

  if (taskStatus === 'pending') return <Loading />;
  if (taskError) return <ErrorMessage message={taskError.message} />;
  if (!task) return <ErrorMessage message="Task not found" />;

  const assigneeUser = task.assignee_id ? users[task.assignee_id] : null;
  const createdByUser = task.created_by_id ? users[task.created_by_id] : null;
  const assignedByUser = task.assigned_by_id ? users[task.assigned_by_id] : null;

  return (
    <div className="min-h-screen bg-white">
      <div className="container pt-8 pb-12 px-2 sm:px-4">
        {/* Header with Back Button */}
        <div className="mb-8">
          <button
            onClick={() => navigate(-1)}
            className="text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors mb-6"
          >
            ← Back to project
          </button>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Left Column - Task + Comments */}
          <div className="lg:col-span-3 space-y-6">
            {/* Task Header Card */}
            <div className="border border-gray-200 rounded-lg p-8 bg-white">
              {/* Title Row */}
              <div className="mb-8">
                {editingField === 'title' ? (
                  <input
                    type="text"
                    value={editValues.title !== undefined ? editValues.title : task.title}
                    onChange={(e) => setEditValues({ ...editValues, title: e.target.value })}
                    onBlur={() => handleInlineEdit('title')}
                    onKeyDown={(e) => e.key === 'Enter' && handleInlineEdit('title')}
                    disabled={isSaving}
                    className="w-full text-4xl font-bold text-gray-900 outline-none disabled:opacity-50 bg-gray-50 border border-gray-300 px-3 py-2 rounded"
                    placeholder="Enter task title"
                  />
                ) : (
                  <h1
                    onClick={() => {
                      setEditingField('title');
                      setEditValues({ title: task.title });
                    }}
                    className="text-4xl font-bold text-gray-900 cursor-text"
                    title="Click to edit"
                  >
                    {task.title}
                  </h1>
                )}
                <p className="text-gray-600 text-sm mt-2">Task #{task.id}</p>
              </div>

              {/* Description */}
              <div className="mb-8 pb-8 border-t border-gray-200 pt-8">
                <h3 className="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">Description</h3>
                {editingField === 'description' ? (
                  <textarea
                    value={editValues.description !== undefined ? editValues.description : (task.description || '')}
                    onChange={(e) => setEditValues({ ...editValues, description: e.target.value })}
                    onBlur={() => handleInlineEdit('description')}
                    disabled={isSaving}
                    className="w-full min-h-[120px] p-3 outline-none font-base text-gray-700 disabled:opacity-50 bg-gray-50 border border-gray-300 rounded"
                    placeholder="Click to add description"
                  />
                ) : (
                  <div
                    onClick={() => {
                      setEditingField('description');
                      setEditValues({ description: task.description || '' });
                    }}
                    className={`text-base leading-relaxed whitespace-pre-wrap cursor-text min-h-[60px] p-3 rounded border border-gray-200 bg-gray-50 ${
                      task.description ? 'text-gray-700' : 'text-gray-400 italic'
                    }`}
                    title="Click to edit"
                  >
                    {task.description || 'Click to add description'}
                  </div>
                )}
              </div>

            </div>

            {/* Comments Section */}
            <div className="border border-gray-200 rounded-lg p-8 bg-white">
              <h2 className="text-xl font-bold text-gray-900 mb-6">Comments</h2>
              
              <CommentForm 
                taskId={taskId} 
                onCommentAdded={handleCommentAdded}
              />

              {comments.length > 0 && (
                <div className="mt-8">
                  <CommentList
                    comments={comments}
                    users={users}
                    onDelete={handleCommentDeleted}
                    onEdit={handleEditComment}
                    onUpdate={handleCommentUpdated}
                    currentUserId={currentUser?.id}
                  />
                </div>
              )}
            </div>
          </div>

            {/* Right Column - Compact Details */}
            <div className="h-fit sticky top-24">
              {/* Status */}
              <div className="mb-4">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1">Status</p>
                {editingField === 'status' ? (
                  <select
                    value={task.status}
                    onChange={(e) => handleInlineEdit('status', e.target.value)}
                    onBlur={() => setEditingField(null)}
                    disabled={isSaving}
                    className="w-full px-3 py-2 rounded text-sm font-semibold disabled:opacity-50 bg-gray-50 border border-gray-300"
                  >
                    {TASK_STATUSES.map(s => (
                      <option key={s} value={s}>{s}</option>
                    ))}
                  </select>
                ) : (
                  <button
                    onClick={() => setEditingField('status')}
                    disabled={isSaving}
                    className="w-full px-3 py-2 rounded text-sm font-semibold text-left disabled:opacity-50 cursor-pointer transition-colors bg-gray-100 text-gray-800 hover:bg-gray-150"
                    title="Click to edit"
                  >
                    {task.status}
                  </button>
                )}
              </div>

              {/* Priority */}
              <div className="mb-4">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1">Priority</p>
                {editingField === 'priority' ? (
                  <select
                    autoFocus
                    value={task.priority}
                    onChange={(e) => handleInlineEdit('priority', e.target.value)}
                    onBlur={() => setEditingField(null)}
                    disabled={isSaving}
                    className="w-full px-3 py-2 rounded text-sm font-semibold disabled:opacity-50 bg-gray-50 border border-gray-300"
                  >
                    {TASK_PRIORITIES.map(p => (
                      <option key={p} value={p}>{p}</option>
                    ))}
                  </select>
                ) : (
                  <button
                    onClick={() => {
                      setEditingField('priority');
                      setEditValues({ priority: task.priority });
                    }}
                    disabled={isSaving}
                    className="w-full px-3 py-2 rounded text-sm font-semibold text-left disabled:opacity-50 cursor-pointer transition-colors bg-gray-100 text-gray-800 hover:bg-gray-150"
                    title="Click to edit"
                  >
                    {task.priority}
                  </button>
                )}
              </div>

              {/* Due Date */}
              <div className="mb-4 pb-4 border-b border-gray-200">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1">Due Date</p>
                {editingField === 'due_date' ? (
                  <input
                    autoFocus
                    type="date"
                    value={editValues.due_date !== undefined ? editValues.due_date : (task?.due_date ? (typeof task.due_date === 'string' ? task.due_date.split('T')[0] : task.due_date) : '')}
                    min={new Date().toISOString().split('T')[0]}
                    onChange={(e) => setEditValues({ ...editValues, due_date: e.target.value })}
                    onBlur={() => {
                      if (editValues.due_date) {
                        handleInlineEdit('due_date', editValues.due_date);
                      } else {
                        setEditingField(null);
                      }
                    }}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' && editValues.due_date) {
                        handleInlineEdit('due_date', editValues.due_date);
                      }
                    }}
                    disabled={isSaving}
                    className="px-3 py-2 rounded disabled:opacity-50 w-full bg-gray-50 border border-gray-300"
                  />
                ) : (
                  <p
                    onClick={() => {
                      if (!isSaving) {
                        setEditingField('due_date');
                        const dateValue = task?.due_date ? (typeof task.due_date === 'string' ? task.due_date.split('T')[0] : task.due_date) : '';
                        setEditValues({ due_date: dateValue });
                      }
                    }}
                    className="text-sm text-gray-800 font-medium cursor-text disabled:opacity-50"
                    title="Click to edit"
                  >
                    {task?.due_date ? (typeof task.due_date === 'string' ? new Date(task.due_date).toLocaleDateString() : 'Invalid date') : '—'}
                  </p>
                )}
              </div>

              {/* Created/Updated (Fixed - Not Editable) */}
              <div className="mb-4 pb-4 border-b border-gray-200 text-sm">
                <div className="mb-2">
                  <p className="text-xs text-gray-600 mb-0.5">Created</p>
                  <p className="text-gray-800 font-medium">{new Date(task.created_at).toLocaleDateString()}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-600 mb-0.5">Updated</p>
                  <p className="text-gray-800 font-medium">{new Date(task.updated_at).toLocaleDateString()}</p>
                </div>
              </div>

              {/* Assignee */}
              <div className="mb-4 pb-4 border-b border-gray-200">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1">Assignee</p>
                {editingField === 'assignee_id' ? (
                  <select
                    autoFocus
                    value={task.assignee_id || ''}
                    onChange={(e) => handleInlineEdit('assignee_id', e.target.value || null)}
                    onBlur={() => setEditingField(null)}
                    disabled={isSaving}
                    className="w-full px-3 py-2 rounded text-sm border-2 border-slate-300 disabled:opacity-50 bg-slate-50"
                  >
                    <option value="">Unassigned</option>
                    {Object.values(users).map(u => (
                      <option key={u.id} value={u.id}>{u.email}</option>
                    ))}
                  </select>
                ) : (
                  <button
                    onClick={() => setEditingField('assignee_id')}
                    disabled={isSaving}
                    className="w-full text-left hover:bg-slate-100 transition-colors rounded px-2 py-2 -mx-0.5 disabled:opacity-50 cursor-pointer"
                    title="Click to edit"
                  >
                    {assigneeUser ? (
                      <div className="flex items-center gap-2">
                        <div className="w-7 h-7 rounded-full bg-gradient-to-br from-blue-400 to-indigo-600 flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
                          {assigneeUser.email.charAt(0).toUpperCase()}
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-gray-900 truncate">{assigneeUser.email.split('@')[0]}</p>
                        </div>
                      </div>
                    ) : (
                      <span className="text-slate-500">Unassigned</span>
                    )}
                  </button>
                )}
              </div>

              {/* Created By */}
              <div className="mb-4 pb-4 border-b border-gray-200">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-2">Created By</p>
                {createdByUser ? (
                  <div className="flex items-center gap-2 px-2 py-2 rounded bg-purple-50">
                    <div className="w-7 h-7 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
                      {createdByUser.email.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">{createdByUser.email.split('@')[0]}</p>
                      <p className="text-xs text-gray-600 truncate">{createdByUser.email.split('@')[1]}</p>
                    </div>
                  </div>
                ) : (
                  <span className="text-slate-500">Unknown</span>
                )}
              </div>

              {/* Assigned By */}
              <div className="mb-4 pb-4 border-b border-gray-200">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-2">Assigned By</p>
                {assignedByUser ? (
                  <div className="flex items-center gap-2 px-2 py-2 rounded bg-orange-50">
                    <div className="w-7 h-7 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
                      {assignedByUser.email.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">{assignedByUser.email.split('@')[0]}</p>
                      <p className="text-xs text-gray-600 truncate">{assignedByUser.email.split('@')[1]}</p>
                    </div>
                  </div>
                ) : (
                  <span className="text-slate-500">Not assigned</span>
                )}
              </div>

              {/* Delete Button */}
              <button
                onClick={handleDeleteTask}
                className="btn-danger w-full text-sm"
              >
                Delete Task
              </button>
            </div>

            {/* Edit Modal */}
            <Modal
              isOpen={showEditForm}
              onClose={() => setShowEditForm(false)}
              title="Edit Task"
            >
              <TaskForm
                task={task}
                onSubmit={handleUpdateTask}
                onCancel={() => setShowEditForm(false)}
              />
            </Modal>

            {/* Comments Section - removed for debugging */}

        </div>

        <ConfirmDialog
          isOpen={deleteConfirm.isOpen}
          title={deleteConfirm.isComment ? "Delete Comment" : "Delete Task"}
          message={
            deleteConfirm.isComment 
              ? `Are you sure you want to delete this comment? "${deleteConfirm.commentText}..."`
              : `Are you sure you want to delete "${deleteConfirm.commentText}"? This action cannot be undone.`
          }
          confirmText="Delete"
          cancelText="Cancel"
          onConfirm={deleteConfirm.isComment ? handleConfirmDeleteComment : handleConfirmDeleteTask}
          onCancel={() => setDeleteConfirm({ isOpen: false, isComment: false, commentId: null, commentText: '' })}
          isDanger={true}
        />

        {/* User Profile Modal */}
        {assigneeUser && (
          <UserProfileModal
            isOpen={showUserProfile}
            onClose={() => setShowUserProfile(false)}
            user={assigneeUser}
            onProfileUpdate={(updatedUser) => {
              setUsers({ ...users, [updatedUser.id]: updatedUser });
            }}
          />
        )}
      </div>
    </div>
  );
}

export default TaskDetail;
