import { useEffect } from 'react';

function Modal({ isOpen, onClose, title, children }) {
  useEffect(() => {
    if (isOpen) {
      // Disable body scroll when modal is open
      document.body.style.overflow = 'hidden';
    } else {
      // Re-enable body scroll when modal is closed
      document.body.style.overflow = 'unset';
    }

    // Cleanup on unmount
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto border border-gray-200">
        {/* Header */}
        <div className="flex justify-between items-center px-8 py-6 border-b border-gray-200 sticky top-0 bg-white rounded-t-xl">
          <h2 className="text-xl font-bold text-gray-900">{title}</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg p-1 transition-colors text-2xl leading-none"
          >
            ✕
          </button>
        </div>

        {/* Content */}
        <div className="p-8">
          {children}
        </div>
      </div>
    </div>
  );
}

export default Modal;
