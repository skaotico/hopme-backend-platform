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
import { MaterialIcons } from '@expo/vector-icons';
import { useArboles } from '../../hooks/useArboles';
import { CatalogoService } from '../../services/catalogo.service';
import { EspecieArbol, EstadoArbol } from '../../dto/catalogo.dto';
import { Arbol } from '../../dto/arbol.dto';
import { styles } from '../ArbolAddScreen/ArbolAddScreen.styles';

interface ArbolEditScreenProps {
  zonaId: string;
  arbol: Arbol;
  onBack: () => void;
  onSuccess: () => void;
}

export function ArbolEditScreen({ zonaId, arbol, onBack, onSuccess }: ArbolEditScreenProps) {
  const { updateArbol } = useArboles(zonaId);
  const [codigo, setCodigo] = useState(arbol.codigo || '');
  const [alturaM, setAlturaM] = useState(arbol.altura_m ? String(arbol.altura_m) : '');
  const [diametroCm, setDiametroCm] = useState(arbol.diametro_tronco_cm ? String(arbol.diametro_tronco_cm) : '');
  const [anchoCopaM, setAnchoCopaM] = useState(arbol.ancho_copa_m ? String(arbol.ancho_copa_m) : '');
  const [edadEstimadaAnios, setEdadEstimadaAnios] = useState(arbol.edad_estimada_anios ? String(arbol.edad_estimada_anios) : '');
  const [observaciones, setObservaciones] = useState(arbol.observaciones || '');
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

  // Load catalogs and set initial selected values
  useEffect(() => {
    (async () => {
      setIsLoadingCatalogs(true);
      try {
        const [loadedEspecies, loadedEstados] = await Promise.all([
          CatalogoService.list('especies'),
          CatalogoService.list('estados-arbol'),
        ]);

        const activeEspecies = (loadedEspecies as EspecieArbol[]).filter(
          (e) => e.activo !== false || e.id === arbol.especie_id
        );
        setEspecies(activeEspecies);
        setEstados(loadedEstados as EstadoArbol[]);

        // Match initial values
        const currentEspecie = activeEspecies.find((e) => e.id === arbol.especie_id);
        const currentEstado = (loadedEstados as EstadoArbol[]).find((e) => e.id === arbol.estado_id);

        if (currentEspecie) setSelectedEspecie(currentEspecie);
        if (currentEstado) setSelectedEstado(currentEstado);
      } catch (err) {
        console.error('Error cargando catálogos para edición de árboles:', err);
        Alert.alert(
          'Error de Catálogo',
          'No se pudieron cargar los catálogos de especies y estados. Por favor inténtalo más tarde.'
        );
      } finally {
        setIsLoadingCatalogs(false);
      }
    })();
  }, [arbol]);

  const handleSave = async () => {
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

    const result = await updateArbol(arbol.id, {
      zona_id: zonaId,
      codigo: codigo.trim() || undefined,
      especie_id: selectedEspecie.id,
      estado_id: selectedEstado.id,
      latitud: arbol.latitud,
      longitud: arbol.longitud,
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
      Alert.alert('Error al actualizar', result.error || 'Ocurrió un error al guardar los cambios del árbol');
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
        <Text style={styles.title}>Editar Árbol</Text>
      </View>

      <ScrollView contentContainerStyle={styles.formContent} showsVerticalScrollIndicator={false}>
        {/* GPS location info */}
        {arbol.latitud != null && arbol.longitud != null && (
          <View style={styles.locationInfo}>
            <MaterialIcons name="my-location" size={18} color="#14B8A6" />
            <Text style={styles.locationText}>
              Ubicación fija: {arbol.latitud.toFixed(5)}, {arbol.longitud.toFixed(5)}
            </Text>
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
          onPress={handleSave}
          activeOpacity={0.8}
          disabled={isSaving}
        >
          {isSaving ? (
            <ActivityIndicator color="#061114" />
          ) : (
            <Text style={styles.primaryButtonText}>Guardar Cambios</Text>
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
