import React, { useEffect, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Animated,
  Dimensions,
  TouchableWithoutFeedback,
  Platform,
  StatusBar,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { CatalogType, CATALOGS_CONFIG } from '../../dto/catalogo.dto';

const { width } = Dimensions.get('window');
const DRAWER_WIDTH = width * 0.78;

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
  currentScreen: string;
  onNavigate: (screen: string, catalogType?: CatalogType) => void;
  onLogout: () => void;
}

export function Sidebar({ isOpen, onClose, currentScreen, onNavigate, onLogout }: SidebarProps) {
  const slideAnim = useRef(new Animated.Value(-DRAWER_WIDTH)).current;
  const opacityAnim = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    if (isOpen) {
      // open animation
      Animated.parallel([
        Animated.timing(slideAnim, {
          toValue: 0,
          duration: 300,
          useNativeDriver: true,
        }),
        Animated.timing(opacityAnim, {
          toValue: 0.6,
          duration: 300,
          useNativeDriver: true,
        }),
      ]).start();
    } else {
      // close animation
      Animated.parallel([
        Animated.timing(slideAnim, {
          toValue: -DRAWER_WIDTH,
          duration: 250,
          useNativeDriver: true,
        }),
        Animated.timing(opacityAnim, {
          toValue: 0,
          duration: 250,
          useNativeDriver: true,
        }),
      ]).start();
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleLinkPress = (screen: string, catalogType?: CatalogType) => {
    onClose();
    setTimeout(() => {
      onNavigate(screen, catalogType);
    }, 100);
  };

  const isSelected = (screen: string, catalogType?: CatalogType) => {
    if (screen === 'CatalogoList' && catalogType) {
      return currentScreen === 'CatalogoList' && catalogType === catalogType;
    }
    return currentScreen === screen;
  };

  return (
    <View style={styles.container}>
      {/* Dark overlay backdrop */}
      <TouchableWithoutFeedback onPress={onClose}>
        <Animated.View style={[styles.overlay, { opacity: opacityAnim }]} />
      </TouchableWithoutFeedback>

      {/* Slide-out Drawer Panel */}
      <Animated.View
        style={[
          styles.drawer,
          {
            transform: [{ translateX: slideAnim }],
          },
        ]}
      >
        {/* Profile/Header Section */}
        <View style={styles.header}>
          <View style={styles.avatarContainer}>
            <MaterialIcons name="eco" size={32} color="#14B8A6" />
          </View>
          <Text style={styles.appName}>EcoParque Admin</Text>
          <Text style={styles.appVersion}>v2.0 • Plataforma</Text>
        </View>

        {/* Navigation Items */}
        <View style={styles.navSection}>
          <Text style={styles.sectionTitle}>ADMINISTRACIÓN</Text>
          <TouchableOpacity
            style={[
              styles.navItem,
              isSelected('ParqueList') && styles.navItemActive,
            ]}
            onPress={() => handleLinkPress('ParqueList')}
            activeOpacity={0.7}
          >
            <MaterialIcons
              name="nature-people"
              size={22}
              color={isSelected('ParqueList') ? '#14B8A6' : '#8FA3A9'}
            />
            <Text style={[styles.navText, isSelected('ParqueList') && styles.navTextActive]}>
              Mis Ecoparques
            </Text>
          </TouchableOpacity>

          <TouchableOpacity
            style={[
              styles.navItem,
              isSelected('Notifications') && styles.navItemActive,
            ]}
            onPress={() => handleLinkPress('Notifications')}
            activeOpacity={0.7}
          >
            <MaterialIcons
              name="notifications-none"
              size={22}
              color={isSelected('Notifications') ? '#14B8A6' : '#8FA3A9'}
            />
            <Text style={[styles.navText, isSelected('Notifications') && styles.navTextActive]}>
              Notificaciones
            </Text>
          </TouchableOpacity>
        </View>

        <View style={[styles.navSection, { flex: 1 }]}>
          <Text style={styles.sectionTitle}>MANTENIMIENTO DE CATÁLOGOS</Text>

          {Object.values(CATALOGS_CONFIG).map((cat) => {
            const active = currentScreen === 'CatalogoList' && cat.type === cat.type; 
            // Wait, to distinguish between different catalogs we check both screen name and the specific type.
            const isCurrentCatalog = isSelected('CatalogoList') && currentScreen === 'CatalogoList';
            // Let's pass type directly down
            return (
              <TouchableOpacity
                key={cat.type}
                style={[
                  styles.navItem,
                  currentScreen === `CatalogoList_${cat.type}` && styles.navItemActive,
                ]}
                onPress={() => handleLinkPress('CatalogoList', cat.type)}
                activeOpacity={0.7}
              >
                <MaterialIcons
                  name={cat.icon as any}
                  size={22}
                  color={currentScreen === `CatalogoList_${cat.type}` ? '#14B8A6' : '#8FA3A9'}
                />
                <Text
                  style={[
                    styles.navText,
                    currentScreen === `CatalogoList_${cat.type}` && styles.navTextActive,
                  ]}
                >
                  {cat.displayName}
                </Text>
              </TouchableOpacity>
            );
          })}
        </View>

        {/* Footer/Logout Action */}
        <View style={styles.footer}>
          <TouchableOpacity style={styles.logoutButton} onPress={onLogout} activeOpacity={0.7}>
            <MaterialIcons name="logout" size={20} color="#F43F5E" />
            <Text style={styles.logoutText}>Cerrar Sesión</Text>
          </TouchableOpacity>
        </View>
      </Animated.View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    ...StyleSheet.absoluteFill,
    zIndex: 9999,
  },
  overlay: {
    ...StyleSheet.absoluteFill,
    backgroundColor: '#000000',
  },
  drawer: {
    position: 'absolute',
    left: 0,
    top: 0,
    bottom: 0,
    width: DRAWER_WIDTH,
    backgroundColor: '#061114',
    borderRightWidth: 1,
    borderColor: '#1D343B',
    paddingTop: Platform.OS === 'ios' ? 60 : (StatusBar.currentHeight || 24) + 20,
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    paddingHorizontal: 24,
    paddingBottom: 24,
    borderBottomWidth: 1,
    borderBottomColor: '#112226',
    marginBottom: 20,
  },
  avatarContainer: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 16,
  },
  appName: {
    fontSize: 20,
    fontWeight: '800',
    color: '#FFFFFF',
    letterSpacing: -0.3,
  },
  appVersion: {
    fontSize: 12,
    color: '#8FA3A9',
    fontWeight: '500',
    marginTop: 4,
  },
  navSection: {
    paddingHorizontal: 16,
    marginBottom: 24,
  },
  sectionTitle: {
    fontSize: 11,
    fontWeight: '700',
    color: '#8FA3A9',
    letterSpacing: 1.2,
    marginBottom: 12,
    paddingHorizontal: 8,
  },
  navItem: {
    flexDirection: 'row',
    alignItems: 'center',
    height: 48,
    borderRadius: 12,
    paddingHorizontal: 12,
    marginBottom: 4,
  },
  navItemActive: {
    backgroundColor: 'rgba(20, 184, 166, 0.1)',
    borderWidth: 1,
    borderColor: 'rgba(20, 184, 166, 0.2)',
  },
  navText: {
    fontSize: 15,
    fontWeight: '600',
    color: '#8FA3A9',
    marginLeft: 12,
  },
  navTextActive: {
    color: '#14B8A6',
  },
  footer: {
    padding: 24,
    borderTopWidth: 1,
    borderTopColor: '#112226',
  },
  logoutButton: {
    flexDirection: 'row',
    alignItems: 'center',
    height: 48,
    borderRadius: 12,
    backgroundColor: 'rgba(244, 63, 94, 0.1)',
    borderWidth: 1,
    borderColor: 'rgba(244, 63, 94, 0.2)',
    justifyContent: 'center',
  },
  logoutText: {
    color: '#F43F5E',
    fontSize: 15,
    fontWeight: '700',
    marginLeft: 8,
  },
});
