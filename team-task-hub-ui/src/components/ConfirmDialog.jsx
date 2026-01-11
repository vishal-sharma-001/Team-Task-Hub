import { useEffect } from 'react';

function ConfirmDialog({ isOpen, title, message, confirmText = "Delete", cancelText = "Cancel", onConfirm, onCancel, isDanger = false }) {
  useEffect(() => {
    if (isOpen) {
      // Disable body scroll when dialog is open
      document.body.style.overflow = 'hidden';
    } else {
      // Re-enable body scroll when dialog is closed
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
      <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-gray-200">
        {/* Header */}
        <div className="px-8 py-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">{title}</h2>
        </div>

        {/* Content */}
        <div className="px-8 py-6">
          <p className="text-gray-600 text-base">{message}</p>
        </div>

        {/* Footer */}
        <div className="px-8 py-6 border-t border-gray-200 flex gap-3 justify-end">
          <button
            onClick={onCancel}
            className="btn-secondary"
          >
            {cancelText}
          </button>
          <button
            onClick={onConfirm}
            className={isDanger ? "btn-danger" : "btn-primary"}
          >
            {confirmText}
          </button>
        </div>
      </div>
    </div>
  );
}

export default ConfirmDialog;
