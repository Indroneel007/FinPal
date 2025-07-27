import KebabMenu from './KebabMenu';
import { useState } from 'react';

export default function GroupCard({ group, accessToken, onAction }) {
  // Modal state for future: add member, update name, etc.
  // const [modal, setModal] = useState(null);

  const groupId = group.group_id || group.id;

  // Define kebab menu options
  const options = [
    {
      label: 'Add Member',
      onClick: () => onAction && onAction('add-member', group),
    },
    {
      label: 'Group Name',
      onClick: () => onAction && onAction('update-name', group),
    },
    {
      label: 'Leave Group',
      onClick: () => onAction && onAction('leave', group),
    },
    {
      label: 'Delete Group',
      onClick: () => onAction && onAction('delete', group),
    },
  ];

  return (
    <div className="bg-gray-800 rounded-xl shadow-lg p-6 flex flex-col gap-3 border border-gray-700 hover:border-blue-500 transition-colors relative">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-xl font-bold text-purple-300">{group.group_name}</h3>
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-400">{group.currency} | {group.type}</span>
          <KebabMenu options={options} />
        </div>
      </div>
    </div>
  );
}

