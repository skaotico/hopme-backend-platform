import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, ActivityIndicator, Alert, StatusBar } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useZonas } from '../../hooks/useZonas';
import { styles } from './ZonaAddScreen.styles';


interface ZonaAddScreenProps {
  parqueId: string;
  onBack: () => void;
  onSuccess: () => void;
}

export function ZonaAddScreen({ parqueId, onBack, onSuccess }: ZonaAddScreenProps) {
  const { addZona } = useZonas(parqueId);
  const [nombre, setNombre] = useState('');
  const [descripcion, setDescripcion] = useState('');
  const [areaM2, setAreaM2] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  const handleAdd = async () => {
    if (!nombre) {
      Alert.alert('Oops', 'Por favor ingresa al menos el nombre de la zona');
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

    const result = await addZona({ 
      nombre, 
      descripcion: descripcion || undefined,
      area_m2: parsedArea,
    });
    
    setIsSaving(false);
    
    if (result.success) {
      onSuccess();
    } else {
      Alert.alert('Error al crear', result.error || 'Ocurrió un error al guardar la zona');
    }
  };

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={20} color="#8FA3A9" />
        </TouchableOpacity>
        <Text style={styles.title}>Agregar Zona</Text>
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
        
        <TouchableOpacity style={styles.primaryButton} onPress={handleAdd} activeOpacity={0.8} disabled={isSaving}>
          {isSaving ? (
            <ActivityIndicator color="#061114" />
          ) : (
            <Text style={styles.primaryButtonText}>Guardar Zona</Text>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );
}
