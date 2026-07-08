import React from 'react';
import { View, Text, FlatList, ActivityIndicator, Alert, TouchableOpacity, RefreshControl, Platform, StatusBar } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useParques } from '../../hooks/useParques';
import { styles } from './ParqueListScreen.styles';

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
        renderItem={({ item }) => {
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
            <TouchableOpacity style={styles.card} activeOpacity={0.8} onPress={() => onSelect(item.id)}>
              <View style={styles.cardContent}>
                <View style={styles.cardHeader}>
                  <Text style={styles.cardLabel}>Parque registrado</Text>
                  <View style={styles.badge}>
                    <MaterialIcons name="eco" size={14} color="#061114" />
                    <Text style={styles.badgeText}>Activo</Text>
                  </View>
                </View>
                <Text style={styles.cardTitle}>{item.nombre}</Text>
                <Text style={styles.cardSubtitle}>
                  <MaterialIcons name="description" size={14} color="#8FA3A9" /> {item.descripcion || 'Sin descripción'}
                </Text>
                <Text style={styles.cardSubtitle}>
                  <MaterialIcons name="place" size={14} color="#8FA3A9" /> {item.direccion || 'Sin dirección'}
                </Text>
              </View>
              
              {/* Mapa temporalmente deshabilitado */}
              
              <View style={styles.cardDivider} />
            
            <View style={styles.cardFooter}>
               <Text style={styles.statsText}>~10 hectáreas</Text>
               <TouchableOpacity style={styles.deleteButton} onPress={() => handleRemove(item.id)} activeOpacity={0.6}>
                 <Text style={styles.deleteButtonText}>Eliminar</Text>
                 <MaterialIcons name="delete-outline" size={18} color="#F43F5E" />
               </TouchableOpacity>
            </View>
          </TouchableOpacity>
          );
        }}
      />

      <TouchableOpacity style={styles.fab} onPress={onAdd} activeOpacity={0.9}>
        <MaterialIcons name="add" size={28} color="#061114" />
      </TouchableOpacity>
    </View>
  );
}


