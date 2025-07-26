export default function GroupCard({ group, onSendMoney, onAddUser }) {
  return (
    <div className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-3 border border-gray-700 hover:border-blue-500 transition-colors">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-xl font-bold text-purple-300">{group.group_name}</h3>
        <span className="text-sm text-gray-400">{group.currency} | {group.type}</span>
      </div>
      <div className="flex gap-2 mt-2">
        <button
          className="flex-1 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-bold"
          onClick={onSendMoney}
        >
          Send Money
        </button>
        <button
          className="flex-1 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors font-bold"
          onClick={onAddUser}
        >
          Add User
        </button>
      </div>
    </div>
  );
}
