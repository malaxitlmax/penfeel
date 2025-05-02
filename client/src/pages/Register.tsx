import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate, Link } from 'react-router-dom';
import { API_URL } from '@/App';

interface RegisterCredentials {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

interface RegisterResponse {
  userId: string;
  message: string;
}

export default function Register() {
  const [credentials, setCredentials] = useState<RegisterCredentials>({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const registerMutation = useMutation<RegisterResponse, Error, RegisterCredentials>({
    mutationFn: async (credentials) => {
      // Client-side validation
      if (credentials.password !== credentials.confirmPassword) {
        throw new Error('Passwords do not match');
      }

      // Create registration data without confirmPassword
      const registrationData = {
        username: credentials.username,
        email: credentials.email,
        password: credentials.password,
      };

      try {
        const response = await fetch(API_URL + '/auth/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(registrationData),
        });

        if (!response.ok) {
          // Handle different HTTP error status codes
          if (response.status === 409) {
            throw new Error('An account with this email already exists.');
          } else if (response.status === 400) {
            // Try to get validation errors
            try {
              const errorData = await response.json();
              // Check if there are specific field validation errors
              if (errorData.errors) {
                const errorMessages = Object.values(errorData.errors).join('. ');
                throw new Error(errorMessages || 'Invalid registration information provided.');
              }
              throw new Error(errorData.message || 'Registration failed due to invalid data.');
            } catch {
              throw new Error('Registration failed due to invalid data.');
            }
          } else if (response.status === 429) {
            throw new Error('Too many registration attempts. Please try again later.');
          }

          // Try to get detailed error message from response
          try {
            const errorData = await response.json();
            throw new Error(errorData.message || `Registration failed with status: ${response.status}`);
          } catch {
            // If parsing JSON fails, use status text
            throw new Error(`Registration failed: ${response.statusText || response.status}`);
          }
        }

        return response.json();
      } catch (error) {
        // Handle network errors and other exceptions
        if (error instanceof TypeError && error.message === 'Failed to fetch') {
          throw new Error('Unable to connect to the server. Please check your internet connection and try again.');
        } else if (error instanceof Error) {
          throw error; // Re-throw if it's already an Error object with message
        } else {
          throw new Error('An unexpected error occurred. Please try again later.');
        }
      }
    },
    onSuccess: () => {
      // Redirect to login page after successful registration
      navigate('/login');
    },
    onError: (error) => {
      setError(error.message);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    registerMutation.mutate(credentials);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCredentials((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-8 bg-white rounded-lg shadow-md">
        <h2 className="text-3xl font-bold text-center text-gray-800 mb-6">Create Account</h2>
        
        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
              Username
            </label>
            <input
              id="username"
              name="username"
              type="text"
              required
              value={credentials.username}
              onChange={handleInputChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="johndoe"
            />
          </div>
          
          <div className="mb-4">
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              id="email"
              name="email"
              type="email"
              required
              value={credentials.email}
              onChange={handleInputChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="your@email.com"
            />
          </div>
          
          <div className="mb-4">
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              id="password"
              name="password"
              type="password"
              required
              value={credentials.password}
              onChange={handleInputChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="••••••••"
              minLength={8}
            />
          </div>
          
          <div className="mb-6">
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-1">
              Confirm Password
            </label>
            <input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              required
              value={credentials.confirmPassword}
              onChange={handleInputChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="••••••••"
            />
          </div>
          
          <button
            type="submit"
            disabled={registerMutation.isPending}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
          >
            {registerMutation.isPending ? 'Creating account...' : 'Create Account'}
          </button>
        </form>
        
        <div className="mt-4 text-center text-sm text-gray-600">
          Already have an account?{' '}
          <Link to="/login" className="text-blue-600 hover:underline">
            Sign in
          </Link>
        </div>
      </div>
    </div>
  );
} 