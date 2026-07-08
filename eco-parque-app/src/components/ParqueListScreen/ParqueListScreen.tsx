import React from 'react';
import { View, Text, FlatList, ActivityIndicator, Alert, TouchableOpacity, RefreshControl, StatusBar } from 'react-native';
import { MaterialIcons, Feather } from '@expo/vector-icons';
import { useParques } from '../../hooks/useParques';
import { styles } from './ParqueListScreen.styles';
import { FadeSlideCard, PressableScale } from '../ui/AnimatedCards';

interface ParqueListScreenProps {
  onOpenMenu: () => void;
  onAdd: () => void;
  onSelect: (id: string) => void;
}

export function ParqueListScreen({ onOpenMenu, onAdd, onSelect }: ParqueListScreenProps) {
  const { parques, loading, error, refresh, removeParque } = useParques();

  const handleRemove = async (id: string) => {
    Alert.alert(
      'Eliminar Ecoparque',
      '¿Estás seguro de que deseas eliminar este parque?',
      [
        { text: 'Cancelar', style: 'cancel' },
        { 
          text: 'Eliminar', 
          style: 'destructive',
          onPress: async () => {
            const result = await removeParque(id);
            if (!result.success) {
              Alert.alert('Error', result.error);
            }
          }
        }
      ]
    );
  };

  if (loading && parques.length === 0) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#14B8A6" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Estás administrando</Text>
          <Text style={styles.title}>Mis Ecoparques</Text>
        </View>
        <TouchableOpacity style={styles.logoutButton} onPress={onOpenMenu} activeOpacity={0.7}>
          <MaterialIcons name="menu" size={24} color="#14B8A6" />
        </TouchableOpacity>
      </View>
      
      {error && (
        <View style={styles.errorContainer}>
          <MaterialIcons name="error-outline" size={20} color="#F43F5E" />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}


      <FlatList
        data={parques}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        showsVerticalScrollIndicator={false}
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refresh} tintColor="#14B8A6" />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIconCircle}>
               <MaterialIcons name="nature-people" size={48} color="#14B8A6" />
            </View>
            <Text style={styles.emptyText}>No tienes ecoparques todavía</Text>
            <Text style={styles.emptySubtext}>Toca el botón flotante para agregar tu primer parque a la lista.</Text>
          </View>
        }
        renderItem={({ item, index }) => {
          let lat = null;
          let lng = null;
          if (item.direccion) {
            const parts = item.direccion.split(',');
            if (parts.length === 2) {
              const parsedLat = parseFloat(parts[0].trim());
              const parsedLng = parseFloat(parts[1].trim());
              if (!isNaN(parsedLat) && !isNaN(parsedLng)) {
                lat = parsedLat;
                lng = parsedLng;
              }
            }
          }

          return (
            <FadeSlideCard delay={index * 50} style={{ marginBottom: 16 }}>
              <PressableScale onPress={() => onSelect(item.id)}>
                <View style={[styles.card, { position: 'relative', overflow: 'hidden' }]}>
                  {/* Hero Icon de fondo */}
                  <Feather 
                    name="map" 
                    size={100} 
                    color="rgba(20, 184, 166, 0.04)" 
                    style={{ position: 'absolute', right: -20, bottom: -10, transform: [{ rotate: '-10deg' }] }} 
                  />

                  <View style={styles.cardContent}>
                    <View style={styles.cardHeader}>
                      <View style={{ flexDirection: 'row', alignItems: 'center', gap: 6 }}>
                        <Feather name="map" size={14} color="#14B8A6" />
                        <Text style={styles.cardLabel}>Ecoparque</Text>
                      </View>
                      <View style={styles.badge}>
                        <MaterialIcons name="eco" size={14} color="#061114" />
                        <Text style={styles.badgeText}>Activo</Text>
                      </View>
                    </View>
                    <Text style={styles.cardTitle}>{item.nombre}</Text>
                    <Text style={styles.cardSubtitle}>
                      <Feather name="file-text" size={12} color="#8FA3A9" /> {item.descripcion || 'Sin descripción'}
                    </Text>
                    <Text style={styles.cardSubtitle}>
                      <Feather name="map-pin" size={12} color="#8FA3A9" /> {item.direccion || 'Sin dirección'}
                    </Text>
                  </View>
                  
                  <View style={styles.cardDivider} />
                
                  <View style={styles.cardFooter}>
                    <Text style={styles.statsText}>~10 hectáreas</Text>
                    <TouchableOpacity style={styles.deleteButton} onPress={() => handleRemove(item.id)} activeOpacity={0.6}>
                      <Text style={styles.deleteButtonText}>Eliminar</Text>
                      <Feather name="trash-2" size={16} color="#F43F5E" />
                    </TouchableOpacity>
                  </View>
                </View>
              </PressableScale>
            </FadeSlideCard>
          );
        }}
      />

      <TouchableOpacity style={styles.fab} onPress={onAdd} activeOpacity={0.9}>
        <MaterialIcons name="add" size={28} color="#061114" />
      </TouchableOpacity>
    </View>
  );
}


