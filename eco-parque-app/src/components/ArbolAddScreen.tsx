import React, { useState, useEffect } from 'react';
import { View, Text, TextInput, TouchableOpacity, ActivityIndicator, Alert, StatusBar, ScrollView } from 'react-native';
import * as Location from 'expo-location';
import { MaterialIcons } from '@expo/vector-icons';
import { useArboles } from '../hooks/useArboles';
import { styles } from './ArbolAddScreen.styles';

interface ArbolAddScreenProps {
  zonaId: string;
  onBack: () => void;
  onSuccess: () => void;
}

export function ArbolAddScreen({ zonaId, onBack, onSuccess }: ArbolAddScreenProps) {
  const { addArbol } = useArboles(zonaId);
  const [codigo, setCodigo] = useState('');
  const [especieId, setEspecieId] = useState('');
  const [alturaM, setAlturaM] = useState('');
  const [diametroCm, setDiametroCm] = useState('');
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
    if (!formLocation) {
      Alert.alert('Oops', 'Aún no se ha obtenido la ubicación GPS del árbol.');
      return;
    }

    setIsSaving(true);
    
    const parsedAltura = alturaM ? parseFloat(alturaM) : undefined;
    const parsedDiametro = diametroCm ? parseFloat(diametroCm) : undefined;

    const result = await addArbol({ 
      codigo: codigo || undefined,
      especie_id: especieId || undefined,
      latitud: formLocation.lat,
      longitud: formLocation.lng,
      altura_m: !isNaN(parsedAltura!) ? parsedAltura : undefined,
      diametro_tronco_cm: !isNaN(parsedDiametro!) ? parsedDiametro : undefined,
    });
    
    setIsSaving(false);
    
    if (result.success) {
      onSuccess();
    } else {
      Alert.alert('Error al crear', result.error || 'Ocurrió un error al registrar el árbol');
    }
  };

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

        <View style={styles.inputContainer}>
          <MaterialIcons name="eco" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="ID de Especie (Temporal)"
            placeholderTextColor="#8FA3A9"
            value={especieId}
            onChangeText={setEspecieId}
          />
        </View>
        
        <View style={styles.inputContainer}>
          <MaterialIcons name="height" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Altura en metros"
            placeholderTextColor="#8FA3A9"
            value={alturaM}
            onChangeText={setAlturaM}
            keyboardType="numeric"
          />
        </View>

        <View style={styles.inputContainer}>
          <MaterialIcons name="straighten" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Diámetro del tronco (cm)"
            placeholderTextColor="#8FA3A9"
            value={diametroCm}
            onChangeText={setDiametroCm}
            keyboardType="numeric"
          />
        </View>
        
        <TouchableOpacity style={styles.primaryButton} onPress={handleAdd} activeOpacity={0.8} disabled={isSaving || isFetchingLocation}>
          {isSaving ? (
            <ActivityIndicator color="#061114" />
          ) : (
            <Text style={styles.primaryButtonText}>Guardar Árbol</Text>
          )}
        </TouchableOpacity>
        <View style={{ height: 40 }} />
      </ScrollView>
    </View>
  );
}
