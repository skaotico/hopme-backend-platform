import React, { useEffect, useState, useRef } from 'react';
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  Animated,
  Alert,
  StatusBar,
  RefreshControl,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useArboles } from '../../hooks/useArboles';
import { CatalogoService } from '../../services/catalogo.service';
import { EspecieArbol, EstadoArbol } from '../../dto/catalogo.dto';
import { Arbol } from '../../dto/arbol.dto';
import { styles } from './ArbolListScreen.styles';
import { FadeSlideCard, PressableScale } from '../ui/AnimatedCards';

interface ArbolListScreenProps {
  zonaId: string;
  onBack: () => void;
  onAdd: () => void;
  onEdit: (arbol: Arbol) => void;
  onSensores?: (arbol: Arbol) => void;
}

// Animated card item using FadeSlideCard and PressableScale
function AnimatedCard({ item, index, onEdit, onDelete, onSensores, especies, estados }: {
  item: Arbol;
  index: number;
  onEdit: () => void;
  onDelete: () => void;
  onSensores?: () => void;
  especies: EspecieArbol[];
  estados: EstadoArbol[];
}) {
  const matchedEspecie = especies.find((e) => e.id === item.especie_id);
  const matchedEstado = estados.find((e) => e.id === item.estado_id);

  const getEspecieTitle = () => {
    if (matchedEspecie) {
      return matchedEspecie.nombre_cientifico;
    }
    return 'Especie Desconocida';
  };

  const getEspecieSub = () => {
    if (matchedEspecie?.nombre_comun) {
      return matchedEspecie.nombre_comun;
    }
    return matchedEspecie?.familia ? `Familia: ${matchedEspecie.familia}` : null;
  };

  const getEstadoName = () => {
    if (matchedEstado) {
      return matchedEstado.nombre || matchedEstado.codigo;
    }
    return 'Salud N/A';
  };

  return (
    <FadeSlideCard delay={index * 50} style={{ marginBottom: 16 }}>
      <PressableScale onPress={onSensores || onEdit}>
        <View style={[styles.card, { position: 'relative', overflow: 'hidden' }]}>
          {/* Hero Icon de fondo */}
          <MaterialIcons 
            name="park" 
            size={120} 
            color="rgba(20, 184, 166, 0.05)" 
            style={{ position: 'absolute', right: -20, bottom: -20 }} 
          />

          <View style={styles.cardContent}>
            <View style={styles.cardHeader}>
              <View style={{ flexDirection: 'row', alignItems: 'center', gap: 6 }}>
                <MaterialIcons name="park" size={16} color="#14B8A6" />
                <Text style={styles.cardLabel}>CÓDIGO: {item.codigo || 'S/N'}</Text>
              </View>
              <View style={styles.badge}>
                <MaterialIcons name="health-and-safety" size={14} color="#14B8A6" />
                <Text style={styles.badgeText}>{getEstadoName()}</Text>
              </View>
            </View>

            <Text style={styles.cardTitle}>{getEspecieTitle()}</Text>

            {getEspecieSub() && (
              <Text style={styles.cardSubtitle}>
                <MaterialIcons name="spa" size={14} color="#8FA3A9" /> {getEspecieSub()}
              </Text>
            )}

            <Text style={styles.cardSubtitle}>
              <MaterialIcons name="height" size={14} color="#8FA3A9" /> Altura: {item.altura_m ? `${item.altura_m}m` : '-'} | Copa: {item.ancho_copa_m ? `${item.ancho_copa_m}m` : '-'} | Tronco: {item.diametro_tronco_cm ? `${item.diametro_tronco_cm}cm` : '-'}
            </Text>

            {item.observaciones && (
              <Text style={styles.cardSubtitle}>
                <MaterialIcons name="chat-bubble-outline" size={14} color="#8FA3A9" /> {item.observaciones}
              </Text>
            )}
          </View>

          <View style={styles.cardFooter}>
            <Text style={styles.statsText}>
              Plantación: {item.fecha_plantacion ? new Date(item.fecha_plantacion).toLocaleDateString() : 'N/A'}
            </Text>

            <View style={styles.cardActions}>
              {onSensores && (
                <TouchableOpacity style={styles.actionButton} onPress={onSensores} activeOpacity={0.7}>
                  <MaterialIcons name="sensors" size={16} color="#F59E0B" />
                  <Text style={[styles.actionButtonText, { color: '#F59E0B' }]}>Sensores</Text>
                </TouchableOpacity>
              )}
              <TouchableOpacity style={styles.actionButton} onPress={onEdit} activeOpacity={0.7}>
                <MaterialIcons name="edit" size={16} color="#14B8A6" />
                <Text style={styles.actionButtonText}>Editar</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.actionButton, styles.actionButtonDanger]}
                onPress={onDelete}
                activeOpacity={0.7}
              >
                <MaterialIcons name="delete-outline" size={16} color="#F43F5E" />
                <Text style={[styles.actionButtonText, styles.actionButtonTextDanger]}>Borrar</Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </PressableScale>
    </FadeSlideCard>
  );
}

export function ArbolListScreen({ zonaId, onBack, onAdd, onEdit, onSensores }: ArbolListScreenProps) {
  const { arboles, loading, error, refresh, removeArbol } = useArboles(zonaId);

  // Catalogs to resolve names
  const [especies, setEspecies] = useState<EspecieArbol[]>([]);
  const [estados, setEstados] = useState<EstadoArbol[]>([]);

  useEffect(() => {
    refresh();

    // Fetch catalog list to map names in render
    (async () => {
      try {
        const [loadedEspecies, loadedEstados] = await Promise.all([
          CatalogoService.list('especies'),
          CatalogoService.list('estados-arbol'),
        ]);
        setEspecies(loadedEspecies as EspecieArbol[]);
        setEstados(loadedEstados as EstadoArbol[]);
      } catch (err) {
        console.warn('Error fetching catalogs in ArbolListScreen:', err);
      }
    })();
  }, [refresh]);

  const handleDelete = (item: Arbol) => {
    const matched = especies.find((e) => e.id === item.especie_id);
    const labelName = matched ? matched.nombre_cientifico : item.codigo || 'Árbol';

    Alert.alert(
      'Eliminar Árbol',
      `¿Estás seguro de que deseas eliminar este espécimen ("${labelName}") del ecoparque?`,
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Eliminar',
          style: 'destructive',
          onPress: async () => {
            const result = await removeArbol(item.id);
            if (!result.success) {
              Alert.alert('Error', result.error || 'No se pudo eliminar el árbol.');
            }
          },
        },
      ]
    );
  };

  if (loading && arboles.length === 0) {
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
        refreshControl={<RefreshControl refreshing={loading} onRefresh={refresh} tintColor="#14B8A6" />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIconCircle}>
              <MaterialIcons name="forest" size={48} color="#14B8A6" />
            </View>
            <Text style={styles.emptyText}>No hay árboles registrados</Text>
            <Text style={styles.emptySubtext}>Comienza a registrar la flora plantada en esta zona.</Text>
          </View>
        }
        renderItem={({ item, index }) => (
          <AnimatedCard
            item={item}
            index={index}
            onEdit={() => onEdit(item)}
            onDelete={() => handleDelete(item)}
            onSensores={onSensores ? () => onSensores(item) : undefined}
            especies={especies}
            estados={estados}
          />
        )}
      />

      <TouchableOpacity style={styles.fab} onPress={onAdd} activeOpacity={0.9}>
        <MaterialIcons name="add" size={28} color="#061114" />
      </TouchableOpacity>
    </View>
  );
}
