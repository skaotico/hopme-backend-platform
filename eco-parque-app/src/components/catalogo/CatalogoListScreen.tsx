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
  StyleSheet,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useCatalogo } from '../../hooks/useCatalogo';
import { CatalogItem, CatalogType, CATALOGS_CONFIG } from '../../dto/catalogo.dto';

interface CatalogoListScreenProps {
  type: CatalogType;
  onOpenMenu: () => void;
  onAdd: () => void;
  onEdit: (item: CatalogItem) => void;
}

export function CatalogoListScreen({ type, onOpenMenu, onAdd, onEdit }: CatalogoListScreenProps) {
  const config = CATALOGS_CONFIG[type];
  const { items, loading, error, refresh, removeItem } = useCatalogo(type);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleDelete = (item: any) => {
    const itemName = item.nombre_cientifico || item.nombre || item.codigo || 'Registro';
    Alert.alert(
      'Eliminar Registro',
      `¿Estás seguro de que deseas eliminar "${itemName}" del catálogo?`,
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Eliminar',
          style: 'destructive',
          onPress: async () => {
            const result = await removeItem(item.id);
            if (!result.success) {
              Alert.alert('Error', result.error || 'No se pudo eliminar el elemento.');
            }
          },
        },
      ],
    );
  };

  const getTitle = (item: any) => {
    if (type === 'especies') {
      return item.nombre_cientifico;
    }
    return item.nombre || item.codigo;
  };

  const getSubtitle = (item: any) => {
    if (type === 'especies') {
      return item.nombre_comun ? `Común: ${item.nombre_comun}` : item.familia ? `Familia: ${item.familia}` : null;
    }
    return item.codigo ? `Código: ${item.codigo}` : null;
  };

  const getDescription = (item: any) => {
    return item.descripcion || (item.unidad_medida ? `Unidad: ${item.unidad_medida}` : null);
  };

  if (loading && items.length === 0) {
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
        <TouchableOpacity style={styles.menuButton} onPress={onOpenMenu} activeOpacity={0.7}>
          <MaterialIcons name="menu" size={24} color="#14B8A6" />
        </TouchableOpacity>
        <Text style={styles.title}>{config.displayName}</Text>
      </View>

      {error && (
        <View style={styles.errorContainer}>
          <MaterialIcons name="error-outline" size={20} color="#F43F5E" />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        showsVerticalScrollIndicator={false}
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refresh} tintColor="#14B8A6" />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIconCircle}>
              <MaterialIcons name={config.icon as any} size={48} color="#14B8A6" />
            </View>
            <Text style={styles.emptyText}>Catálogo vacío</Text>
            <Text style={styles.emptySubtext}>Crea el primer elemento de este mantenedor.</Text>
          </View>
        }
        renderItem={({ item }: { item: any }) => (
          <View style={styles.card}>
            <View style={styles.cardHeader}>
              <Text style={styles.cardLabel}>Registro</Text>
              {type === 'especies' && (
                <View style={[styles.badge, { backgroundColor: item.activo ? 'rgba(20, 184, 166, 0.2)' : 'rgba(244, 63, 94, 0.2)' }]}>
                  <Text style={[styles.badgeText, { color: item.activo ? '#14B8A6' : '#F43F5E' }]}>
                    {item.activo ? 'Activo' : 'Inactivo'}
                  </Text>
                </View>
              )}
            </View>

            <Text style={styles.cardTitle}>{getTitle(item)}</Text>

            {getSubtitle(item) && (
              <Text style={styles.cardSubtitle}>
                <MaterialIcons name="label-outline" size={14} color="#8FA3A9" /> {getSubtitle(item)}
              </Text>
            )}

            {getDescription(item) && (
              <Text style={styles.cardSubtitle}>
                <MaterialIcons name="description" size={14} color="#8FA3A9" /> {getDescription(item)}
              </Text>
            )}

            <View style={styles.cardActions}>
              <TouchableOpacity
                style={styles.actionButton}
                onPress={() => onEdit(item)}
                activeOpacity={0.7}
              >
                <MaterialIcons name="edit" size={18} color="#14B8A6" />
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
          </View>
        )}
      />

      <TouchableOpacity style={styles.fab} onPress={onAdd} activeOpacity={0.9}>
        <MaterialIcons name="add" size={28} color="#061114" />
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#061114',
  },
  container: {
    flex: 1,
    backgroundColor: '#061114',
    paddingTop: 50,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 24,
    marginBottom: 24,
  },
  menuButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 16,
  },
  title: {
    fontSize: 24,
    fontWeight: '800',
    color: '#FFFFFF',
    letterSpacing: -0.5,
  },
  errorContainer: {
    flexDirection: 'row',
    backgroundColor: 'rgba(244, 63, 94, 0.1)',
    marginHorizontal: 24,
    marginBottom: 16,
    padding: 16,
    borderRadius: 16,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: 'rgba(244, 63, 94, 0.2)',
  },
  errorText: {
    color: '#F43F5E',
    marginLeft: 8,
    fontWeight: '500',
  },
  listContent: {
    paddingHorizontal: 24,
    paddingBottom: 120,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 80,
    paddingHorizontal: 20,
  },
  emptyIconCircle: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: '#112226',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 24,
    borderWidth: 1,
    borderColor: '#1D343B',
  },
  emptyText: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFFFFF',
    marginBottom: 8,
    textAlign: 'center',
  },
  emptySubtext: {
    fontSize: 15,
    color: '#8FA3A9',
    textAlign: 'center',
    lineHeight: 22,
  },
  card: {
    backgroundColor: '#112226',
    marginBottom: 16,
    borderRadius: 24,
    borderWidth: 1,
    borderColor: '#1D343B',
    padding: 20,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 10,
  },
  cardLabel: {
    fontSize: 12,
    color: '#8FA3A9',
    fontWeight: '500',
  },
  badge: {
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 8,
  },
  badgeText: {
    fontSize: 11,
    fontWeight: 'bold',
  },
  cardTitle: {
    fontSize: 20,
    fontWeight: '800',
    color: '#FFFFFF',
    marginBottom: 6,
  },
  cardSubtitle: {
    fontSize: 14,
    color: '#8FA3A9',
    fontWeight: '500',
    marginBottom: 4,
  },
  cardActions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 8,
    marginTop: 12,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: '#1D343B',
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
    backgroundColor: 'rgba(20, 184, 166, 0.1)',
    borderWidth: 1,
    borderColor: 'rgba(20, 184, 166, 0.2)',
  },
  actionButtonText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#14B8A6',
  },
  actionButtonDanger: {
    backgroundColor: 'rgba(244, 63, 94, 0.1)',
    borderColor: 'rgba(244, 63, 94, 0.2)',
  },
  actionButtonTextDanger: {
    color: '#F43F5E',
  },
  fab: {
    position: 'absolute',
    bottom: 32,
    alignSelf: 'center',
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#14B8A6',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#14B8A6',
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.3,
    shadowRadius: 10,
    elevation: 8,
    borderWidth: 4,
    borderColor: '#061114',
  },
});
