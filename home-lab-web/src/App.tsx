import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { Layout } from './components/layout/Layout';
import { LoginScreen } from './components/auth/LoginScreen';
import { RegisterScreen } from './components/auth/RegisterScreen';
import { ParqueListScreen } from './components/parques/ParqueListScreen';

const NotificationsPlaceholder = () => <div className="screen-container"><h1 className="screen-header">Notificaciones</h1></div>;
const CatalogsPlaceholder = () => <div className="screen-container"><h1 className="screen-header">Catálogos</h1></div>;

const AuthRoute = ({ children }: { children: React.ReactNode }) => {
  const { user, loading } = useAuth();
  const [isRegistering, setIsRegistering] = useState(false);
  
  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', backgroundColor: 'var(--bg-primary)' }}>
        <div className="animate-spin" style={{ width: '40px', height: '40px', border: '4px solid var(--accent-light)', borderTopColor: 'var(--accent-primary)', borderRadius: '50%' }} />
      </div>
    );
  }
  
  if (user) {
    return <>{children}</>;
  }
  
  return isRegistering ? (
    <RegisterScreen onGoToLogin={() => setIsRegistering(false)} />
  ) : (
    <LoginScreen onGoToRegister={() => setIsRegistering(true)} />
  );
};

export const App = () => {
  return (
    <AuthProvider>
      <Router>
        <AuthRoute>
          <Routes>
            <Route path="/" element={<Layout />}>
              <Route index element={<ParqueListScreen />} />
              <Route path="notifications" element={<NotificationsPlaceholder />} />
              <Route path="catalogos/:type" element={<CatalogsPlaceholder />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Routes>
        </AuthRoute>
      </Router>
    </AuthProvider>
  );
};

export default App;
