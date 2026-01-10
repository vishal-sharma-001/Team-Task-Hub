import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authAPI } from '../api/client';
import { useForm } from '../hooks/useAsync';
import ErrorMessage from '../components/ErrorMessage';

function Signup({ setUser }) {
  const navigate = useNavigate();
  const [apiError, setApiError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { values, errors, touched, handleChange, handleBlur, handleSubmit } =
    useForm(
      {
        email: '',
        password: '',
        confirmPassword: '',
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
          if (value.length < 8) return 'Password must be at least 8 characters';
          return '';
        },
        confirmPassword: (value) => {
          if (!value) return 'Please confirm your password';
          return '';
        },
      },
      async (formData) => {
        if (formData.password !== formData.confirmPassword) {
          setApiError('Passwords do not match');
          return;
        }

        setIsLoading(true);
        setApiError('');
        try {
          console.log('Attempting signup with:', formData.email);
          const response = await authAPI.signup(formData);
          console.log('Signup response:', response);
          
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
          console.error('Signup error:', error);
          setApiError(error.response?.data?.message || 'Signup failed');
          setIsLoading(false);
        }
      }
    );

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-12">
          <div className="mb-10">
            <h1 className="text-4xl font-bold text-gray-900 mb-3">Get started</h1>
            <p className="text-gray-600 text-base">Create an account to start organizing your tasks</p>
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

            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700">
                Confirm Password
              </label>
              <input
                id="confirmPassword"
                type="password"
                name="confirmPassword"
                value={values.confirmPassword}
                onChange={handleChange}
                onBlur={handleBlur}
                className={`input w-full ${
                  touched.confirmPassword && errors.confirmPassword
                    ? 'border-red-500 focus:border-red-500 focus:ring-red-200'
                    : ''
                }`}
                placeholder="••••••••"
              />
              {touched.confirmPassword && errors.confirmPassword && (
                <p className="text-red-600 text-sm mt-1">{errors.confirmPassword}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Creating account...' : 'Create account'}
            </button>
          </form>

          <div className="mt-8 pt-6 border-t border-gray-200">
            <p className="text-center text-gray-600 text-sm">
              Already have an account?{' '}
              <Link to="/login" className="text-blue-600 font-medium hover:text-blue-700 transition-colors">
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Signup;
