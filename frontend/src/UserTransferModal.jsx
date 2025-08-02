import { useState } from 'react';
import SendIcon from '@mui/icons-material/Send';

const SUPPORTED_CURRENCIES = ['USD', 'Euros', 'Rupees']; // Adjust as per backend util.IsSupportedCurrency
const SUPPORTED_TYPES = ['rent', 'food', 'travel', 'savings', 'bills','medical','shopping','misc']; // Adjust as per backend util.IsSupportedType

export default function UserTransferModal({ open, onClose, toUsername, onSubmit, loading, error }) {
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState(SUPPORTED_CURRENCIES[0]);
  const [type, setType] = useState(SUPPORTED_TYPES[0]);
  const [submitError, setSubmitError] = useState('');

  if (!open) return null;

  const handleSubmit = (e) => {
    e.preventDefault();
    setSubmitError('');
    if (!amount || isNaN(amount) || Number(amount) <= 0) {
      setSubmitError('Please enter a valid amount.');
      return;
    }
    onSubmit({
      to_username: toUsername,
      amount: Number(amount),
      currency,
      type,
    });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
      <div className="bg-gray-900 rounded-xl shadow-lg p-8 max-w-md w-full border border-gray-700 relative">
        <button
          className="absolute top-2 right-2 text-gray-400 hover:text-white text-xl"
          onClick={onClose}
          title="Close"
        >
          &times;
        </button>
        <h2 className="text-xl font-bold mb-4 text-center text-white">Send Money to <span className="text-blue-400">{toUsername}</span></h2>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <input
            type="number"
            min="1"
            step="1"
            placeholder="Amount"
            value={amount}
            onChange={e => setAmount(e.target.value)}
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
            required
          />
          <select
            value={currency}
            onChange={e => setCurrency(e.target.value)}
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          >
            {SUPPORTED_CURRENCIES.map(cur => (
              <option key={cur} value={cur}>{cur}</option>
            ))}
          </select>
          <select
            value={type}
            onChange={e => setType(e.target.value)}
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          >
            {SUPPORTED_TYPES.map(t => (
              <option key={t} value={t}>{t.charAt(0).toUpperCase() + t.slice(1)}</option>
            ))}
          </select>
          {submitError && <div className="text-red-400 text-sm text-center">{submitError}</div>}
          {error && <div className="text-red-400 text-sm text-center">{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-60"
          >
            {loading ? (<><SendIcon style={{ marginRight: 6, fontSize: 20 }} />Sending...</>) : (<><SendIcon style={{ marginRight: 6, fontSize: 20 }} />Send Money</>)}
          </button>
        </form>
      </div>
    </div>
  );
}
