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

function App() {
  const [authMode, setAuthMode] = useState(null);

  return (
    <div className="min-h-screen w-full relative overflow-x-hidden">
      <Navbar onLogin={() => setAuthMode('login')} />
      <Routes>
        <Route
          path="/"
          element={
            <>
              <LandingSection />
              {authMode && (
                <AuthCard
                  mode={authMode}
                  onClose={() => setAuthMode(null)}
                  onSwitch={() => setAuthMode(authMode === 'login' ? 'register' : 'login')}
                />
              )}
            </>
          }
        />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/otp-verify" element={<OtpVerify />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route path="/location" element={<Location />} />
      </Routes>
    </div>
  );
}

export default App;
