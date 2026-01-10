function ErrorMessage({ message }) {
  return (
    <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <p className="font-semibold text-red-700 text-sm">Error</p>
      <p className="text-sm mt-2 text-red-600">{message || 'An error occurred'}</p>
    </div>
  );
}

export default ErrorMessage;
