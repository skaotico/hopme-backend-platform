import React, { useEffect } from 'react';
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  RefreshControl,
  StatusBar,
  Alert,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useZonas } from '../../hooks/useZonas';
import { Zona } from '../../dto/zona.dto';
import { styles } from './ZonaListScreen.styles';

interface ZonaListScreenProps {
  parqueId: string;
  onBack: () => void;
  onAdd: () => void;
  onSelect: (id: string) => void;
  onEdit: (zona: Zona) => void;
}

export function ZonaListScreen({ parqueId, onBack, onAdd, onSelect, onEdit }: ZonaListScreenProps) {
  const { zonas, loading, error, refresh, deleteZona } = useZonas(parqueId);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleDelete = (zona: Zona) => {
    Alert.alert(
      'Eliminar Zona',
      `¿Estás seguro de que deseas eliminar "${zona.nombre}"? Esta acción no se puede deshacer.`,
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Eliminar',
          style: 'destructive',
          onPress: async () => {
            const result = await deleteZona(zona.id);
            if (!result.success) {
              Alert.alert('Error', result.error || 'No se pudo eliminar la zona.');
            }
          },
        },
      ],
    );
  };

  if (loading && zonas.length === 0) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#3B82F6" />
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
        <Text style={styles.title}>Zonas del Parque</Text>
      </View>

      {error && (
        <View style={styles.errorContainer}>
          <MaterialIcons name="error-outline" size={20} color="#F43F5E" />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      <FlatList
        data={zonas}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        showsVerticalScrollIndicator={false}
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refresh} tintColor="#3B82F6" />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIconCircle}>
              <MaterialIcons name="map" size={48} color="#3B82F6" />
            </View>
            <Text style={styles.emptyText}>Este parque no tiene zonas</Text>
            <Text style={styles.emptySubtext}>Crea una zona para subdividir el parque.</Text>
          </View>
        }
        renderItem={({ item }) => (
          <TouchableOpacity style={styles.card} activeOpacity={0.8} onPress={() => onSelect(item.id)}>
            <View style={styles.cardHeader}>
              <Text style={styles.cardLabel}>Zona</Text>
              <View style={styles.badge}>
                <MaterialIcons name="fullscreen" size={14} color="#FFFFFF" />
                <Text style={styles.badgeText}>{item.area_m2 ? `${item.area_m2} m²` : 'N/A'}</Text>
              </View>
            </View>
            <Text style={styles.cardTitle}>{item.nombre}</Text>
            <Text style={styles.cardSubtitle}>
              <MaterialIcons name="description" size={14} color="#8FA3A9" /> {item.descripcion || 'Sin descripción'}
            </Text>
            <View style={styles.cardActions}>
              <TouchableOpacity
                style={styles.actionButton}
                onPress={() => onEdit(item)}
                activeOpacity={0.7}
              >
                <MaterialIcons name="edit" size={18} color="#3B82F6" />
                <Text style={styles.actionButtonText}>Editar</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.actionButton, styles.actionButtonDanger]}
                onPress={() => handleDelete(item)}
                activeOpacity={0.7}
              >
                <MaterialIcons name="delete-outline" size={18} color="#F43F5E" />
                <Text style={[styles.actionButtonText, styles.actionButtonTextDanger]}>Eliminar</Text>
              </TouchableOpacity>
            </View>
          </TouchableOpacity>
        )}
      />

      <TouchableOpacity style={styles.fab} onPress={onAdd} activeOpacity={0.9}>
        <MaterialIcons name="add" size={28} color="#061114" />
      </TouchableOpacity>
    </View>
  );
}
