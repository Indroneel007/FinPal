import React from 'react';

export default function PromptSidebar({ open, onClose, prompt, loading, error, mindset, setMindset, onMindsetChange }) {
  return (
    <div
      className={`fixed top-0 left-0 h-full bg-gradient-to-br from-purple-800 via-blue-900 to-gray-900 text-white shadow-2xl z-50 transition-transform duration-300 sidebar-scroll ${open ? 'translate-x-0' : '-translate-x-full'}`}

      style={{ width: '40vw', minWidth: 320, maxWidth: '600px', left: 0, top: 0 }}
    >
      <div className="flex flex-col h-full">
        <div className="flex items-center justify-between p-4 border-b border-purple-700">
          <h2 className="text-xl font-bold text-purple-200">✨ AI Prompt</h2>
          <button
            className="text-purple-400 hover:text-red-400 text-xl font-bold ml-2"
            onClick={onClose}
            aria-label="Close prompt sidebar"
          >
            ×
          </button>
        </div>
        {/* Mindset Dropdown, placed below header for no overlap */}
        <div className="flex items-center gap-2 bg-gray-900 bg-opacity-80 px-4 py-2 border-b border-purple-700" style={{zIndex: 10}}>
          <label htmlFor="mindset-select-sidebar" className="text-purple-200 font-semibold text-sm">Saving Mindset:</label>
          <select
            id="mindset-select-sidebar"
            value={mindset}
            onChange={e => {
              setMindset(e.target.value);
              if (onMindsetChange) onMindsetChange(e.target.value);
            }}
            className="rounded bg-gray-800 text-white border border-purple-400 px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
            style={{ minWidth: 80 }}
          >
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
          </select>
        </div>
        <div className="flex-1 p-6 flex flex-col items-center justify-center">
          {loading ? (
            <div className="animate-pulse text-lg text-blue-200">Thinking...</div>
          ) : error ? (
            <div className="text-red-400">{error}</div>
          ) : prompt ? (
            <div className="text-base md:text-lg text-center font-small text-blue-100 drop-shadow animate-fade-in" style={{ wordBreak: 'break-word', lineHeight: '1.5' }}>
              "{prompt}"
            </div>
          ) : (
            <div className="text-gray-400 text-lg">Click the button to get a prompt!</div>
          )}
        </div>
        <div className="p-4 border-t border-purple-700 text-center text-xs text-purple-300">
          Powered by FinPal AI
        </div>
      </div>
    </div>
  );
}
