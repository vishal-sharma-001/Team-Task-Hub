import { commentAPI } from '../api/client';
import { useState } from 'react';

function CommentForm({ taskId, onCommentAdded }) {
  const [content, setContent] = useState('');
  const [errors, setErrors] = useState('');
  const [touched, setTouched] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateContent = (value) => {
    if (!value) return 'Comment cannot be empty';
    return '';
  };

  const handleChange = (e) => {
    setContent(e.target.value);
    if (touched) {
      setErrors(validateContent(e.target.value));
    }
  };

  const handleBlur = () => {
    setTouched(true);
    setErrors(validateContent(content));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const error = validateContent(content);
    if (error) {
      setErrors(error);
      setTouched(true);
      return;
    }

    setIsSubmitting(true);
    try {
      const newComment = await commentAPI.create(taskId, { content });
      onCommentAdded(newComment);
      setContent('');
      setErrors('');
      setTouched(false);
    } catch (err) {
      console.error('Failed to post comment:', err);
      setErrors('Failed to post comment. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={content}
        onChange={handleChange}
        onBlur={handleBlur}
        placeholder="Add a comment..."
        disabled={isSubmitting}
        className={`input w-full resize-none min-h-24 text-sm ${
          touched && errors ? 'border-red-500 focus:border-red-500' : ''
        }`}
        rows="4"
      />
      {touched && errors && (
        <p className="text-red-600 text-xs font-medium">{errors}</p>
      )}
      <button 
        type="submit" 
        disabled={isSubmitting}
        className="btn-primary text-sm"
      >
        {isSubmitting ? 'Posting...' : 'Post Comment'}
      </button>
    </form>
  );
}

export default CommentForm;
