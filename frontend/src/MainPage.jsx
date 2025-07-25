import { useCallback, useEffect, useState } from 'react';
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
      <Navbar username={username} showLogin={false} />
      <div className="mt-8 mb-6 flex justify-between items-center w-full">
        <h1 className="text-2xl font-bold text-white">Your Transactions</h1>
        <AddUserTransferModal
          accessToken={accessToken}
          onTransferSuccess={() => {
            setPage(1);
            setLoading(true);
            // Refetch transactions after successful transfer
            fetchUserTransactions();
          }}
        />
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
                .map((user) => {
                  const totalSent = user.total_sent || 0;
                  const totalReceived = user.total_received || 0;
                  const netAmount = totalReceived - totalSent;
                  return (
                    <div 
                      key={user.username}
                      className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-3 border border-gray-700 hover:border-blue-500 transition-colors cursor-pointer"
                      onClick={() => {
                        setHistoryModal({
                          open: true,
                          username: user.username,
                          paid: user.paid,
                          loading: false,
                          error: ''
                        });
                      }}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <h3 className="text-xl font-bold text-blue-300">
                          {user.username}
                        </h3>
                        <span className={`text-lg font-bold ${netAmount >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                          {netAmount >= 0 ? '+' : ''}{netAmount.toLocaleString('en-IN')} ₹
                        </span>
                      </div>
                      <div className="flex justify-between text-sm text-gray-400">
                        <div>
                          <div className="text-red-400">Sent: -₹{totalSent.toLocaleString('en-IN')}</div>
                          <div className="text-green-400">Received: +₹{totalReceived.toLocaleString('en-IN')}</div>
                        </div>
                      </div>
                      <div className="mt-3 pt-2 border-t border-gray-700">
                        <button
                          className="w-full py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                          onClick={(e) => {
                            e.stopPropagation();
                            // Open transfer modal with this user
                            setTransferModal({
                              open: true,
                              toUsername: user.username,
                              loading: false,
                              error: ''
                            });
                          }}
                        >
                          Send Money
                        </button>
                      </div>
                    </div>
                  );
                })
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
                  {historyModal.paid.length === 0 ? (
                    <div className="text-gray-400 text-sm">No payments made to this user.</div>
                  ) : (
                    <ul className="space-y-2">
                      {historyModal.paid.map((t, idx) => (
                        <li key={idx} className="flex flex-col p-3 bg-gray-800 rounded-lg">
                          <div className="flex justify-between items-center">
                            <span className="text-red-300 font-medium">To: {t.to_username || 'Unknown'}</span>
                            <span className="text-red-400 font-bold">-₹{t.amount}</span>
                          </div>
                          <div className="flex justify-between text-xs text-gray-400 mt-1">
                            <span>{t.type}</span>
                            <span>{new Date(t.created_at).toLocaleString()}</span>
                          </div>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
                <div className="mt-6">
                  <h3 className="text-md font-semibold text-green-400 mb-2">Received Money</h3>
                  {historyModal.received.length === 0 ? (
                    <div className="text-gray-400 text-sm">No payments received from this user.</div>
                  ) : (
                    <ul className="space-y-2">
                      {historyModal.received.map((t, idx) => (
                        <li key={idx} className="flex flex-col p-3 bg-gray-800 rounded-lg">
                          <div className="flex justify-between items-center">
                            <span className="text-green-300 font-medium">From: {t.from_username || 'Unknown'}</span>
                            <span className="text-green-400 font-bold">+₹{t.amount}</span>
                          </div>
                          <div className="flex justify-between text-xs text-gray-400 mt-1">
                            <span>{t.type}</span>
                            <span>{new Date(t.created_at).toLocaleString()}</span>
                          </div>
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
