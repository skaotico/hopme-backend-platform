import React, { useState } from 'react';
import { View, ActivityIndicator, StyleSheet } from 'react-native';
import { useAuth } from './src/hooks/useAuth';
import { LoginScreen } from './src/components/LoginScreen';
import { RegisterScreen } from './src/components/RegisterScreen';
import { ParqueListScreen } from './src/components/ParqueListScreen';
import { ParqueAddScreen } from './src/components/ParqueAddScreen';
import { ZonaListScreen } from './src/components/ZonaListScreen';
import { ZonaAddScreen } from './src/components/ZonaAddScreen';
import { ArbolListScreen } from './src/components/ArbolListScreen';
import { ArbolAddScreen } from './src/components/ArbolAddScreen';

type ScreenName = 'ParqueList' | 'ParqueAdd' | 'ZonaList' | 'ZonaAdd' | 'ArbolList' | 'ArbolAdd';

export default function App() {
  const { user, loading, login, register, logout } = useAuth();
  const [isRegistering, setIsRegistering] = useState(false);
  
  // Custom router state
  const [currentScreen, setCurrentScreen] = useState<ScreenName>('ParqueList');
  const [selectedParqueId, setSelectedParqueId] = useState<string | null>(null);
  const [selectedZonaId, setSelectedZonaId] = useState<string | null>(null);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#14B8A6" />
      </View>
    );
  }

  const renderAuthenticatedApp = () => {
    switch (currentScreen) {
      case 'ParqueAdd':
        return (
          <ParqueAddScreen 
            onBack={() => setCurrentScreen('ParqueList')} 
            onSuccess={() => setCurrentScreen('ParqueList')} 
          />
        );
      case 'ZonaList':
        if (!selectedParqueId) {
          setCurrentScreen('ParqueList');
          return null;
        }
        return (
          <ZonaListScreen 
            parqueId={selectedParqueId} 
            onBack={() => setCurrentScreen('ParqueList')} 
            onAdd={() => setCurrentScreen('ZonaAdd')} 
            onSelect={(zonaId) => {
              setSelectedZonaId(zonaId);
              setCurrentScreen('ArbolList');
            }}
          />
        );
      case 'ZonaAdd':
        if (!selectedParqueId) {
          setCurrentScreen('ParqueList');
          return null;
        }
        return (
          <ZonaAddScreen 
            parqueId={selectedParqueId} 
            onBack={() => setCurrentScreen('ZonaList')} 
            onSuccess={() => setCurrentScreen('ZonaList')} 
          />
        );
      case 'ArbolList':
        if (!selectedZonaId) {
          setCurrentScreen('ZonaList');
          return null;
        }
        return (
          <ArbolListScreen 
            zonaId={selectedZonaId} 
            onBack={() => setCurrentScreen('ZonaList')} 
            onAdd={() => setCurrentScreen('ArbolAdd')} 
          />
        );
      case 'ArbolAdd':
        if (!selectedZonaId) {
          setCurrentScreen('ZonaList');
          return null;
        }
        return (
          <ArbolAddScreen 
            zonaId={selectedZonaId} 
            onBack={() => setCurrentScreen('ArbolList')} 
            onSuccess={() => setCurrentScreen('ArbolList')} 
          />
        );
      case 'ParqueList':
      default:
        return (
          <ParqueListScreen 
            onLogout={logout} 
            onAdd={() => setCurrentScreen('ParqueAdd')}
            onSelect={(id) => {
              setSelectedParqueId(id);
              setCurrentScreen('ZonaList');
            }}
          />
        );
    }
  };

  return (
    <View style={styles.container}>
      {user ? (
        renderAuthenticatedApp()
      ) : isRegistering ? (
        <RegisterScreen 
          onRegisterSuccess={() => setIsRegistering(false)} 
          onGoToLogin={() => setIsRegistering(false)} 
          register={register} 
          loading={loading} 
        />
      ) : (
        <LoginScreen 
          onLoginSuccess={() => {}} 
          onGoToRegister={() => setIsRegistering(true)}
          login={login} 
          loading={loading} 
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#061114',
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#061114',
  },
});
