import { useState } from 'react';
import CancelIcon from '@mui/icons-material/Cancel';
import EditIcon from '@mui/icons-material/Edit';

export default function UpdateGroupNameModal({ open, onClose, onSubmit, loading, error, group }) {
  const [groupname, setGroupname] = useState('');

  if (!open) return null;

  const handleSubmit = (e) => {
    e.preventDefault();
    if (groupname.trim()) {
      onSubmit({
        new_name: groupname.trim(),
      });
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <form onSubmit={handleSubmit} className="bg-gray-900 rounded-xl p-8 w-full max-w-md shadow-xl flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-blue-300 mb-2">Update Group Name</h2>
        <label className="flex flex-col gap-1">
          <input
            type="text"
            value={groupname}
            onChange={e => setGroupname(e.target.value)}
            required
            className="bg-gray-800 border border-gray-700 rounded px-3 py-2 text-white focus:outline-none focus:border-blue-400"
            placeholder="Enter group name"
          />
        </label>
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
          >{loading ? (<><EditIcon style={{ marginRight: 6, fontSize: 20 }} />Updating...</>) : (<><EditIcon style={{ marginRight: 6, fontSize: 20 }} />Update Name</>)}</button>
        </div>
      </form>
    </div>
  );
}
