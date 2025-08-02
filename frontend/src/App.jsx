// NOTE: App must be wrapped in <BrowserRouter> in main.jsx for navigation (useNavigate, etc) to work.
import { useState } from 'react';
import LandingSection from './LandingSection';
import AuthCard from './AuthCard';
import Navbar from './Navbar';
import ForgotPassword from './ForgotPassword';
import OtpVerify from './OtpVerify';
import { Routes, Route } from 'react-router-dom';
import ResetPassword from './ResetPassword';
import Location from './Location';
import MainPage from './MainPage';

function App() {
  const [authMode, setAuthMode] = useState(null);
  const [username, setUsername] = useState(null);

  return (
    <div className="min-h-screen w-full relative overflow-x-hidden">
      <Routes>
        <Route
          path="/"
          element={
            <>
              <Navbar onLogin={() => setAuthMode('login')} username={username} showLogin={true} />
              <LandingSection />
              {authMode && (
                <AuthCard
                  mode={authMode}
                  onClose={() => setAuthMode(null)}
                  onSwitch={() => setAuthMode(authMode === 'login' ? 'register' : 'login')}
                  setUsername={setUsername}
                />
              )}
            </>
          }
        />
        <Route path="/forgot-password" element={<><Navbar username={username} showLogin={false} /><ForgotPassword /></>} />
        <Route path="/otp-verify" element={<><Navbar username={username} showLogin={false} /><OtpVerify /></>} />
        <Route path="/reset-password" element={<><Navbar username={username} showLogin={false} /><ResetPassword /></>} />
        <Route path="/location" element={<><Navbar username={username} showLogin={false} /><Location /></>} />
        <Route path="/main" element={<MainPage />} />
      </Routes>
    </div>
  );
}

export default App;
