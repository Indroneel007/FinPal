import { useState } from 'react';

export default function TransactionTypeSelector({ value, onChange }) {
  return (
    <div className="flex flex-col items-end gap-4">
      <div className="flex gap-2 bg-gray-800 rounded-lg p-2 shadow-md">
        <button
          className={`px-4 py-2 rounded-lg font-semibold transition-colors ${value === 'user' ? 'bg-blue-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-blue-700'}`}
          onClick={() => onChange('user')}
        >
          Single User
        </button>
        <button
          className={`px-4 py-2 rounded-lg font-semibold transition-colors ${value === 'group' ? 'bg-green-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-green-700'}`}
          onClick={() => onChange('group')}
        >
          Group
        </button>
      </div>
    </div>
  );
}
