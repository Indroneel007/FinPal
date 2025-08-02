import { useState, useRef } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { MapContainer, TileLayer, Marker, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import Navbar from './Navbar';

function ChangeView({ center }) {
  const map = useMap();
  map.setView(center);
  return null;
}

export default function Location() {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selected, setSelected] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const accessToken = location.state?.access_token;
  // Username from registration (passed via navigation state)
  const username = location.state?.username;

  // Guard: If missing token or username, show error and do not render page
  if (!accessToken || !username) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#18181b]">
        <div className="bg-gray-900 p-8 rounded-2xl shadow-lg max-w-md border border-gray-700 text-center text-white">
          <h2 className="text-2xl font-bold mb-4">Invalid Access</h2>
          <p className="mb-4">Please register again to access this page.</p>
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

  // Default map center (India)
  const defaultCenter = [22.5937, 78.9629];

  const handleInput = async (e) => {
    const value = e.target.value;
    setQuery(value);
    setSuggestions([]);
    setSelected(null);
    if (value.length < 3) return;
    const res = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(value)}`);
    const data = await res.json();
    setSuggestions(data);
  };

  const handleSelect = (s) => {
    setQuery(s.display_name);
    setSelected({
      address: s.display_name,
      lat: parseFloat(s.lat),
      lng: parseFloat(s.lon),
    });
    setSuggestions([]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (!selected || !selected.address) {
      setError('Please select a location from suggestions.');
      return;
    }
    if (!accessToken) {
      setError('Missing access token. Please register again.');
      return;
    }
    setLoading(true);
    try {
      const res = await fetch('https://finpal-1.onrender.com/location', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify({
          address: selected.address,
          lattitude: selected.lat,
          longitude: selected.lng
        })
      });
      if (!res.ok) {
        const err = await res.json();
        setError(err.error || 'Failed to save location');
        setLoading(false);
        return;
      }
      // Success: redirect to main page
      navigate('/main', { state: { access_token: accessToken, username } });
    } catch (err) {
      setError('Network error: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#18181b]">
      <Navbar username={username} />
      <div className="bg-gray-900 p-8 rounded-2xl shadow-lg w-full max-w-md border border-gray-700 mb-8">
        <h2 className="text-2xl font-bold text-white mb-6 text-center">Select Your Location</h2>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <input
            type="text"
            value={query}
            onChange={handleInput}
            placeholder="Search for a location"
            className="px-4 py-3 rounded-lg bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
            autoComplete="off"
            required
          />
          {suggestions.length > 0 && (
            <ul className="bg-gray-800 border border-gray-700 rounded-lg max-h-40 overflow-y-auto">
              {suggestions.map((s, i) => (
                <li
                  key={s.place_id}
                  className="px-4 py-2 cursor-pointer hover:bg-violet-700 text-white"
                  onClick={() => handleSelect(s)}
                >
                  {s.display_name}
                </li>
              ))}
            </ul>
          )}
          {error && <div className="text-red-400 text-sm text-center">{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className="text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:bg-gradient-to-l focus:ring-4 focus:outline-none focus:ring-purple-200 dark:focus:ring-purple-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-60"
          >
            {loading ? 'Saving...' : 'Submit Location'}
          </button>
        </form>
      </div>
      <div className="w-full max-w-xl h-80 rounded-2xl overflow-hidden border border-gray-700">
        <MapContainer
          center={selected ? [selected.lat, selected.lng] : defaultCenter}
          zoom={selected ? 13 : 4}
          scrollWheelZoom={true}
          style={{ height: '100%', width: '100%' }}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          {selected && <Marker position={[selected.lat, selected.lng]} />} 
          <ChangeView center={selected ? [selected.lat, selected.lng] : defaultCenter} />
        </MapContainer>
      </div>
    </div>
  );
}
