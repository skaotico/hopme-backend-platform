import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  StatusBar,
  StyleSheet,
  Switch,
  ScrollView,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useCatalogo } from '../../hooks/useCatalogo';
import { CatalogItem, CatalogType, CATALOGS_CONFIG } from '../../dto/catalogo.dto';

interface CatalogoFormScreenProps {
  type: CatalogType;
  item?: CatalogItem; // If present, edit mode
  onBack: () => void;
  onSuccess: () => void;
}

export function CatalogoFormScreen({ type, item, onBack, onSuccess }: CatalogoFormScreenProps) {
  const config = CATALOGS_CONFIG[type];
  const { addItem, updateItem } = useCatalogo(type);
  const [formState, setFormState] = useState<Record<string, any>>({});
  const [isSaving, setIsSaving] = useState(false);

  // Initialize form state
  useEffect(() => {
    const initialState: Record<string, any> = {};
    config.fields.forEach((field) => {
      if (item) {
        initialState[field.name] = (item as any)[field.name];
      } else {
        // Default values
        initialState[field.name] = field.type === 'boolean' ? true : '';
      }
    });
    setFormState(initialState);
  }, [type, item]);

  const handleChange = (name: string, value: any) => {
    setFormState((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    // Validate required fields
    for (const field of config.fields) {
      if (field.required && !String(formState[field.name] || '').trim()) {
        Alert.alert('Oops', `El campo "${field.label}" es obligatorio.`);
        return;
      }
    }

    setIsSaving(true);
    let result;

    const payload = { ...formState };
    // Trim string inputs
    Object.keys(payload).forEach((key) => {
      if (typeof payload[key] === 'string') {
        payload[key] = payload[key].trim();
      }
    });

    if (item) {
      result = await updateItem(item.id, payload);
    } else {
      result = await addItem(payload);
    }

    setIsSaving(false);

    if (result.success) {
      onSuccess();
    } else {
      Alert.alert(
        item ? 'Error al actualizar' : 'Error al crear',
        result.error || 'Ocurrió un error al guardar los cambios'
      );
    }
  };

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={20} color="#8FA3A9" />
        </TouchableOpacity>
        <Text style={styles.title}>{item ? `Editar` : `Agregar`} {config.displayName}</Text>
      </View>

      <ScrollView contentContainerStyle={styles.formContent} keyboardShouldPersistTaps="handled">
        {config.fields.map((field) => {
          if (field.type === 'boolean') {
            return (
              <View key={field.name} style={styles.switchContainer}>
                <Text style={styles.switchLabel}>{field.label}</Text>
                <Switch
                  value={!!formState[field.name]}
                  onValueChange={(val) => handleChange(field.name, val)}
                  trackColor={{ false: '#112226', true: '#14B8A6' }}
                  thumbColor={formState[field.name] ? '#061114' : '#8FA3A9'}
                />
              </View>
            );
          }

          const isTextArea = field.type === 'textarea';
          return (
            <View key={field.name} style={[styles.inputContainer, isTextArea && styles.textAreaContainer]}>
              <MaterialIcons
                name={isTextArea ? 'description' : 'edit'}
                size={20}
                color="#8FA3A9"
                style={[styles.icon, isTextArea && styles.textAreaIcon]}
              />
              <TextInput
                style={[styles.input, isTextArea && styles.textAreaInput]}
                placeholder={field.label}
                placeholderTextColor="#8FA3A9"
                value={formState[field.name] ? String(formState[field.name]) : ''}
                onChangeText={(val) => handleChange(field.name, val)}
                multiline={isTextArea}
                numberOfLines={isTextArea ? 4 : 1}
              />
            </View>
          );
        })}

        <TouchableOpacity
          style={styles.primaryButton}
          onPress={handleSave}
          activeOpacity={0.8}
          disabled={isSaving}
        >
          {isSaving ? (
            <ActivityIndicator color="#061114" />
          ) : (
            <Text style={styles.primaryButtonText}>
              {item ? 'Guardar Cambios' : 'Crear Registro'}
            </Text>
          )}
        </TouchableOpacity>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#061114',
    paddingTop: 50,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 24,
    marginBottom: 32,
  },
  backButton: {
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
    fontSize: 22,
    fontWeight: '800',
    color: '#FFFFFF',
    letterSpacing: -0.5,
  },
  formContent: {
    paddingHorizontal: 24,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    borderRadius: 16,
    marginBottom: 16,
    paddingHorizontal: 16,
    height: 60,
  },
  textAreaContainer: {
    height: 120,
    alignItems: 'flex-start',
    paddingTop: 16,
  },
  icon: {
    marginRight: 12,
  },
  textAreaIcon: {
    marginTop: 2,
  },
  input: {
    flex: 1,
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  textAreaInput: {
    height: '100%',
    textAlignVertical: 'top',
  },
  switchContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    borderRadius: 16,
    marginBottom: 16,
    paddingHorizontal: 16,
    height: 60,
  },
  switchLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  primaryButton: {
    backgroundColor: '#14B8A6',
    height: 60,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 16,
  },
  primaryButtonText: {
    color: '#061114',
    fontSize: 16,
    fontWeight: 'bold',
  },
});
