export default function LeaveGroupConfirmModal({ open, group, loading, error, onClose, onConfirm }) {
  if (!open) return null;
  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-xl p-8 w-full max-w-md shadow-xl flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-red-400 mb-2">Leave Group</h2>
        <div className="text-gray-200 mb-4">
          Are you sure you want to leave <span className="font-semibold text-red-300">{group?.group_name}</span>?
        </div>
        {error && <div className="text-red-400 text-sm">{error}</div>}
        <div className="flex gap-2 mt-2">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
            disabled={loading}
          >Cancel</button>
          <button
            type="button"
            onClick={onConfirm}
            className="flex-1 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors font-bold"
            disabled={loading}
          >{loading ? 'Leaving...' : 'Yes, Leave'}</button>
        </div>
      </div>
    </div>
  );
}
