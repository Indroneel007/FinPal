import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

export default function ResetPassword() {
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const accessToken = location.state?.access_token;

  if (!accessToken) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#18181b]">
        <div className="bg-gray-900 p-8 rounded-2xl shadow-lg max-w-md border border-gray-700 text-center text-white">
          <h2 className="text-2xl font-bold mb-4">Invalid Access</h2>
          <p className="mb-4">Please restart the forgot password flow to reset your password.</p>
          <button
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 font-medium rounded-lg px-5 py-2.5"
            onClick={() => navigate('/forgot-password')}
          >
            Go to Forgot Password
          </button>
        </div>
      </div>
    );
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (!newPassword || !confirmPassword) {
      setError('Please fill all fields');
      return;
    }
    if (newPassword !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    if (!accessToken) {
      setError('Missing access token. Please restart the forgot password flow.');
      return;
    }
    setLoading(true);
    try {
      const res = await fetch('http://localhost:9090/resetpassword', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify({ new_password: newPassword })
      });
      if (!res.ok) {
        const err = await res.json();
        setError(err.error || 'Failed to reset password');
        setLoading(false);
        return;
      }
      // Success: redirect to landing page
      navigate('/');
    } catch (err) {
      setError('Network error: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#18181b]">
      <div className="bg-gray-900 p-8 rounded-2xl shadow-lg w-full max-w-md border border-gray-700">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">Reset Password</h2>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <input
            type="password"
            name="newPassword"
            placeholder="New Password"
            value={newPassword}
            onChange={e => setNewPassword(e.target.value)}
            required
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
          <input
            type="password"
            name="confirmPassword"
            placeholder="Confirm Password"
            value={confirmPassword}
            onChange={e => setConfirmPassword(e.target.value)}
            required
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
          {error && <div className="text-red-400 text-sm text-center">{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-60"
          >
            {loading ? 'Resetting...' : 'Reset Password'}
          </button>
        </form>
      </div>
    </div>
  );
}
