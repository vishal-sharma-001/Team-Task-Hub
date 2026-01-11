import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authAPI } from '../api/client';
import { useForm } from '../hooks/useAsync';
import ErrorMessage from '../components/ErrorMessage';

function Login({ setUser }) {
  const navigate = useNavigate();
  const [apiError, setApiError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { values, errors, touched, handleChange, handleBlur, handleSubmit } =
    useForm(
      {
        email: '',
        password: '',
      },
      {
        email: (value) => {
          if (!value) return 'Email is required';
          if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
            return 'Please enter a valid email';
          }
          return '';
        },
        password: (value) => {
          if (!value) return 'Password is required';
          return '';
        },
      },
      async (formData) => {
        setIsLoading(true);
        setApiError('');
        try {
          console.log('Attempting login with:', formData.email);
          const response = await authAPI.login(formData);
          console.log('Login response:', response);
          
          if (!response || !response.token || !response.user) {
            throw new Error('Invalid response from server');
          }
          
          // Store data
          localStorage.setItem('authToken', response.token);
          localStorage.setItem('user', JSON.stringify(response.user));
          
          console.log('Stored user:', response.user);
          console.log('Calling setUser...');
          setUser(response.user);
          
          // Navigate after state update
          console.log('Navigating to dashboard...');
          setTimeout(() => {
            console.log('Navigation timeout executed');
            navigate('/dashboard');
          }, 100);
        } catch (error) {
          console.error('Login error:', error);
          setApiError(error.response?.data?.message || error.message || 'Login failed');
          setIsLoading(false);
        }
      }
    );

  return (
    <div className="min-h-screen bg-white flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-12">
          <div className="mb-10">
            <h1 className="text-4xl font-bold text-gray-900 mb-3">Welcome back</h1>
            <p className="text-gray-600 text-base">Sign in to Task Hub to manage your projects and tasks</p>
          </div>

          {apiError && <ErrorMessage message={apiError} />}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                Email
              </label>
              <input
                id="email"
                type="email"
                name="email"
                value={values.email}
                onChange={handleChange}
                onBlur={handleBlur}
                className={`input w-full ${
                  touched.email && errors.email ? 'border-red-500 focus:border-red-500 focus:ring-red-200' : ''
                }`}
                placeholder="your@email.com"
              />
              {touched.email && errors.email && (
                <p className="text-red-600 text-sm mt-1">{errors.email}</p>
              )}
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                Password
              </label>
              <input
                id="password"
                type="password"
                name="password"
                value={values.password}
                onChange={handleChange}
                onBlur={handleBlur}
                className={`input w-full ${
                  touched.password && errors.password ? 'border-red-500 focus:border-red-500 focus:ring-red-200' : ''
                }`}
                placeholder="••••••••"
              />
              {touched.password && errors.password && (
                <p className="text-red-600 text-sm mt-1">{errors.password}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Signing in...' : 'Sign in'}
            </button>
          </form>

          <div className="mt-8 pt-6 border-t border-gray-200">
            <p className="text-center text-gray-600 text-sm">
              Don't have an account?{' '}
              <Link to="/signup" className="text-gray-900 font-medium hover:opacity-80 transition-opacity">
                Sign up
              </Link>
            </p>
          </div>
        </div>

        <div className="mt-8 text-center text-gray-600 text-xs">
          <p>Demo: user@example.com / password</p>
        </div>
      </div>
    </div>
  );
}

export default Login;
