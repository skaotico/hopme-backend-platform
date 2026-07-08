import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  StatusBar,
  ScrollView,
  Modal,
  FlatList,
  TouchableWithoutFeedback,
} from 'react-native';
import * as Location from 'expo-location';
import { MaterialIcons } from '@expo/vector-icons';
import { useArboles } from '../../hooks/useArboles';
import { CatalogoService } from '../../services/catalogo.service';
import { EspecieArbol, EstadoArbol } from '../../dto/catalogo.dto';
import { styles } from './ArbolAddScreen.styles';

interface ArbolAddScreenProps {
  zonaId: string;
  onBack: () => void;
  onSuccess: () => void;
}

export function ArbolAddScreen({ zonaId, onBack, onSuccess }: ArbolAddScreenProps) {
  const { addArbol } = useArboles(zonaId);
  const [codigo, setCodigo] = useState('');
  const [alturaM, setAlturaM] = useState('');
  const [diametroCm, setDiametroCm] = useState('');
  const [anchoCopaM, setAnchoCopaM] = useState('');
  const [edadEstimadaAnios, setEdadEstimadaAnios] = useState('');
  const [observaciones, setObservaciones] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  // Catalogs state
  const [especies, setEspecies] = useState<EspecieArbol[]>([]);
  const [estados, setEstados] = useState<EstadoArbol[]>([]);
  const [selectedEspecie, setSelectedEspecie] = useState<EspecieArbol | null>(null);
  const [selectedEstado, setSelectedEstado] = useState<EstadoArbol | null>(null);
  const [isLoadingCatalogs, setIsLoadingCatalogs] = useState(true);

  // Pickers visibility
  const [isEspeciePickerOpen, setIsEspeciePickerOpen] = useState(false);
  const [isEstadoPickerOpen, setIsEstadoPickerOpen] = useState(false);

  // GPS state
  const [formLocation, setFormLocation] = useState<{ lat: number; lng: number } | null>(null);
  const [isFetchingLocation, setIsFetchingLocation] = useState(false);

  // Load catalogs and GPS position
  useEffect(() => {
    // 1. Fetch Location
    (async () => {
      setIsFetchingLocation(true);
      let { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        setIsFetchingLocation(false);
        return;
      }
      try {
        let location = await Location.getCurrentPositionAsync({});
        setFormLocation({
          lat: location.coords.latitude,
          lng: location.coords.longitude,
        });
      } catch (error) {
        console.warn(error);
      }
      setIsFetchingLocation(false);
    })();

    // 2. Fetch species and states from catalog service
    (async () => {
      setIsLoadingCatalogs(true);
      try {
        const [loadedEspecies, loadedEstados] = await Promise.all([
          CatalogoService.list('especies'),
          CatalogoService.list('estados-arbol'),
        ]);
        // Only keep active species if active flag exists
        const activeEspecies = (loadedEspecies as EspecieArbol[]).filter(
          (e) => e.activo !== false
        );
        setEspecies(activeEspecies);
        setEstados(loadedEstados as EstadoArbol[]);

        // Auto select first options as default if available
        if (activeEspecies.length > 0) setSelectedEspecie(activeEspecies[0]);
        if (loadedEstados.length > 0) setSelectedEstado(loadedEstados[0] as EstadoArbol);
      } catch (err) {
        console.error('Error cargando catálogos para registro de árboles:', err);
        Alert.alert(
          'Error de Catálogo',
          'No se pudieron cargar los catálogos de especies y estados. Por favor inténtalo más tarde.'
        );
      } finally {
        setIsLoadingCatalogs(false);
      }
    })();
  }, []);

  const handleAdd = async () => {
    if (!formLocation) {
      Alert.alert('Oops', 'Aún no se ha obtenido la ubicación GPS del árbol.');
      return;
    }

    if (!selectedEspecie) {
      Alert.alert('Oops', 'Por favor selecciona la especie del árbol.');
      return;
    }

    if (!selectedEstado) {
      Alert.alert('Oops', 'Por favor selecciona el estado de salud del árbol.');
      return;
    }

    setIsSaving(true);

    const parsedAltura = alturaM ? parseFloat(alturaM) : undefined;
    const parsedDiametro = diametroCm ? parseFloat(diametroCm) : undefined;
    const parsedAnchoCopa = anchoCopaM ? parseFloat(anchoCopaM) : undefined;
    const parsedEdad = edadEstimadaAnios ? parseInt(edadEstimadaAnios, 10) : undefined;

    const result = await addArbol({
      codigo: codigo.trim() || undefined,
      especie_id: selectedEspecie.id,
      estado_id: selectedEstado.id,
      latitud: formLocation.lat,
      longitud: formLocation.lng,
      altura_m: parsedAltura != null && !isNaN(parsedAltura) ? parsedAltura : undefined,
      diametro_tronco_cm: parsedDiametro != null && !isNaN(parsedDiametro) ? parsedDiametro : undefined,
      ancho_copa_m: parsedAnchoCopa != null && !isNaN(parsedAnchoCopa) ? parsedAnchoCopa : undefined,
      edad_estimada_anios: parsedEdad != null && !isNaN(parsedEdad) ? parsedEdad : undefined,
      observaciones: observaciones.trim() || undefined,
    });

    setIsSaving(false);

    if (result.success) {
      onSuccess();
    } else {
      Alert.alert('Error al crear', result.error || 'Ocurrió un error al registrar el árbol');
    }
  };

  if (isLoadingCatalogs) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#14B8A6" />
        <Text style={{ color: '#8FA3A9', marginTop: 12, fontWeight: '600' }}>
          Cargando catálogos...
        </Text>
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
        <Text style={styles.title}>Registrar Árbol</Text>
      </View>

      <ScrollView contentContainerStyle={styles.formContent} showsVerticalScrollIndicator={false}>
        {isFetchingLocation ? (
          <View style={styles.locationInfo}>
            <ActivityIndicator size="small" color="#14B8A6" />
            <Text style={styles.locationText}>Capturando GPS del árbol...</Text>
          </View>
        ) : formLocation ? (
          <View style={styles.locationInfo}>
            <MaterialIcons name="my-location" size={18} color="#14B8A6" />
            <Text style={styles.locationText}>
              Ubicación guardada: {formLocation.lat.toFixed(5)}, {formLocation.lng.toFixed(5)}
            </Text>
          </View>
        ) : (
          <View style={[styles.locationInfo, { borderColor: '#F43F5E', backgroundColor: 'rgba(244,63,94,0.1)' }]}>
            <MaterialIcons name="location-disabled" size={18} color="#F43F5E" />
            <Text style={[styles.locationText, { color: '#F43F5E' }]}>No se pudo obtener el GPS</Text>
          </View>
        )}

        {/* Codigo */}
        <View style={styles.inputContainer}>
          <MaterialIcons name="qr-code" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Código del árbol (Opcional)"
            placeholderTextColor="#8FA3A9"
            value={codigo}
            onChangeText={setCodigo}
          />
        </View>

        {/* Selector de Especie */}
        <TouchableOpacity
          style={styles.selectorButton}
          onPress={() => setIsEspeciePickerOpen(true)}
          activeOpacity={0.7}
        >
          <View style={styles.selectorLeft}>
            <MaterialIcons name="nature" size={20} color="#14B8A6" />
            {selectedEspecie ? (
              <Text style={styles.selectorValue}>{selectedEspecie.nombre_cientifico}</Text>
            ) : (
              <Text style={styles.selectorPlaceholder}>Seleccionar Especie</Text>
            )}
          </View>
          <MaterialIcons name="arrow-drop-down" size={24} color="#8FA3A9" />
        </TouchableOpacity>

        {/* Selector de Estado */}
        <TouchableOpacity
          style={styles.selectorButton}
          onPress={() => setIsEstadoPickerOpen(true)}
          activeOpacity={0.7}
        >
          <View style={styles.selectorLeft}>
            <MaterialIcons name="health-and-safety" size={20} color="#14B8A6" />
            {selectedEstado ? (
              <Text style={styles.selectorValue}>{selectedEstado.nombre || selectedEstado.codigo}</Text>
            ) : (
              <Text style={styles.selectorPlaceholder}>Seleccionar Estado de Salud</Text>
            )}
          </View>
          <MaterialIcons name="arrow-drop-down" size={24} color="#8FA3A9" />
        </TouchableOpacity>

        {/* Altura */}
        <View style={styles.inputContainer}>
          <MaterialIcons name="height" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Altura en metros (Opcional)"
            placeholderTextColor="#8FA3A9"
            value={alturaM}
            onChangeText={setAlturaM}
            keyboardType="numeric"
          />
        </View>

        {/* Diametro */}
        <View style={styles.inputContainer}>
          <MaterialIcons name="straighten" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Diámetro del tronco en cm (Opcional)"
            placeholderTextColor="#8FA3A9"
            value={diametroCm}
            onChangeText={setDiametroCm}
            keyboardType="numeric"
          />
        </View>

        {/* Ancho Copa */}
        <View style={styles.inputContainer}>
          <MaterialIcons name="fullscreen" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Ancho de la copa en metros (Opcional)"
            placeholderTextColor="#8FA3A9"
            value={anchoCopaM}
            onChangeText={setAnchoCopaM}
            keyboardType="numeric"
          />
        </View>

        {/* Edad Estimada */}
        <View style={styles.inputContainer}>
          <MaterialIcons name="schedule" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Edad estimada en años (Opcional)"
            placeholderTextColor="#8FA3A9"
            value={edadEstimadaAnios}
            onChangeText={setEdadEstimadaAnios}
            keyboardType="numeric"
          />
        </View>

        {/* Observaciones */}
        <View style={[styles.inputContainer, styles.textAreaContainer]}>
          <MaterialIcons name="description" size={20} color="#8FA3A9" style={[styles.icon, styles.textAreaIcon]} />
          <TextInput
            style={[styles.input, styles.textAreaInput]}
            placeholder="Observaciones (Opcional)"
            placeholderTextColor="#8FA3A9"
            value={observaciones}
            onChangeText={setObservaciones}
            multiline
            numberOfLines={3}
          />
        </View>

        <TouchableOpacity
          style={styles.primaryButton}
          onPress={handleAdd}
          activeOpacity={0.8}
          disabled={isSaving || isFetchingLocation}
        >
          {isSaving ? (
            <ActivityIndicator color="#061114" />
          ) : (
            <Text style={styles.primaryButtonText}>Guardar Árbol</Text>
          )}
        </TouchableOpacity>
        <View style={{ height: 40 }} />
      </ScrollView>

      {/* Modal Picker Especie */}
      <Modal
        visible={isEspeciePickerOpen}
        transparent
        animationType="slide"
        onRequestClose={() => setIsEspeciePickerOpen(false)}
      >
        <View style={styles.modalOverlay}>
          <TouchableWithoutFeedback onPress={() => setIsEspeciePickerOpen(false)}>
            <View style={{ flex: 1 }} />
          </TouchableWithoutFeedback>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Especie del Árbol</Text>
              <TouchableOpacity
                style={styles.modalCloseButton}
                onPress={() => setIsEspeciePickerOpen(false)}
              >
                <MaterialIcons name="close" size={24} color="#8FA3A9" />
              </TouchableOpacity>
            </View>

            <FlatList
              data={especies}
              keyExtractor={(item) => item.id}
              contentContainerStyle={styles.modalList}
              renderItem={({ item }) => {
                const isActive = selectedEspecie?.id === item.id;
                return (
                  <TouchableOpacity
                    style={[styles.modalItem, isActive && styles.modalItemActive]}
                    onPress={() => {
                      setSelectedEspecie(item);
                      setIsEspeciePickerOpen(false);
                    }}
                    activeOpacity={0.7}
                  >
                    <Text style={styles.modalItemText}>{item.nombre_cientifico}</Text>
                    {item.nombre_comun && (
                      <Text style={styles.modalItemSubtext}>{item.nombre_comun}</Text>
                    )}
                  </TouchableOpacity>
                );
              }}
            />
          </View>
        </View>
      </Modal>

      {/* Modal Picker Estado */}
      <Modal
        visible={isEstadoPickerOpen}
        transparent
        animationType="slide"
        onRequestClose={() => setIsEstadoPickerOpen(false)}
      >
        <View style={styles.modalOverlay}>
          <TouchableWithoutFeedback onPress={() => setIsEstadoPickerOpen(false)}>
            <View style={{ flex: 1 }} />
          </TouchableWithoutFeedback>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Estado de Salud</Text>
              <TouchableOpacity
                style={styles.modalCloseButton}
                onPress={() => setIsEstadoPickerOpen(false)}
              >
                <MaterialIcons name="close" size={24} color="#8FA3A9" />
              </TouchableOpacity>
            </View>

            <FlatList
              data={estados}
              keyExtractor={(item) => item.id}
              contentContainerStyle={styles.modalList}
              renderItem={({ item }) => {
                const isActive = selectedEstado?.id === item.id;
                return (
                  <TouchableOpacity
                    style={[styles.modalItem, isActive && styles.modalItemActive]}
                    onPress={() => {
                      setSelectedEstado(item);
                      setIsEstadoPickerOpen(false);
                    }}
                    activeOpacity={0.7}
                  >
                    <Text style={styles.modalItemText}>{item.nombre || item.codigo}</Text>
                    {item.descripcion && (
                      <Text style={styles.modalItemSubtext}>{item.descripcion}</Text>
                    )}
                  </TouchableOpacity>
                );
              }}
            />
          </View>
        </View>
      </Modal>
    </View>
  );
}
