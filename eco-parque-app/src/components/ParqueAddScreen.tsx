import React, { useState, useEffect } from 'react';
import { View, Text, TextInput, TouchableOpacity, ActivityIndicator, Alert, StatusBar } from 'react-native';
import * as Location from 'expo-location';
import { MaterialIcons } from '@expo/vector-icons';
import { useParques } from '../hooks/useParques';
import { styles } from './ParqueAddScreen.styles';

interface ParqueAddScreenProps {
  onBack: () => void;
  onSuccess: () => void;
}

export function ParqueAddScreen({ onBack, onSuccess }: ParqueAddScreenProps) {
  const { addParque } = useParques();
  const [nuevoNombre, setNuevoNombre] = useState('');
  const [nuevaDescripcion, setNuevaDescripcion] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  
  const [formLocation, setFormLocation] = useState<{lat: number, lng: number} | null>(null);
  const [isFetchingLocation, setIsFetchingLocation] = useState(false);

  useEffect(() => {
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
          lng: location.coords.longitude
        });
      } catch (error) {
        console.warn(error);
      }
      setIsFetchingLocation(false);
    })();
  }, []);

  const handleAdd = async () => {
    if (!nuevoNombre || !nuevaDescripcion) {
      Alert.alert('Oops', 'Por favor completa nombre y descripción del parque');
      return;
    }
    
    if (!formLocation) {
      Alert.alert('Oops', 'Aún no se ha obtenido la ubicación, espera un momento.');
      return;
    }

    setIsSaving(true);
    const direccion = `${formLocation.lat}, ${formLocation.lng}`;

    const result = await addParque({ 
      nombre: nuevoNombre, 
      descripcion: nuevaDescripcion,
      direccion: direccion,
    });
    setIsSaving(false);
    
    if (result.success) {
      onSuccess();
    } else {
      Alert.alert('Error al crear', result.error || 'Ocurrió un error al guardar el parque');
    }
  };

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={20} color="#8FA3A9" />
        </TouchableOpacity>
        <Text style={styles.title}>Agregar Parque</Text>
      </View>
      
      <View style={styles.formContent}>
        <View style={styles.inputContainer}>
          <MaterialIcons name="park" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Nombre del parque"
            placeholderTextColor="#8FA3A9"
            value={nuevoNombre}
            onChangeText={setNuevoNombre}
          />
        </View>
        
        <View style={styles.inputContainer}>
          <MaterialIcons name="description" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Descripción del parque"
            placeholderTextColor="#8FA3A9"
            value={nuevaDescripcion}
            onChangeText={setNuevaDescripcion}
          />
        </View>
        
        {isFetchingLocation ? (
          <View style={styles.locationInfo}>
            <ActivityIndicator size="small" color="#14B8A6" />
            <Text style={styles.locationText}>Obteniendo ubicación del dispositivo...</Text>
          </View>
        ) : formLocation ? (
          <View style={styles.locationInfo}>
            <MaterialIcons name="my-location" size={18} color="#14B8A6" />
            <Text style={styles.locationText}>
              Ubicación lista: {formLocation.lat.toFixed(5)}, {formLocation.lng.toFixed(5)}
            </Text>
          </View>
        ) : (
          <View style={[styles.locationInfo, { borderColor: '#F43F5E', backgroundColor: 'rgba(244,63,94,0.1)' }]}>
            <MaterialIcons name="location-disabled" size={18} color="#F43F5E" />
            <Text style={[styles.locationText, { color: '#F43F5E' }]}>No se pudo obtener la ubicación</Text>
          </View>
        )}
        
        <TouchableOpacity style={styles.primaryButton} onPress={handleAdd} activeOpacity={0.8} disabled={isSaving || isFetchingLocation}>
          {isSaving ? (
            <ActivityIndicator color="#061114" />
          ) : (
            <Text style={styles.primaryButtonText}>Guardar Parque</Text>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );
}
