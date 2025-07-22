import { motion } from 'framer-motion';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function AuthCard({ mode, onClose, onSwitch }) {
  const navigate = useNavigate();
  const [form, setForm] = useState({ username: '', email: '', password: '', confirmPassword: '', fullName: '', salary: '' });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (mode === 'login') {
      try {
        const res = await fetch('http://localhost:9090/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: form.username,
            password: form.password
          })
        });
        if (!res.ok) {
          const err = await res.json();
          alert(err.error || 'Login failed');
          return;
        }
        const data = await res.json();
        alert('Login successful!');
        onClose();
        // Optionally: store token, update user state, etc.
      } catch (err) {
        alert('Network error: ' + err.message);
      }
    }
    if (mode === 'register') {
      if (form.password !== form.confirmPassword) {
        alert('Passwords do not match!');
        return;
      }
      if (!form.username || !form.email || !form.fullName || !form.salary) {
        alert('Please fill in all fields.');
        return;
      }
      try {
        const res = await fetch('http://localhost:9090/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: form.username,
            password: form.password,
            full_name: form.fullName,
            email: form.email,
            salary: Number(form.salary)
          })
        });
        if (!res.ok) {
          const err = await res.json();
          alert(err.error || 'Registration failed');
          return;
        }
        alert('Registration successful!');
        onClose();
      } catch (err) {
        alert('Network error: ' + err.message);
      }
    }
  };

  return (
    <motion.div
      initial={{ y: -100, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      exit={{ y: -100, opacity: 0 }}
      transition={{ type: 'spring', stiffness: 300, damping: 30 }}
      className="fixed inset-0 z-30 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    >
      <motion.div
        initial={{ scale: 0.9 }}
        animate={{ scale: 1 }}
        className="bg-[#18181b] rounded-2xl shadow-2xl p-8 w-[90vw] max-w-md border border-gray-700 relative"
      >
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-white text-2xl"
          aria-label="Close"
        >
          &times;
        </button>
        <h2 className="text-2xl font-bold text-white mb-6 text-center">
          {mode === 'login' ? 'Login to FinPal' : 'Create your FinPal account'}
        </h2>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          {mode === 'login' ? (
            <input
              type="text"
              name="username"
              placeholder="Username"
              value={form.username}
              onChange={handleChange}
              required
              className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
            />
          ) : (
            <>
              <input
                type="text"
                name="username"
                placeholder="Username"
                value={form.username}
                onChange={handleChange}
                required
                className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              />
              <input
                type="email"
                name="email"
                placeholder="Email"
                value={form.email}
                onChange={handleChange}
                required
                className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              />
            </>
          )}
          <input
            type="password"
            name="password"
            placeholder="Password"
            value={form.password}
            onChange={handleChange}
            required
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
          {mode === 'register' && (
            <>
              <input
                type="password"
                name="confirmPassword"
                placeholder="Confirm Password"
                value={form.confirmPassword}
                onChange={handleChange}
                required
                className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              />
              <input
                type="text"
                name="fullName"
                placeholder="Full Name"
                value={form.fullName}
                onChange={handleChange}
                required
                className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              />
              <input
                type="number"
                name="salary"
                placeholder="Salary"
                value={form.salary}
                onChange={handleChange}
                required
                min="0"
                className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              />
            </>
          )}
          <button
            type="submit"
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center me-2 mb-2"
          >
            {mode === 'login' ? 'Login' : 'Register'}
          </button>
        </form>
        {mode === 'login' && (
          <div className="text-center mt-4 flex flex-col gap-2">
            <button
              type="button"
              onClick={onSwitch}
              className="text-violet-400 hover:underline text-sm"
            >
              New User?
            </button>
            <button
              type="button"
              onClick={() => {
                onClose();
                navigate('/forgot-password');
              }}
              className="text-violet-400 hover:underline text-sm"
            >
              Forgot password?
            </button>
          </div>
        )}
        {mode === 'register' && (
          <div className="text-center mt-4">
            <button
              type="button"
              onClick={onSwitch}
              className="text-violet-400 hover:underline text-sm"
            >
              Already have an account? Login
            </button>
          </div>
        )}
      </motion.div>
    </motion.div>
  );
}
