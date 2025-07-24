import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import UserTransferModal from './UserTransferModal';
import { MOCK_USERS } from './mockUsers';
import AddUserTransferModal from './AddUserTransferModal';
import Navbar from './Navbar';


export default function MainPage() {
  const location = useLocation();
  const navigate = useNavigate();
  // Try to get token and username from navigation state or fallback to global state if implemented
  const accessToken = location.state?.access_token || localStorage.getItem('access_token');
  const username = location.state?.username || localStorage.getItem('username');
  const [accounts, setAccounts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [historyModal, setHistoryModal] = useState({ open: false, username: '', paid: [], received: [], loading: false, error: '' });
  // User search and transfer modal state
  const [userQuery, setUserQuery] = useState('');
  const [userResults, setUserResults] = useState([]);
  const [transferModal, setTransferModal] = useState({ open: false, toUsername: '', loading: false, error: '' });
  const pageSize = 5;

  // User search effect
  useEffect(() => {
    if (userQuery.length < 2) {
      setUserResults([]);
      return;
    }
    // For now use mock user list; filter and exclude self
    setUserResults(
      MOCK_USERS.filter(
        (u) => u.toLowerCase().includes(userQuery.toLowerCase()) && u !== username
      )
    );
  }, [userQuery, username]);

  useEffect(() => {
    if (!accessToken) {
      setError('Missing access token. Please login again.');
      setLoading(false);
      return;
    }
    setLoading(true);
    fetch(`http://localhost:9090/accounts?page_id=${page}&page_size=${pageSize}`, {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch accounts');
        return res.json();
      })
      .then((data) => {
        setAccounts(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [accessToken, page]);

  if (!accessToken || !username) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#18181b] text-white">
        <div className="bg-gray-900 p-8 rounded-2xl shadow-lg max-w-md border border-gray-700 text-center">
          <h2 className="text-2xl font-bold mb-4">Invalid Access</h2>
          <p className="mb-4">Please login again to access your accounts.</p>
          <button
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 font-medium rounded-lg px-5 py-2.5"
            onClick={() => navigate('/')}
          >
            Go to Home
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-br from-purple-900 via-blue-900 to-gray-900 p-8">
      <Navbar username={username} showLogin={false} />
      <div className="mt-8 mb-6 flex justify-end w-full">
        <AddUserTransferModal
          accessToken={accessToken}
          onTransferSuccess={() => {
            setPage(1);
            setLoading(true);
          }}
        />
      </div>
      {loading ? (
        <div className="text-white text-center">Loading accounts...</div>
      ) : error ? (
        <div className="text-red-400 text-center">{error}</div>
      ) : (
        <div className="flex-1">
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 mx-auto max-w-5xl pt-4">
            {accounts.length === 0 ? (
              <div className="col-span-full text-white text-center">No accounts found.</div>
            ) : (
              accounts.map((acc) => (
                <div key={acc.id} className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-2 border border-gray-700 w-full max-w-lg mx-auto">
                  <div className="flex items-center justify-between mb-2">
                    <span
                      className={`text-lg font-semibold cursor-pointer underline text-blue-300 hover:text-blue-400 transition`}
                      title={acc.owner === username ? 'Your account' : `View transfers with ${acc.owner}`}
                      onClick={() => {
                        if (acc.owner === username) return;
                        setHistoryModal((prev) => ({ ...prev, open: true, username: acc.owner, paid: [], received: [], loading: true, error: '' }));
                        fetch(`http://localhost:9090/transfers/${acc.owner}`, {
                          headers: { 'Authorization': `Bearer ${accessToken}` }
                        })
                          .then(res => {
                            if (!res.ok) throw new Error('Failed to fetch transfer history');
                            return res.json();
                          })
                          .then(data => {
                            setHistoryModal((prev) => ({ ...prev, paid: data.paid, received: data.received, loading: false, error: '' }));
                          })
                          .catch(err => {
                            setHistoryModal((prev) => ({ ...prev, loading: false, error: err.message }));
                          });
                      }}
                    >
                      {acc.owner === username ? 'You' : acc.owner}
                    </span>
                    <span className="text-sm text-gray-400">ID: {acc.id}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-2xl font-bold text-green-400">₹{acc.balance}</span>
                    <span className="text-sm text-gray-400">{acc.currency}</span>
                  </div>
                  <div className="flex items-center justify-between mt-2">
                    <span className="text-sm text-blue-400">Type: {acc.type}</span>
                    {acc.group_id && acc.group_id.Valid && (
                      <span className="text-sm text-pink-400">Group: {acc.group_id.Int64}</span>
                    )}
                    {acc.has_accepted !== undefined && (
                      <span className={`text-xs font-bold rounded px-2 py-1 ${acc.has_accepted ? 'bg-green-700 text-green-200' : 'bg-yellow-700 text-yellow-200'}`}>
                        {acc.has_accepted ? 'Accepted' : 'Pending'}
                      </span>
                    )}
                  </div>
                  <div className="text-xs text-gray-500 mt-2">Created: {new Date(acc.created_at).toLocaleString()}</div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
      {/* Pagination Controls fixed to footer */}
      <div className="w-full fixed bottom-0 left-0 flex justify-center items-center bg-gradient-to-r from-purple-900 via-blue-900 to-gray-900 py-4 gap-4 z-10 border-t border-gray-700">
        <button
          className="px-4 py-2 rounded bg-gray-700 text-white disabled:opacity-50"
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page === 1}
        >
          Previous
        </button>
        <span className="text-white px-4 py-2">Page {page}</span>
        <button
          className="px-4 py-2 rounded bg-gray-700 text-white"
          onClick={() => setPage((p) => p + 1)}
          disabled={accounts.length < pageSize}
        >
          Next
        </button>
      </div>
      {/* Transfer History Modal */}
      {historyModal.open && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
          <div className="bg-gray-900 rounded-xl shadow-lg p-8 max-w-md w-full border border-gray-700 relative">
            <button
              className="absolute top-2 right-2 text-gray-400 hover:text-white text-xl"
              onClick={() => setHistoryModal({ open: false, username: '', paid: [], received: [], loading: false, error: '' })}
              title="Close"
            >
              &times;
            </button>
            <h2 className="text-xl font-bold mb-4 text-center text-white">Transfers with <span className="text-blue-400">{historyModal.username}</span></h2>
            {historyModal.loading ? (
              <div className="text-white text-center">Loading...</div>
            ) : historyModal.error ? (
              <div className="text-red-400 text-center">{historyModal.error}</div>
            ) : (
              <>
                <div className="mb-4">
                  <h3 className="text-md font-semibold text-green-400 mb-2">Paid</h3>
                  {historyModal.paid.length === 0 ? (
                    <div className="text-gray-400 text-sm">No payments made to this user.</div>
                  ) : (
                    <ul className="space-y-1">
                      {historyModal.paid.map((t, idx) => (
                        <li key={idx} className="flex justify-between text-sm">
                          <span>₹{t.amount}</span>
                          <span className="text-gray-400">{new Date(t.created_at).toLocaleString()}</span>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
                <div>
                  <h3 className="text-md font-semibold text-blue-400 mb-2">Received</h3>
                  {historyModal.received.length === 0 ? (
                    <div className="text-gray-400 text-sm">No payments received from this user.</div>
                  ) : (
                    <ul className="space-y-1">
                      {historyModal.received.map((t, idx) => (
                        <li key={idx} className="flex justify-between text-sm">
                          <span>₹{t.amount}</span>
                          <span className="text-gray-400">{new Date(t.created_at).toLocaleString()}</span>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
