import { useState } from 'react';

export default function CreateGroupModal({ open, onClose, onCreated, accessToken, username }) {
  const [groupName, setGroupName] = useState('');
  const [currency, setCurrency] = useState('INR');
  const [type, setType] = useState('savings');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  if (!open) return null;

  async function handleSubmit(e) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const res = await fetch('http://localhost:9090/groups', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          username,
          group_name: groupName,
          currency,
          type,
        })
      });
      if (!res.ok) {
        const errData = await res.json().catch(() => ({}));
        throw new Error(errData.error || 'Failed to create group');
      }
      const data = await res.json();
      setLoading(false);
      onCreated(data);
      onClose();
    } catch (err) {
      setLoading(false);
      setError(err.message);
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <form onSubmit={handleSubmit} className="bg-gray-900 rounded-xl p-8 w-full max-w-md shadow-xl flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-blue-300 mb-2">Create Group</h2>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-300">Group Name</span>
          <input
            type="text"
            value={groupName}
            onChange={e => setGroupName(e.target.value)}
            required
            className="bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white focus:outline-none focus:border-blue-400"
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-300">Currency</span>
          <select
            value={currency}
            onChange={e => setCurrency(e.target.value)}
            required
            className="bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white focus:outline-none focus:border-blue-400"
          >
            <option value="USD">USD</option>
            <option value="Euros">Euros</option>
            <option value="Rupees">Rupees</option>
          </select>
        </label>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-300">Type</span>
          <select
            value={type}
            onChange={e => setType(e.target.value)}
            required
            className="bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white focus:outline-none focus:border-blue-400"
          >
            <option value="rent">Rent</option>
            <option value="food">Food</option>
            <option value="travel">Travel</option>
            <option value="savings">Savings</option>
            <option value="bills">Bills</option>
            <option value="medical">Medical</option>
            <option value="shopping">Shopping</option>
            <option value="misc">Misc</option>
          </select>
        </label>
        {error && <div className="text-red-400 text-sm">{error}</div>}
        <div className="flex gap-2 mt-2">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
            disabled={loading}
          >Cancel</button>
          <button
            type="submit"
            className="flex-1 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-bold"
            disabled={loading}
          >{loading ? 'Creating...' : 'Create'}</button>
        </div>
      </form>
    </div>
  );
}
