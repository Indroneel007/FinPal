import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

export default function OtpVerify() {
  const [otp, setOtp] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const location = useLocation();
  const email = location.state?.email || '';

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!otp) {
      setError('Please enter the OTP');
      return;
    }
    setError('');
    try {
      const accessToken = location.state?.access_token;
      if (!accessToken) {
        setError('Missing access token. Please restart the forgot password flow.');
        return;
      }
      const res = await fetch('https://finpal-1.onrender.com/checkotp', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify({ otp })
      });
      if (!res.ok) {
        const err = await res.json();
        setError(err.error || 'OTP verification failed');
        return;
      }
      // On success, navigate to reset password page
      navigate('/reset-password', { state: { email, access_token: accessToken } });
    } catch (err) {
      setError('Network error: ' + err.message);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#18181b]">
      <div className="bg-gray-900 p-8 rounded-2xl shadow-lg w-full max-w-md border border-gray-700">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">Verify OTP</h2>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <input
            type="text"
            name="otp"
            placeholder="Enter OTP"
            value={otp}
            onChange={e => setOtp(e.target.value)}
            required
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
          {error && <div className="text-red-400 text-sm text-center">{error}</div>}
          <button
            type="submit"
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center"
          >
            Verify OTP
          </button>
        </form>
      </div>
    </div>
  );
}
