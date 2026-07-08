import React, { useState } from 'react';
import { View, ActivityIndicator, StyleSheet } from 'react-native';
import { useAuth } from './src/hooks/useAuth';
import { LoginScreen } from './src/components/LoginScreen/LoginScreen';
import { RegisterScreen } from './src/components/RegisterScreen/RegisterScreen';
import { ParqueListScreen } from './src/components/ParqueListScreen/ParqueListScreen';
import { ParqueAddScreen } from './src/components/ParqueAddScreen/ParqueAddScreen';
import { ZonaListScreen } from './src/components/ZonaListScreen/ZonaListScreen';
import { ZonaAddScreen } from './src/components/ZonaAddScreen/ZonaAddScreen';
import { ZonaEditScreen } from './src/components/ZonaEditScreen/ZonaEditScreen';
import { ArbolListScreen } from './src/components/ArbolListScreen/ArbolListScreen';
import { ArbolAddScreen } from './src/components/ArbolAddScreen/ArbolAddScreen';
import { ArbolEditScreen } from './src/components/ArbolEditScreen/ArbolEditScreen';
import { SensorListScreen } from './src/components/SensorListScreen/SensorListScreen';
import { NotificationListScreen } from './src/components/NotificationListScreen/NotificationListScreen';
import { Sidebar } from './src/components/Sidebar/Sidebar';
import { CatalogoListScreen } from './src/components/catalogo/CatalogoListScreen';
import { CatalogoFormScreen } from './src/components/catalogo/CatalogoFormScreen';
import { ErrorBoundary } from './src/components/ErrorBoundary/ErrorBoundary';
import { Zona } from './src/dto/zona.dto';
import { Arbol } from './src/dto/arbol.dto';
import { CatalogItem, CatalogType } from './src/dto/catalogo.dto';

type ScreenName =
  | 'ParqueList'
  | 'ParqueAdd'
  | 'ZonaList'
  | 'ZonaAdd'
  | 'ZonaEdit'
  | 'ArbolList'
  | 'ArbolAdd'
  | 'ArbolEdit'
  | 'SensorList'
  | 'CatalogoList_especies'
  | 'CatalogoList_estados-arbol'
  | 'CatalogoList_estados-estanque'
  | 'CatalogoList_estados-agua'
  | 'CatalogoList_tipos-sensor'
  | 'CatalogoAdd'
  | 'CatalogoEdit'
  | 'Notifications';

export default function App() {
  const { user, loading, login, register, logout } = useAuth();
  const [isRegistering, setIsRegistering] = useState(false);

  // Custom router state
  const [currentScreen, setCurrentScreen] = useState<ScreenName>('ParqueList');
  const [selectedParqueId, setSelectedParqueId] = useState<string | null>(null);
  const [selectedZonaId, setSelectedZonaId] = useState<string | null>(null);
  const [selectedArbolId, setSelectedArbolId] = useState<string | null>(null);
  const [selectedArbolCodigo, setSelectedArbolCodigo] = useState<string | undefined>(undefined);
  const [zonaToEdit, setZonaToEdit] = useState<Zona | null>(null);
  const [arbolToEdit, setArbolToEdit] = useState<Arbol | null>(null);

  // Sidebar and Catalog navigation state
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [activeCatalogType, setActiveCatalogType] = useState<CatalogType>('especies');
  const [catalogItemToEdit, setCatalogItemToEdit] = useState<CatalogItem | null>(null);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#14B8A6" />
      </View>
    );
  }

  const handleNavigateFromSidebar = (screen: string, catalogType?: CatalogType) => {
    if (screen === 'CatalogoList' && catalogType) {
      setActiveCatalogType(catalogType);
      setCurrentScreen(`CatalogoList_${catalogType}` as ScreenName);
    } else {
      setCurrentScreen(screen as ScreenName);
    }
  };

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
            onEdit={(zona) => {
              setZonaToEdit(zona);
              setCurrentScreen('ZonaEdit');
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
      case 'ZonaEdit':
        if (!selectedParqueId || !zonaToEdit) {
          setCurrentScreen('ZonaList');
          return null;
        }
        return (
          <ZonaEditScreen
            parqueId={selectedParqueId}
            zona={zonaToEdit}
            onBack={() => setCurrentScreen('ZonaList')}
            onSuccess={() => {
              setZonaToEdit(null);
              setCurrentScreen('ZonaList');
            }}
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
            onEdit={(arbol) => {
              setArbolToEdit(arbol);
              setCurrentScreen('ArbolEdit');
            }}
            onSensores={(arbol) => {
              setSelectedArbolId(arbol.id);
              setSelectedArbolCodigo(arbol.codigo);
              setCurrentScreen('SensorList');
            }}
          />
        );

      case 'SensorList':
        if (!selectedArbolId) {
          setCurrentScreen('ArbolList');
          return null;
        }
        return (
          <SensorListScreen
            arbolId={selectedArbolId}
            arbolCodigo={selectedArbolCodigo}
            onBack={() => setCurrentScreen('ArbolList')}
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
      case 'ArbolEdit':
        if (!selectedZonaId || !arbolToEdit) {
          setCurrentScreen('ArbolList');
          return null;
        }
        return (
          <ArbolEditScreen
            zonaId={selectedZonaId}
            arbol={arbolToEdit}
            onBack={() => {
              setArbolToEdit(null);
              setCurrentScreen('ArbolList');
            }}
            onSuccess={() => {
              setArbolToEdit(null);
              setCurrentScreen('ArbolList');
            }}
          />
        );

      // Catalogs mapping
      case 'CatalogoList_especies':
      case 'CatalogoList_estados-arbol':
      case 'CatalogoList_estados-estanque':
      case 'CatalogoList_estados-agua':
      case 'CatalogoList_tipos-sensor':
        return (
          <CatalogoListScreen
            type={activeCatalogType}
            onOpenMenu={() => setIsSidebarOpen(true)}
            onAdd={() => setCurrentScreen('CatalogoAdd')}
            onEdit={(item) => {
              setCatalogItemToEdit(item);
              setCurrentScreen('CatalogoEdit');
            }}
          />
        );

      case 'CatalogoAdd':
        return (
          <CatalogoFormScreen
            type={activeCatalogType}
            onBack={() => setCurrentScreen(`CatalogoList_${activeCatalogType}` as ScreenName)}
            onSuccess={() => setCurrentScreen(`CatalogoList_${activeCatalogType}` as ScreenName)}
          />
        );

      case 'CatalogoEdit':
        if (!catalogItemToEdit) {
          setCurrentScreen(`CatalogoList_${activeCatalogType}` as ScreenName);
          return null;
        }
        return (
          <CatalogoFormScreen
            type={activeCatalogType}
            item={catalogItemToEdit}
            onBack={() => {
              setCatalogItemToEdit(null);
              setCurrentScreen(`CatalogoList_${activeCatalogType}` as ScreenName);
            }}
            onSuccess={() => {
              setCatalogItemToEdit(null);
              setCurrentScreen(`CatalogoList_${activeCatalogType}` as ScreenName);
            }}
          />
        );

      case 'Notifications':
        return (
          <NotificationListScreen 
            onBack={() => setCurrentScreen('ParqueList')} 
          />
        );

      case 'ParqueList':
      default:
        return (
          <ParqueListScreen
            onOpenMenu={() => setIsSidebarOpen(true)}
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
      {user && (
        <Sidebar
          isOpen={isSidebarOpen}
          onClose={() => setIsSidebarOpen(false)}
          currentScreen={currentScreen}
          onNavigate={handleNavigateFromSidebar}
          onLogout={logout}
        />
      )}
      {user ? (
        <ErrorBoundary>
          {renderAuthenticatedApp()}
        </ErrorBoundary>
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
