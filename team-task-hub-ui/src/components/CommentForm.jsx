import { useForm } from '../hooks/useAsync';
import { commentAPI } from '../api/client';

function CommentForm({ taskId, onCommentAdded }) {
  const { values, errors, touched, handleChange, handleBlur, handleSubmit } =
    useForm(
      {
        content: '',
      },
      {
        content: (value) => {
          if (!value) return 'Comment cannot be empty';
          if (value.length < 1) return 'Comment is required';
          return '';
        },
      },
      async (formData) => {
        try {
          const newComment = await commentAPI.create(taskId, formData);
          onCommentAdded(newComment);
          // Reset form
          document.querySelector('textarea[name="content"]').value = '';
        } catch (err) {
          console.error('Failed to add comment:', err);
        }
      }
    );

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <textarea
        name="content"
        value={values.content}
        onChange={handleChange}
        onBlur={handleBlur}
        placeholder="Add a comment..."
        className={`input w-full resize-none min-h-24 ${
          touched.content && errors.content ? 'border-red-500 focus:border-red-500 focus:ring-red-200' : ''
        }`}
        rows="4"
      />
      {touched.content && errors.content && (
        <p className="text-red-600 text-sm">{errors.content}</p>
      )}
      <button type="submit" className="btn-primary text-sm">
        Post Comment
      </button>
    </form>
  );
}

export default CommentForm;
