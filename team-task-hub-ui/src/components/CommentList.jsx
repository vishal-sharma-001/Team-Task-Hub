import { useState } from 'react';
import { commentAPI } from '../api/client';

function CommentList({ comments, onDelete, onEdit, onUpdate, currentUserId }) {
  const [editingId, setEditingId] = useState(null);
  const [editContent, setEditContent] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  const handleEditClick = (comment) => {
    setEditingId(comment.id);
    setEditContent(comment.content);
    if (onEdit) onEdit(comment);
  };

  const handleSaveEdit = async (commentId) => {
    if (!editContent.trim()) {
      alert('Comment cannot be empty');
      return;
    }

    setIsSaving(true);
    try {
      const updated = await commentAPI.update(commentId, { content: editContent });
      setEditingId(null);
      setEditContent('');
      // Notify parent to update the comment in state
      if (onUpdate) onUpdate(updated || { id: commentId, content: editContent });
    } catch (err) {
      console.error('Failed to update comment:', err);
      alert('Failed to update comment');
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditContent('');
  };

  if (!comments || comments.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <p>No comments yet</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {comments.map((comment) => {
        const createdAt = new Date(comment.created_at).toLocaleString();
        const isOwner = currentUserId === comment.user_id;
        const isEditing = editingId === comment.id;

        return (
          <div key={comment.id} className="border border-gray-200 rounded p-4">
            <div className="flex justify-between items-start mb-2">
              <div>
                <p className="font-medium text-gray-900 text-sm">{comment.author_name || 'Unknown'}</p>
                <p className="text-xs text-gray-500">{createdAt}</p>
              </div>
              {isOwner && !isEditing && (
                <div className="flex gap-2">
                  {onEdit && (
                    <button
                      onClick={() => handleEditClick(comment)}
                      className="text-gray-400 hover:text-gray-600 p-1"
                      title="Edit"
                    >
                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                      </svg>
                    </button>
                  )}
                  {onDelete && (
                    <button
                      onClick={() => onDelete(comment.id)}
                      className="text-gray-400 hover:text-red-600 p-1"
                      title="Delete"
                    >
                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                      </svg>
                    </button>
                  )}
                </div>
              )}
            </div>
            
            {isEditing ? (
              <div className="space-y-2">
                <textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  className="input w-full resize-none min-h-20 text-sm"
                  rows="3"
                />
                <div className="flex gap-2">
                  <button
                    onClick={() => handleSaveEdit(comment.id)}
                    disabled={isSaving}
                    className="text-gray-400 hover:text-green-600 p-1 disabled:opacity-50"
                    title="Save"
                  >
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </button>
                  <button
                    onClick={handleCancelEdit}
                    disabled={isSaving}
                    className="text-gray-400 hover:text-gray-600 p-1 disabled:opacity-50"
                    title="Cancel"
                  >
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                    </svg>
                  </button>
                </div>
              </div>
            ) : (
              <p className="text-gray-700 text-sm">{comment.content}</p>
            )}
          </div>
        );
      })}
    </div>
  );
}

export default CommentList;
