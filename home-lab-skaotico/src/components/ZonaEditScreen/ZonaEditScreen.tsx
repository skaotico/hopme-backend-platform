import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  StatusBar,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useZonas } from '../../hooks/useZonas';
import { Zona } from '../../dto/zona.dto';
import { styles } from '../ZonaAddScreen/ZonaAddScreen.styles';

interface ZonaEditScreenProps {
  parqueId: string;
  zona: Zona;
  onBack: () => void;
  onSuccess: () => void;
}

export function ZonaEditScreen({ parqueId, zona, onBack, onSuccess }: ZonaEditScreenProps) {
  const { updateZona } = useZonas(parqueId);
  const [nombre, setNombre] = useState(zona.nombre);
  const [descripcion, setDescripcion] = useState(zona.descripcion ?? '');
  const [areaM2, setAreaM2] = useState(zona.area_m2 != null ? String(zona.area_m2) : '');
  const [isSaving, setIsSaving] = useState(false);

  const handleSave = async () => {
    if (!nombre.trim()) {
      Alert.alert('Oops', 'Por favor ingresa el nombre de la zona');
      return;
    }

    setIsSaving(true);

    let parsedArea: number | undefined = undefined;
    if (areaM2) {
      const parsed = parseFloat(areaM2);
      if (!isNaN(parsed)) {
        parsedArea = parsed;
      }
    }

    const result = await updateZona(zona.id, {
      nombre: nombre.trim(),
      descripcion: descripcion.trim() || undefined,
      area_m2: parsedArea,
    });

    setIsSaving(false);

    if (result.success) {
      onSuccess();
    } else {
      Alert.alert('Error al actualizar', result.error || 'Ocurrió un error al guardar los cambios');
    }
  };

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={20} color="#8FA3A9" />
        </TouchableOpacity>
        <Text style={styles.title}>Editar Zona</Text>
      </View>

      <View style={styles.formContent}>
        <View style={styles.inputContainer}>
          <MaterialIcons name="fullscreen" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Nombre de la zona"
            placeholderTextColor="#8FA3A9"
            value={nombre}
            onChangeText={setNombre}
          />
        </View>

        <View style={styles.inputContainer}>
          <MaterialIcons name="description" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Descripción (opcional)"
            placeholderTextColor="#8FA3A9"
            value={descripcion}
            onChangeText={setDescripcion}
          />
        </View>

        <View style={styles.inputContainer}>
          <MaterialIcons name="straighten" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Área en m² (opcional)"
            placeholderTextColor="#8FA3A9"
            value={areaM2}
            onChangeText={setAreaM2}
            keyboardType="numeric"
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
      </View>
    </View>
  );
}
