import { useCallback, useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import UserTotalsCard from './UserTotalsCard';
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
  const [userTransactions, setUserTransactions] = useState({});
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

  const fetchUserTransactions = useCallback(async () => {
    if (!accessToken) {
      setError('Missing access token. Please login again.');
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const res = await fetch(`http://localhost:9090/transfers/user?page_id=${page}&page_size=${pageSize}`, {
        headers: { 'Authorization': `Bearer ${accessToken}` }
      });
      if (!res.ok) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error || 'Failed to fetch users with transactions');
      }
      const users = await res.json();
      if (!users || users.length === 0) {
        setUserTransactions({});
        return;
      }
      // Map to object for easier access by username
      const transactionsByUser = {};
      users.forEach(user => {
        transactionsByUser[user.username] = {
          username: user.username,
          total_sent: user.total_sent,
          total_received: user.total_received
        };
      });
      setUserTransactions(transactionsByUser);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [accessToken, page]);

  // Fetch all users with transactions
  useEffect(() => {
    fetchUserTransactions();
  }, [fetchUserTransactions, page]);

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
      {transferModal.open && (
        <AddUserTransferModal
          accessToken={accessToken}
          toUsername={transferModal.toUsername}
          onTransferSuccess={async () => {
            setTransferModal({ open: false, toUsername: '', loading: false, error: '' });
            if (historyModal.username) {
              setHistoryModal((prev) => ({ ...prev, loading: true }));
              try {
                const res = await fetch(`http://localhost:9090/transfers/${historyModal.username}`, {
                  headers: { 'Authorization': `Bearer ${accessToken}` }
                });
                if (!res.ok) {
                  const errorData = await res.json().catch(() => ({}));
                  throw new Error(errorData.error || 'Failed to fetch transaction history');
                }
                const data = await res.json();
                setHistoryModal({
                  open: true,
                  username: historyModal.username,
                  paid: data.paid || [],
                  received: data.received || [],
                  loading: false,
                  error: ''
                });
              } catch (err) {
                setHistoryModal((prev) => ({ ...prev, loading: false, error: err.message || 'Failed to fetch transaction history' }));
              }
            }
            fetchUserTransactions();
          }}
        />
      )}
      <Navbar username={username} showLogin={false} />
      <div className="flex flex-col gap-4 mt-6 mb-6 w-full">
        <button
          className="self-end px-5 py-2.5 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold shadow hover:from-blue-600 hover:to-purple-700 transition-colors"
          onClick={() => setTransferModal({ open: true, toUsername: '', loading: false, error: '' })}
        >
          + Add User
        </button>
      </div>
        
      {loading ? (
        <div className="text-white text-center">Loading accounts...</div>
      ) : error ? (
        <div className="text-red-400 text-center">{error}</div>
      ) : (
        <div className="flex-1">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mx-auto max-w-7xl pt-4">
            {Object.keys(userTransactions).length === 0 ? (
              <div className="col-span-full text-white text-center">No transactions found. Send or receive money to see transaction history.</div>
            ) : (
              Object.values(userTransactions)
                .filter((user) => user.username !== username)
                .map((user) => (
                  <UserTotalsCard
                    key={user.username}
                    username={user.username}
                    accessToken={accessToken}
                    onClick={async () => {
                      setHistoryModal({
                        open: true,
                        username: user.username,
                        paid: [],
                        received: [],
                        loading: true,
                        error: ''
                      });
                      try {
                        const res = await fetch(`http://localhost:9090/transfers/${user.username}`, {
                          headers: { 'Authorization': `Bearer ${accessToken}` }
                        });
                        if (!res.ok) {
                          const errorData = await res.json().catch(() => ({}));
                          throw new Error(errorData.error || 'Failed to fetch transaction history');
                        }
                        const data = await res.json();
                        setHistoryModal({
                          open: true,
                          username: user.username,
                          paid: data.paid || [],
                          received: data.received || [],
                          loading: false,
                          error: ''
                        });
                      } catch (err) {
                        setHistoryModal((prev) => ({
                          ...prev,
                          loading: false,
                          error: err.message || 'Failed to fetch transaction history'
                        }));
                      }
                    }}
                  />
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
          disabled={Object.keys(userTransactions).length < pageSize} 
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
                <div className="mb-6">
                  <h3 className="text-md font-semibold text-red-400 mb-2">Sent Money</h3>
                  
                </div>
                <div className="flex justify-between items-center mt-6 mb-2">
                  <h2 className="text-lg font-bold">Transaction History</h2>
                  <button
                    className="px-4 py-2 text-sm bg-blue-700 hover:bg-blue-800 text-white rounded"
                    onClick={() => setTransferModal({ open: true, toUsername: historyModal.username, loading: false, error: '' })}
                  >
                    Send Money
                  </button>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Paid Column */}
                  <div>
                    <h3 className="text-md font-semibold text-red-400 mb-2">Paid</h3>
                    {historyModal.paid.length === 0 ? (
                      <div className="text-gray-400 text-sm">No payments sent to this user.</div>
                    ) : (
                      <ul className="space-y-2">
                        {historyModal.paid.map((t, idx) => (
                          <li key={idx} className="flex flex-col p-3 bg-gray-800 rounded-lg">
                            <span className="text-xs text-gray-400 mb-1">{t.type}</span>
                            <span className="text-red-400 font-bold">-₹{t.amount}</span>
                            <span className="text-xs text-gray-500">{new Date(t.created_at).toLocaleString()}</span>
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                  {/* Received Column */}
                  <div>
                    <h3 className="text-md font-semibold text-green-400 mb-2">Received</h3>
                    {historyModal.received.length === 0 ? (
                      <div className="text-gray-400 text-sm">No payments received from this user.</div>
                    ) : (
                      <ul className="space-y-2">
                        {historyModal.received.map((t, idx) => (
                          <li key={idx} className="flex flex-col p-3 bg-gray-800 rounded-lg">
                            <span className="text-xs text-gray-400 mb-1">{t.type}</span>
                            <span className="text-green-400 font-bold">+₹{t.amount}</span>
                            <span className="text-xs text-gray-500">{new Date(t.created_at).toLocaleString()}</span>
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
