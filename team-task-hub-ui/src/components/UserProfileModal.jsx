import { useState } from 'react';
import { authAPI } from '../api/client';
import Modal from './Modal';

function UserProfileModal({ isOpen, onClose, user, onProfileUpdate }) {
  const [name, setName] = useState(user?.name || '');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const updatedUser = await authAPI.updateProfile({ name });
      if (onProfileUpdate) {
        onProfileUpdate(updatedUser);
      }
      onClose();
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to update profile');
      console.error('Failed to update profile:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setName(user?.name || '');
    setError('');
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Edit Profile">
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">
            Email Address
          </label>
          <div className="p-3 rounded-lg bg-gray-50 border border-gray-200">
            <p className="text-sm text-gray-600 font-medium">{user?.email || 'No email'}</p>
            <p className="text-xs text-gray-500 mt-1">Email cannot be changed</p>
          </div>
        </div>

        <div className="space-y-2">
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">
            Display Name (Optional)
          </label>
          <input
            id="name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="input w-full"
            placeholder="Enter your name"
            maxLength={255}
          />
        </div>

        {error && (
          <div className="p-3 rounded-lg bg-red-50 border border-red-200">
            <p className="text-sm text-red-700">{error}</p>
          </div>
        )}

        <div className="flex gap-4 pt-4 border-t border-gray-200">
          <button
            type="button"
            onClick={handleClose}
            className="btn-secondary flex-1"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={loading}
            className="btn-primary flex-1"
          >
            {loading ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>
    </Modal>
  );
}

export default UserProfileModal;
