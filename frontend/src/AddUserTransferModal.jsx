import { useState, useEffect } from 'react';

const SUPPORTED_CURRENCIES = ['USD', 'Euros', 'Rupees'];
const SUPPORTED_TYPES = ['rent', 'food', 'travel', 'savings', 'bills','medical','shopping','misc'];

export default function AddUserTransferModal({ accessToken, onTransferSuccess, toUsername }) {
  // Modal is now fully controlled by parent; always open if rendered
  // Remove internal open state

  const [step, setStep] = useState(toUsername ? 1 : 0); // skip username input if toUsername is provided
  const [username, setUsername] = useState(toUsername || '');
  const [userCheckLoading, setUserCheckLoading] = useState(false);
  const [userCheckError, setUserCheckError] = useState('');
  const [transferLoading, setTransferLoading] = useState(false);
  const [transferError, setTransferError] = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState(SUPPORTED_CURRENCIES[0]);
  const [type, setType] = useState(SUPPORTED_TYPES[0]);

  // Reset state when modal closes or toUsername changes
  useEffect(() => {
    setStep(toUsername ? 1 : 0);
    setUsername(toUsername || '');
    setUserCheckLoading(false);
    setUserCheckError('');
    setTransferLoading(false);
    setTransferError('');
    setSuccessMsg('');
    setAmount('');
    setCurrency(SUPPORTED_CURRENCIES[0]);
    setType(SUPPORTED_TYPES[0]);
  }, [toUsername]);

  // If toUsername is provided, skip user check
  const handleUserCheck = async (e) => {
    e.preventDefault();
    setUserCheckError('');
    setUserCheckLoading(true);
    setSuccessMsg('');
    try {
      // Use new public endpoint for recipient lookup
      const res = await fetch(`http://localhost:9090/users/${toUsername || username}`);
      if (!res.ok) {
        setUserCheckError('User not found.');
        setUserCheckLoading(false);
        return;
      }
      setUserCheckLoading(false);
      setStep(1);
      if (toUsername) return; // Already have username, skip further checks
    } catch (err) {
      setUserCheckError('Network error: ' + err.message);
      setUserCheckLoading(false);
    }
  };

  const handleTransfer = async (e) => {
    e.preventDefault();
    setTransferError('');
    setTransferLoading(true);
    setSuccessMsg('');
    if (!amount || isNaN(amount) || Number(amount) <= 0) {
      setTransferError('Please enter a valid amount.');
    return;
      setTransferLoading(false);
      return;
    }
    try {
      const res = await fetch('http://localhost:9090/transfers', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          to_username: toUsername || username,
          amount: Number(amount),
          currency,
          type,
        }),
      });
      if (!res.ok) {
        const err = await res.json();
        setTransferError(err.error || 'Transfer failed');
        setTransferLoading(false);
        return;
      }
      setSuccessMsg('Transfer successful!');
      setTransferLoading(false);
      setTimeout(() => {
        setSuccessMsg('');
        if (onTransferSuccess) onTransferSuccess();
        onClose();
      }, 1200);
      setAmount('');
      setCurrency(SUPPORTED_CURRENCIES[0]);
      setType(SUPPORTED_TYPES[0]);
    } catch (err) {
      setTransferError('Network error: ' + err.message);
      setTransferLoading(false);
    }
  };



  return (
    <>

      {/* Modal always open if rendered by parent */}
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
        <div className="bg-gray-900 rounded-xl shadow-lg p-8 max-w-md w-full border border-gray-700 relative">
          <button
              className="absolute top-2 right-2 text-gray-400 hover:text-white text-xl"
              onClick={() => onTransferSuccess && onTransferSuccess()}
              title="Close"
            >
              &times;
            </button>
            {step === 0 && !toUsername && (
              <form onSubmit={handleUserCheck} className="flex flex-col gap-4">
                <h2 className="text-xl font-bold mb-2 text-center text-white">Send Money</h2>
                <input
                  type="text"
                  className="p-2 rounded bg-gray-800 border border-gray-600 text-white"
                  placeholder="Recipient username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                  disabled={userCheckLoading}
                />
                <button
                  type="submit"
                  className="bg-blue-600 hover:bg-blue-700 text-white py-2 rounded font-semibold"
                  disabled={userCheckLoading}
                >
                  {userCheckLoading ? 'Checking...' : 'Next'}
                </button>
                {userCheckError && <div className="text-red-400 text-sm text-center">{userCheckError}</div>}
              </form>
            )}
            {step === 1 && (
              <form onSubmit={handleTransfer} className="flex flex-col gap-4">
                <h2 className="text-xl font-bold mb-2 text-center text-white">Send to <span className="text-blue-400">{username}</span></h2>
                <input
                  type="number"
                  min="1"
                  step="1"
                  placeholder="Amount"
                  value={amount}
                  onChange={e => setAmount(e.target.value)}
                  className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                  required
                />
                <select
                  value={currency}
                  onChange={e => setCurrency(e.target.value)}
                  className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                >
                  {SUPPORTED_CURRENCIES.map(cur => (
                    <option key={cur} value={cur}>{cur}</option>
                  ))}
                </select>
                <select
                  value={type}
                  onChange={e => setType(e.target.value)}
                  className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                >
                  {SUPPORTED_TYPES.map(t => (
                    <option key={t} value={t}>{t.charAt(0).toUpperCase() + t.slice(1)}</option>
                  ))}
                </select>
                {transferError && <div className="text-red-400 text-sm text-center">{transferError}</div>}
                <button
                  type="submit"
                  disabled={transferLoading}
                  className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-60"
                >
                  {transferLoading ? 'Sending...' : 'Send Money'}
                </button>
              </form>
            )}
            {successMsg && (
              <div className="text-green-400 text-center mt-4">{successMsg}</div>
            )}
          </div>
        </div>
    </>
  );
}
