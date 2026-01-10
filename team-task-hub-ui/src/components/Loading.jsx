function Loading() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-b from-white to-gray-50">
      <div className="text-center">
        <div className="inline-block animate-spin rounded-full h-12 w-12 border-2 border-gray-200 border-t-blue-600"></div>
        <p className="mt-4 text-gray-500 font-medium">Loading...</p>
      </div>
    </div>
  );
}

export default Loading;
