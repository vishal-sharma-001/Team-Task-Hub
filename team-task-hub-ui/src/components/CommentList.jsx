function CommentList({ comments, onDelete }) {
  if (!comments || comments.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <p>No comments yet</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {comments.map((comment) => {
        const createdAt = new Date(comment.created_at).toLocaleString();
        return (
          <div key={comment.id} className="bg-gray-50 rounded-lg p-4 border border-gray-200">
            <div className="flex justify-between items-start mb-2">
              <div>
                <p className="font-semibold text-sm">{comment.author_name}</p>
                <p className="text-xs text-gray-500">{createdAt}</p>
              </div>
              {onDelete && (
                <button
                  onClick={() => onDelete(comment.id)}
                  className="text-red-600 hover:text-red-800 text-sm"
                >
                  Delete
                </button>
              )}
            </div>
            <p className="text-gray-700 text-sm">{comment.content}</p>
          </div>
        );
      })}
    </div>
  );
}

export default CommentList;
