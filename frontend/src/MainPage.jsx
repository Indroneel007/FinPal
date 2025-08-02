import { useCallback, useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import UserTotalsCard from './UserTotalsCard';
import UserTransferModal from './UserTransferModal';
import { MOCK_USERS } from './mockUsers';
import AddUserTransferModal from './AddUserTransferModal';
import Navbar from './Navbar';
import CreateGroupModal from './CreateGroupModal';
import GroupCard from './GroupCard';
import TransactionTypeSelector from './TransactionTypeSelector';
import AddMemberToGroupModal from './AddMemberToGroupModal';
import UpdateGroupNameModal from './UpdateGroupNameModal';
import LeaveGroupConfirmModal from './LeaveGroupConfirmModal';
import GroupHistoryModal from './GroupHistoryModal';
import GroupSendMoneyModal from './GroupSendMoneyModal';
import PromptSidebar from './PromptSidebar';

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

  // Group modal and state
  const [groupModalOpen, setGroupModalOpen] = useState(false);
  const [groups, setGroups] = useState([]);

  // Modal state for Add Member
  const [addMemberModal, setAddMemberModal] = useState({ open: false, group: null, loading: false, error: '' });
  // Modal state to update groupname
  const [newGroupNameModal, setNewGroupNameModal] = useState({open: false, group: null, loading: false, error:''});
  // Modal state for leave group
  const [leaveGroupModal, setLeaveGroupModal] = useState({ open: false, group: null, loading: false, error: '' });
  // Right-side selector state
  const [transactionType, setTransactionType] = useState('user');
  // Group History state
  const [groupHistoryModal, setGroupHistoryModal] = useState({ open: false, group: null, history: [], loading: false, error: '' });
  // Group Send Money state
  const [groupSendMoneyModal, setGroupSendMoneyModal] = useState({ open: false, group: null, members: [], loading: false, error: '' });
  // Prompt sidebar state
  const [promptSidebar, setPromptSidebar] = useState({ open: false, loading: false, prompt: '', error: '' });
  // Mindset dropdown state
  const [mindset, setMindset] = useState('medium');

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
    if (transactionType === 'user') {
      fetchUserTransactions();
    } else if (transactionType === 'group') {
      // Fetch groups from backend
      const fetchGroups = async () => {
        setLoading(true);
        setError('');
        try {
          const res = await fetch(`http://localhost:9090/groups?page_id=1&page_size=5`, {
            headers: { 'Authorization': `Bearer ${accessToken}` },
          });
          if (!res.ok) {
            const errorData = await res.json().catch(() => ({}));
            throw new Error(errorData.error || 'Failed to fetch groups');
          }
          const data = await res.json();
          setGroups(Array.isArray(data) ? data : []);
        } catch (err) {
          setError(err.message);
        } finally {
          setLoading(false);
        }
      };
      fetchGroups();
    }
  }, [fetchUserTransactions, page, transactionType, accessToken]);

  if (!accessToken || !username) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#18181b] text-white">
        <div className="bg-gray-900 p-8 rounded-2xl shadow-lg max-w-md w-full border border-gray-700 text-center">
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
      <div className="flex flex-row gap-4 mt-6 mb-6 w-full">
        <div className="flex-1 flex gap-4 justify-end">
          <button
            className="px-5 py-2.5 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold shadow hover:from-blue-600 hover:to-purple-700 transition-colors"
            onClick={() => setTransferModal({ open: true, toUsername: '', loading: false, error: '' })}
          >
            + Add User
          </button>
          <button
            className="px-5 py-2.5 rounded-lg bg-gradient-to-r from-green-500 to-teal-600 text-white font-semibold shadow hover:from-green-600 hover:to-teal-700 transition-colors"
            onClick={() => setGroupModalOpen(true)}
          >
            + Create Group
          </button>
        </div>
        <div className="flex items-center">
          <TransactionTypeSelector value={transactionType} onChange={setTransactionType} />
        </div>
      </div>
      <CreateGroupModal
        open={groupModalOpen}
        onClose={() => setGroupModalOpen(false)}
        accessToken={accessToken}
        username={username}
        onCreated={async () => {
          // Always fetch the latest groups after creation
          const res = await fetch(`http://localhost:9090/groups?page_id=1&page_size=5`, {
            headers: { 'Authorization': `Bearer ${accessToken}` },
          });
          const data = await res.json();
          setGroups(Array.isArray(data) ? data : []);
        }}
      />
      <AddMemberToGroupModal
        open={addMemberModal.open}
        group={addMemberModal.group || {}}
        loading={addMemberModal.loading}
        error={addMemberModal.error}
        onClose={() => setAddMemberModal({ open: false, group: null, loading: false, error: '' })}
        onSubmit={async ({ username, currency, type }) => {
          setAddMemberModal(m => ({ ...m, loading: true, error: '' }));
          try {
            const groupId = (addMemberModal.group.group_id || addMemberModal.group.id);
            const res = await fetch(`http://localhost:9090/groups/${groupId}/add`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`,
              },
              body: JSON.stringify({ username, currency, type })
            });
            if (!res.ok) {
              const errData = await res.json().catch(() => ({}));
              throw new Error(errData.error || 'Failed to add member');
            }
            setAddMemberModal({ open: false, group: null, loading: false, error: '' });
            // Refresh groups list
            if (transactionType === 'group') {
              const res = await fetch(`http://localhost:9090/groups?page_id=1&page_size=5`, {
                headers: { 'Authorization': `Bearer ${accessToken}` },
              });
              const data = await res.json();
              setGroups(Array.isArray(data) ? data : []);
            }
          } catch (err) {
            setAddMemberModal(m => ({ ...m, loading: false, error: err.message }));
          }
        }}
      />

      <UpdateGroupNameModal
        open={newGroupNameModal.open}
        group={newGroupNameModal.group || {}}
        loading={newGroupNameModal.loading}
        error={newGroupNameModal.error}
        onClose={() => setNewGroupNameModal({ open: false, group: null, loading: false, error: '' })}
        onSubmit={async ({ new_name }) => {
          setNewGroupNameModal(m => ({ ...m, loading: true, error: '' }));
          try {
            const groupId = newGroupNameModal.group.group_id || newGroupNameModal.group.id;
            if(!groupId){
              setNewGroupNameModal(m => ({ ...m, loading: false, error: "Group ID is missing." }));
              return;
            }
            const res = await fetch(`http://localhost:9090/groups/${groupId}/updatename`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`,
              },
              body: JSON.stringify({ new_name })
            })
            if (!res.ok) {
              const errData = await res.json().catch(() => ({}));
              throw new Error(errData.error || 'Failed to update group name');
            }
            setNewGroupNameModal({open: false, group: null, loading: false, error: ''})
            if(transactionType === 'group'){
              const res = await fetch(`http://localhost:9090/groups?page_id=1&page_size=5`, {
                headers: { 'Authorization': `Bearer ${accessToken}` },
              });
              const data = await res.json();
              setGroups(Array.isArray(data) ? data : []);
            }
          }catch(err){
            setNewGroupNameModal(m => ({ ...m, loading: false, error: err.message }));
          }
        }}
      />
      <LeaveGroupConfirmModal
        open={leaveGroupModal.open}
        group={leaveGroupModal.group || {}}
        loading={leaveGroupModal.loading}
        error={leaveGroupModal.error}
        onClose={() => setLeaveGroupModal({ open: false, group: null, loading: false, error: '' })}
        onConfirm={async () => {
          setLeaveGroupModal(m => ({ ...m, loading: true, error: '' }));
          try {
            const groupId = leaveGroupModal.group.group_id || leaveGroupModal.group.id;
            if(!groupId){
              setLeaveGroupModal(m => ({ ...m, loading: false, error: "Group ID is missing." }));
              return;
            }
            const res = await fetch(`http://localhost:9090/groups/${groupId}/leave`, {
              method: 'POST',
              headers: {
                'Authorization': `Bearer ${accessToken}`,
              },
            });
            if (!res.ok) {
              const errData = await res.json().catch(() => ({}));
              throw new Error(errData.error || 'Failed to leave group');
            }
            setLeaveGroupModal({ open: false, group: null, loading: false, error: '' });
            // Refresh groups list
            if (transactionType === 'group') {
              const res = await fetch(`http://localhost:9090/groups?page_id=1&page_size=5`, {
                headers: { 'Authorization': `Bearer ${accessToken}` },
              });
              const data = await res.json();
              setGroups(Array.isArray(data) ? data : []);
            }
          } catch (err) {
            setLeaveGroupModal(m => ({ ...m, loading: false, error: err.message }));
          }
        }}
      />
      <GroupHistoryModal
        open={groupHistoryModal.open}
        group={groupHistoryModal.group || {}}
        history={groupHistoryModal.history}
        loading={groupHistoryModal.loading}
        error={groupHistoryModal.error}
        onClose={() => setGroupHistoryModal({ open: false, group: null, history: [], loading: false, error: '' })}
      />
      <GroupSendMoneyModal
        open={groupSendMoneyModal.open}
        group={groupSendMoneyModal.group || {}}
        members={groupSendMoneyModal.members}
        loading={groupSendMoneyModal.loading}
        error={groupSendMoneyModal.error}
        onClose={() => setGroupSendMoneyModal({ open: false, group: null, members: [], loading: false, error: '' })}
        onSubmit={async ({ to_username, amount }) => {
          const groupId = groupSendMoneyModal.group.group_id || groupSendMoneyModal.group.id;
          try {
            const res = await fetch(`http://localhost:9090/groups/${groupId}/transaction`, {
              method: 'POST',
              headers: {
                'Authorization': `Bearer ${accessToken}`,
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                to_username,
                amount: Number(amount),
                currency: groupSendMoneyModal.group.currency,
                type: groupSendMoneyModal.group.type,
              }),
            });
            if (!res.ok) {
              const errorData = await res.json().catch(() => ({}));
              throw new Error(errorData.error || 'Failed to send money');
            }
            // Optionally refresh group history or UI here
          } catch (err) {
            throw err;
          }
        }}
      />
      <PromptSidebar
        open={promptSidebar.open}
        prompt={promptSidebar.prompt}
        loading={promptSidebar.loading}
        error={promptSidebar.error}
        mindset={mindset}
        setMindset={setMindset}
        onMindsetChange={async (newMindset) => {
          setPromptSidebar(ps => ({ ...ps, loading: true, error: '', prompt: '' }));
          try {
            const res = await fetch('http://localhost:9090/prompt', {
              method: 'POST',
              headers: { 'Authorization': `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
              body: JSON.stringify({ saving_mindset: newMindset })
            });
            if (!res.ok) {
              const errorData = await res.json().catch(() => ({}));
              throw new Error(errorData.error || 'Failed to get prompt');
            }
            const data = await res.json();
            setPromptSidebar(ps => ({ ...ps, loading: false, prompt: data.sentence || data.prompt || JSON.stringify(data), error: '' }));
          } catch (err) {
            setPromptSidebar(ps => ({ ...ps, loading: false, prompt: '', error: err.message }));
          }
        }}
        onClose={() => setPromptSidebar({ ...promptSidebar, open: false })}
      />
      {/* Floating AI Prompt Button */}
      <button
        className="fixed left-4 top-1/2 z-40 bg-gradient-to-br from-purple-600 via-blue-700 to-indigo-800 text-white rounded-full shadow-xl p-4 hover:scale-110 hover:shadow-2xl transition-all duration-200 border-2 border-purple-400"
        style={{ transform: 'translateY(-50%)' }}
        onClick={async () => {
          setPromptSidebar({ open: true, loading: true, prompt: '', error: '' });
          try {
            const res = await fetch('http://localhost:9090/prompt', {
              method: 'POST',
              headers: { 'Authorization': `Bearer ${accessToken}`, 'Content-Type': 'application/json' },
              body: JSON.stringify({ saving_mindset: mindset })
            });
            if (!res.ok) {
              const errorData = await res.json().catch(() => ({}));
              throw new Error(errorData.error || 'Failed to get prompt');
            }
            const data = await res.json();
            setPromptSidebar({ open: true, loading: false, prompt: data.sentence || data.prompt || JSON.stringify(data), error: '' });
          } catch (err) {
            setPromptSidebar({ open: true, loading: false, prompt: '', error: err.message });
          }
        }}
        aria-label="Open AI Prompt"
      >
        <span className="text-2xl">✨</span>
      </button>

      {loading ? (
        <div className="text-white text-center">Loading accounts...</div>
      ) : error ? (
        <div className="text-red-400 text-center">{error}</div>
      ) : (
        <div className="flex-1">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mx-auto max-w-7xl pt-4">
            {/* Conditional grid rendering based on transactionType */}
            {transactionType === 'group' ? (
              groups.length === 0 ? (
                <div className="col-span-full text-white text-center">No groups found. Create a group to get started.</div>
              ) : (
                groups.map((group, idx) => (
                  <GroupCard
                    key={group.group_id ? `group-${group.group_id}` : `group-${group.group_name}-${idx}`}
                    group={group}
                    accessToken={accessToken}
                    onAction={async (action, groupObj) => {
                      if (action === 'view-history') {
                        setGroupHistoryModal({ open: true, group: groupObj, history: [], loading: true, error: '' });
                        try {
                          const res = await fetch(`http://localhost:9090/groups/${groupObj.group_id || groupObj.id}/history`, {
                            headers: { 'Authorization': `Bearer ${accessToken}` },
                          });
                          if (!res.ok) {
                            const errorData = await res.json().catch(() => ({}));
                            throw new Error(errorData.error || 'Failed to fetch group history');
                          }
                          const data = await res.json();
                          setGroupHistoryModal({ open: true, group: groupObj, history: data, loading: false, error: '' });
                        } catch (err) {
                          setGroupHistoryModal({ open: true, group: groupObj, history: [], loading: false, error: err.message });
                        }
                        return;
                      }
                      if (action === 'send-money') {
                        setGroupSendMoneyModal({ open: true, group: groupObj, members: [], loading: true, error: '' });
                        // Fetch group members
                        try {
                          const res = await fetch(`http://localhost:9090/groups/${groupObj.group_id || groupObj.id}/accounts?page_id=1&page_size=5`, {
                            headers: { 'Authorization': `Bearer ${accessToken}` },
                          });
                          if (!res.ok) {
                            const errorData = await res.json().catch(() => ({}));
                            throw new Error(errorData.error || 'Failed to fetch group members');
                          }
                          const data = await res.json();
                          setGroupSendMoneyModal({ open: true, group: groupObj, members: data, loading: false, error: '' });
                        } catch (err) {
                          setGroupSendMoneyModal({ open: true, group: groupObj, members: [], loading: false, error: err.message });
                        }
                        return;
                      }
                      if (action === 'add-member') {
                        setAddMemberModal({ open: true, group: groupObj, loading: false, error: '' });
                      } else if (action === 'update-name') {
                        setNewGroupNameModal({ open: true, group: groupObj, loading: false, error: '' });
                      } else if (action === 'leave') {
                        setLeaveGroupModal({ open: true, group: groupObj, loading: false, error: '' });
                      } else if (action === 'delete') {
                        // Confirm, then POST to /groups/:id/delete
                      }
                    }}
                  />
                ))
              )
            ) : (
              Object.keys(userTransactions).length === 0 ? (
                <div className="col-span-full text-white text-center">No transactions found. Send or receive money to see transaction history.</div>
              ) : (
                Object.values(userTransactions)
                  .filter((user) => user.username !== username)
                  .map((user) => (
                    <UserTotalsCard
                      key={user.username}
                      username={user.username}
                      accessToken={accessToken}
                      onCardClick={async () => {
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
                      onSendMoney={() => setTransferModal({ open: true, toUsername: user.username, loading: false, error: '' })}
                    />
                  ))
              )
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
