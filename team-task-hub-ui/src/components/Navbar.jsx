import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import UserProfileModal from './UserProfileModal';

function Navbar({ user, setUser }) {
  const navigate = useNavigate();
  const [showProfileModal, setShowProfileModal] = useState(false);

  const handleLogout = () => {
    setUser(null);
    localStorage.removeItem('user');
    localStorage.removeItem('authToken');
    navigate('/login');
  };

  const handleProfileUpdate = (updatedUser) => {
    setUser(updatedUser);
  };

  if (!user) return null;

  return (
    <>
      <nav className="fixed top-0 left-0 right-0 bg-white border-b border-gray-200 z-50 shadow-sm">
        <div className="container px-6 py-4 flex justify-between items-center">
          <div className="flex items-center gap-12">
            <button 
              onClick={() => navigate('/')} 
              className="text-xl font-bold text-gray-900 hover:text-blue-600 transition-colors cursor-pointer bg-none border-none p-0"
            >
              📋 Task Hub
            </button>
            <div className="hidden md:flex gap-8">
              <button 
                onClick={() => navigate('/')} 
                className="text-gray-600 hover:text-gray-900 font-medium transition-colors cursor-pointer bg-none border-none p-0 text-sm"
              >
                Dashboard
              </button>
              <button 
                onClick={() => navigate('/projects')} 
                className="text-gray-600 hover:text-gray-900 font-medium transition-colors cursor-pointer bg-none border-none p-0 text-sm"
              >
                Projects
              </button>
            </div>
          </div>

          <div className="flex items-center gap-6">
            <button
              onClick={() => setShowProfileModal(true)}
              className="flex items-center gap-2 px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer"
            >
              <span className="text-gray-600">👤</span>
              <span className="text-gray-700 text-sm font-medium">{user.email}</span>
            </button>
            <button
              onClick={handleLogout}
              className="btn-secondary"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <UserProfileModal
        isOpen={showProfileModal}
        onClose={() => setShowProfileModal(false)}
        user={user}
        onProfileUpdate={handleProfileUpdate}
      />
    </>
  );
}

export default Navbar;
