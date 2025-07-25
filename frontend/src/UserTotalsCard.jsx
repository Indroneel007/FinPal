import { useEffect, useState } from 'react';

export default function UserTotalsCard({ username, accessToken, onClick }) {
  const [totals, setTotals] = useState({ paid: 0, received: 0, net: 0, loading: true, error: '' });

  useEffect(() => {
    let isMounted = true;
    async function fetchTotals() {
      setTotals(t => ({ ...t, loading: true, error: '' }));
      try {
        const res = await fetch(`http://localhost:9090/transfers/${username}`, {
          headers: { 'Authorization': `Bearer ${accessToken}` },
        });
        if (!res.ok) {
          const errorData = await res.json().catch(() => ({}));
          throw new Error(errorData.error || 'Failed to fetch transaction totals');
        }
        const data = await res.json();
        const paid = Array.isArray(data.paid) ? data.paid.reduce((sum, tx) => sum + (tx.Amount || tx.amount || 0), 0) : 0;
        const received = Array.isArray(data.received) ? data.received.reduce((sum, tx) => sum + (tx.Amount || tx.amount || 0), 0) : 0;
        if (isMounted) {
          setTotals({ paid, received, net: received - paid, loading: false, error: '' });
        }
      } catch (err) {
        if (isMounted) setTotals(t => ({ ...t, loading: false, error: err.message }));
      }
    }
    fetchTotals();
    return () => { isMounted = false; };
  }, [username, accessToken]);

  if (totals.loading) {
    return (
      <div className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-3 border border-gray-700 animate-pulse">
        <div className="h-5 bg-gray-700 rounded w-1/2 mb-2" />
        <div className="h-4 bg-gray-700 rounded w-1/3 mb-1" />
        <div className="h-4 bg-gray-700 rounded w-1/3 mb-1" />
        <div className="h-4 bg-gray-700 rounded w-1/3" />
      </div>
    );
  }
  if (totals.error) {
    return (
      <div className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-3 border border-gray-700 text-red-400">
        Error: {totals.error}
      </div>
    );
  }
  return (
    <div 
      className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-3 border border-gray-700 hover:border-blue-500 transition-colors cursor-pointer"
      onClick={onClick}
    >
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-xl font-bold text-blue-300">{username}</h3>
        <span className={`text-lg font-bold ${totals.net >= 0 ? 'text-green-400' : 'text-red-400'}`}>
          {totals.net >= 0 ? '+' : ''}{totals.net.toLocaleString('en-IN')} ₹
        </span>
      </div>
      <div className="flex flex-col gap-1 text-sm text-gray-400">
        <div>
          <span className="text-green-400 font-semibold">Received: +₹{totals.received.toLocaleString('en-IN')}</span>
        </div>
        <div>
          <span className="text-red-400 font-semibold">Sent: -₹{totals.paid.toLocaleString('en-IN')}</span>
        </div>
      </div>
      <div className="mt-3 pt-2 border-t border-gray-700">
        <button
          className="w-full py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          onClick={e => { e.stopPropagation(); onClick && onClick(); }}
        >
          Send Money
        </button>
      </div>
    </div>
  );
}
