import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function ForgotPassword() {
  const [email, setEmail] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitted(false);
    try {
      const res = await fetch('https://finpal-1.onrender.com/forgotpassword', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      });
      if (!res.ok) {
        const err = await res.json();
        alert(err.error || 'Failed to send OTP');
        return;
      }
      // If email is valid, backend returns loginUserResponse
      const data = await res.json();
      const accessToken = data.access_token;
      navigate('/otp-verify', { state: { email, access_token: accessToken } });
    } catch (err) {
      alert('Network error: ' + err.message);
    }
  };


  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#18181b]">
      <div className="bg-gray-900 p-8 rounded-2xl shadow-lg w-full max-w-md border border-gray-700">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">Forgot Password</h2>
        {submitted ? (
          <div className="text-green-400 text-center">
            If the email exists, a reset link has been sent.
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <input
              type="email"
              name="email"
              placeholder="Enter your email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
              className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
            />
            <button
              type="submit"
              className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center"
            >
              Send OTP
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
