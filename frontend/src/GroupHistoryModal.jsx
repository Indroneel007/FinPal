import React from 'react';

export default function GroupHistoryModal({ open, onClose, history, loading, error }) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-xl p-8 w-full max-w-lg shadow-xl flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-blue-300 mb-2">Group Transaction History</h2>
        {loading ? (
          <div className="text-white">Loading...</div>
        ) : error ? (
          <div className="text-red-400">{error}</div>
        ) : history.length === 0 ? (
          <div className="text-gray-400">No transactions found for this group.</div>
        ) : (
          <div className="overflow-y-auto max-h-96">
            <table className="min-w-full text-white border-separate border-spacing-y-2">
              <thead>
                <tr>
                  <th className="text-left">Amount</th>
                  <th className="text-left">From</th>
                  <th className="text-left">To</th>
                  <th className="text-left">Time</th>
                </tr>
              </thead>
              <tbody>
                {history.map(tx => (
                  <tr key={tx.transfer_id} className="bg-gray-800 hover:bg-gray-700">
                    <td className="py-2 px-3">{tx.amount}</td>
                    <td className="py-2 px-3">{tx.from_username}</td>
                    <td className="py-2 px-3">{tx.to_username}</td>
                    <td className="py-2 px-3 text-xs">{new Date(tx.created_at).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <button
          type="button"
          className="mt-4 py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-bold"
          onClick={onClose}
        >Close</button>
      </div>
    </div>
  );
}
