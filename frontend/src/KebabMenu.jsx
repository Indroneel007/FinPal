import { useState, useRef, useEffect } from 'react';

export default function KebabMenu({ options = [] }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    function handleClickOutside(event) {
      if (ref.current && !ref.current.contains(event.target)) {
        setOpen(false);
      }
    }
    if (open) {
      document.addEventListener('mousedown', handleClickOutside);
    } else {
      document.removeEventListener('mousedown', handleClickOutside);
    }
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [open]);

  return (
    <div className="relative" ref={ref}>
      <button
        className="p-2 rounded-full hover:bg-gray-700 focus:outline-none"
        onClick={() => setOpen((v) => !v)}
        aria-label="Open menu"
      >
        <svg width="24" height="24" fill="none" viewBox="0 0 24 24">
          <circle cx="12" cy="5" r="1.5" fill="#aaa" />
          <circle cx="12" cy="12" r="1.5" fill="#aaa" />
          <circle cx="12" cy="19" r="1.5" fill="#aaa" />
        </svg>
      </button>
      {open && (
        <div className="absolute right-0 mt-2 w-48 bg-gray-900 border border-gray-700 rounded-lg shadow-xl z-50">
          {options.map((option, idx) => (
            <button
              key={option.label}
              className="block w-full text-left px-4 py-2 text-sm text-gray-200 hover:bg-gray-800 hover:text-blue-400 transition-colors"
              onClick={() => {
                setOpen(false);
                option.onClick();
              }}
            >
              {option.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
