import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import './index.css';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Dashboard from './pages/Dashboard';
import Projects from './pages/Projects';
import TaskBoard from './pages/TaskBoard';
import TaskDetail from './pages/TaskDetail';
import ProtectedRoute from './components/ProtectedRoute';
import Navbar from './components/Navbar';

function App() {
  const [user, setUser] = useState(() => {
    try {
      const storedUser = localStorage.getItem('user');
      // Only parse if it's a non-empty string and not "undefined"
      if (storedUser && storedUser !== 'undefined' && storedUser !== 'null') {
        return JSON.parse(storedUser);
      }
    } catch (error) {
      console.error('Error parsing stored user:', error);
    }
    // Clear any bad data
    localStorage.removeItem('user');
    localStorage.removeItem('authToken');
    return null;
  });

  // Update localStorage whenever user changes
  useEffect(() => {
    if (user && typeof user === 'object') {
      localStorage.setItem('user', JSON.stringify(user));
    } else {
      localStorage.removeItem('user');
      localStorage.removeItem('authToken');
    }
  }, [user]);

  return (
    <Router>
      {user && <Navbar user={user} setUser={setUser} />}
      <main className={user ? 'pt-20' : ''}>
        <Routes>
          <Route path="/login" element={user ? <Navigate to="/dashboard" /> : <Login setUser={setUser} />} />
          <Route path="/signup" element={user ? <Navigate to="/dashboard" /> : <Signup setUser={setUser} />} />
          
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/projects"
            element={
              <ProtectedRoute>
                <Projects />
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/projects/:projectId/tasks"
            element={
              <ProtectedRoute>
                <TaskBoard />
              </ProtectedRoute>
            }
          />

          <Route
            path="/tasks/:taskId"
            element={
              <ProtectedRoute>
                <TaskDetail />
              </ProtectedRoute>
            }
          />

          <Route path="/" element={user ? <Navigate to="/dashboard" /> : <Navigate to="/login" />} />
        </Routes>
      </main>
    </Router>
  );
}

export default App;
