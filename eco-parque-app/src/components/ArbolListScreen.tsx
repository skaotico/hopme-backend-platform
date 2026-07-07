import React, { useEffect } from 'react';
import { View, Text, FlatList, ActivityIndicator, TouchableOpacity, RefreshControl, StatusBar } from 'react-native';
import MapView, { Marker } from 'react-native-maps';
import { MaterialIcons } from '@expo/vector-icons';
import { useArboles } from '../hooks/useArboles';
import { styles } from './ArbolListScreen.styles';

interface ArbolListScreenProps {
  zonaId: string;
  onBack: () => void;
  onAdd: () => void;
}

export function ArbolListScreen({ zonaId, onBack, onAdd }: ArbolListScreenProps) {
  const { arboles, loading, error, refresh } = useArboles(zonaId);

  useEffect(() => {
    refresh();
  }, [refresh]);

  if (loading && arboles.length === 0) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#10B981" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={20} color="#8FA3A9" />
        </TouchableOpacity>
        <Text style={styles.title}>Árboles de la Zona</Text>
      </View>
      
      {error && (
        <View style={styles.errorContainer}>
          <MaterialIcons name="error-outline" size={20} color="#F43F5E" />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      <FlatList
        data={arboles}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        showsVerticalScrollIndicator={false}
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refresh} tintColor="#10B981" />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIconCircle}>
               <MaterialIcons name="forest" size={48} color="#10B981" />
            </View>
            <Text style={styles.emptyText}>No hay árboles registrados</Text>
            <Text style={styles.emptySubtext}>Comienza a registrar la flora plantada en esta zona.</Text>
          </View>
        }
        renderItem={({ item }) => (
          <View style={styles.card}>
            <View style={styles.cardContent}>
              <View style={styles.cardHeader}>
                <Text style={styles.cardLabel}>Código: {item.codigo || 'S/N'}</Text>
                <View style={styles.badge}>
                  <MaterialIcons name="eco" size={14} color="#FFFFFF" />
                  <Text style={styles.badgeText}>{item.altura_m ? `${item.altura_m}m` : 'N/A'}</Text>
                </View>
              </View>
              <Text style={styles.cardTitle}>Especie: {item.especie_id ? item.especie_id.substring(0,8) : 'Desconocida'}</Text>
              
              <Text style={styles.cardSubtitle}>
                <MaterialIcons name="height" size={14} color="#8FA3A9" /> Altura: {item.altura_m || '-'} m | Diámetro: {item.diametro_tronco_cm || '-'} cm
              </Text>
              <Text style={styles.cardSubtitle}>
                <MaterialIcons name="description" size={14} color="#8FA3A9" /> {item.observaciones || 'Sin observaciones'}
              </Text>
            </View>
            
            {item.latitud != null && item.longitud != null && (
              <View style={styles.mapContainer}>
                <MapView
                  style={styles.map}
                  initialRegion={{
                    latitude: item.latitud,
                    longitude: item.longitud,
                    latitudeDelta: 0.002,
                    longitudeDelta: 0.002,
                  }}
                  scrollEnabled={false}
                  zoomEnabled={false}
                >
                  <Marker coordinate={{ latitude: item.latitud, longitude: item.longitud }} />
                </MapView>
              </View>
            )}
            
            <View style={styles.cardFooter}>
               <Text style={styles.statsText}>Registrado: {item.fecha_creacion ? new Date(item.fecha_creacion).toLocaleDateString() : 'N/A'}</Text>
            </View>
          </View>
        )}
      />

      <TouchableOpacity style={styles.fab} onPress={onAdd} activeOpacity={0.9}>
        <MaterialIcons name="add" size={28} color="#061114" />
      </TouchableOpacity>
    </View>
  );
}
