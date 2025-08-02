import React, { useState } from 'react';
import SendIcon from '@mui/icons-material/Send';
import CancelIcon from '@mui/icons-material/Cancel';

export default function GroupSendMoneyModal({ open, onClose, group, members, loading, error, onSubmit }) {
  const [toUsername, setToUsername] = useState('');
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState(group?.currency || '');
  const [type, setType] = useState(group?.type || '');
  const [submitError, setSubmitError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  if (!open) return null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitError('');
    setSubmitting(true);
    try {
      await onSubmit({ to_username: toUsername, amount: Number(amount), currency, type });
      setSubmitting(false);
      onClose();
    } catch (err) {
      setSubmitError(err.message || 'Failed to send money');
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-xl p-8 w-full max-w-md shadow-xl flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-green-300 mb-2">Send Money to Group Member</h2>
        {loading ? (
          <div className="text-white">Loading group members...</div>
        ) : error ? (
          <div className="text-red-400">{error}</div>
        ) : (
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <label className="text-white font-semibold">Select Member:
              <select
                className="block w-full mt-1 p-2 rounded bg-gray-800 text-white border border-gray-700"
                value={toUsername}
                onChange={e => setToUsername(e.target.value)}
                required
              >
                <option value="" disabled>Select user</option>
                {members.map((m) => (
                  <option key={m.owner} value={m.owner}>{m.owner}</option>
                ))}
              </select>
            </label>
            <label className="text-white font-semibold">Amount:
              <input
                type="number"
                min="1"
                className="block w-full mt-1 p-2 rounded bg-gray-800 text-white border border-gray-700"
                value={amount}
                onChange={e => setAmount(e.target.value)}
                required
              />
            </label>
            {submitError && <div className="text-red-400 text-sm">{submitError}</div>}
            <div className="flex gap-2 mt-2">
              <button
                type="submit"
                className="py-2 px-4 bg-green-600 hover:bg-green-700 text-white rounded-lg font-bold transition-colors disabled:opacity-60"
                disabled={submitting}
              >
                {submitting ? (<><SendIcon style={{ marginRight: 6, fontSize: 20 }} />Sending...</>) : (<><SendIcon style={{ marginRight: 6, fontSize: 20 }} />Send Money</>)}
              </button>
              <button
                type="button"
                className="py-2 px-4 bg-gray-600 hover:bg-gray-700 text-white rounded-lg font-bold transition-colors"
                onClick={onClose}
                disabled={submitting}
              > <CancelIcon style={{ marginRight: 6, fontSize: 20 }} />Cancel</button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
