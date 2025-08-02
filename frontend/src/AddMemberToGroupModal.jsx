import { useState } from 'react';
import CancelIcon from '@mui/icons-material/Cancel';
import PersonAddIcon from '@mui/icons-material/PersonAdd';

export default function AddMemberToGroupModal({ open, onClose, onSubmit, loading, error, group }) {
  const [username, setUsername] = useState('');

  if (!open) return null;

  const handleSubmit = (e) => {
    e.preventDefault();
    if (username.trim()) {
      onSubmit({
        username: username.trim(),
        currency: group.currency,
        type: group.type,
      });
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <form onSubmit={handleSubmit} className="bg-gray-900 rounded-xl p-8 w-full max-w-md shadow-xl flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-blue-300 mb-2">Add Member to {group.group_name}</h2>
        <label className="flex flex-col gap-1">
          <span className="text-sm text-gray-300">Username</span>
          <input
            type="text"
            value={username}
            onChange={e => setUsername(e.target.value)}
            required
            className="bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white focus:outline-none focus:border-blue-400"
            placeholder="Enter username"
          />
        </label>
        <div className="flex gap-2">
          <div className="flex-1">
            <span className="block text-xs text-gray-400">Currency</span>
            <span className="block text-white font-semibold">{group.currency}</span>
          </div>
          <div className="flex-1">
            <span className="block text-xs text-gray-400">Type</span>
            <span className="block text-white font-semibold">{group.type}</span>
          </div>
        </div>
        {error && <div className="text-red-400 text-sm">{error}</div>}
        <div className="flex gap-2 mt-2">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
            disabled={loading}
          > <CancelIcon style={{ marginRight: 6, fontSize: 20 }} />Cancel</button>
          <button
            type="submit"
            className="flex-1 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-bold"
            disabled={loading}
          >{loading ? (<><PersonAddIcon style={{ marginRight: 6, fontSize: 20 }} />Adding...</>) : (<><PersonAddIcon style={{ marginRight: 6, fontSize: 20 }} />Add Member</>)}</button>
        </div>
      </form>
    </div>
  );
}
