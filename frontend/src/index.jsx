import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import App from './App';
import ForgotPassword from './ForgotPassword';
import OtpVerify from './OtpVerify';
import ResetPassword from './ResetPassword';
import Location from './Location';
import './index.css';

const root = createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/otp-verify" element={<OtpVerify />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route path="/location" element={<Location />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);