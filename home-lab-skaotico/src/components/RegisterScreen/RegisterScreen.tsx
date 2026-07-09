import React, { useState } from 'react';
import { View, Text, TextInput, StyleSheet, Alert, ActivityIndicator, TouchableOpacity, KeyboardAvoidingView, Platform, StatusBar } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { styles } from './RegisterScreen.styles';

interface RegisterScreenProps {
  onRegisterSuccess: () => void;
  onGoToLogin: () => void;
  register: (username: string, email: string, pass: string) => Promise<{ success: boolean; error?: string }>;
  loading: boolean;
}

export function RegisterScreen({ onRegisterSuccess, onGoToLogin, register, loading }: RegisterScreenProps) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleRegister = async () => {
    if (!username || !email || !password) {
      Alert.alert('Oops', 'Por favor completa todos los campos');
      return;
    }

    const result = await register(username, email, password);
    if (result.success) {
      Alert.alert('¡Bienvenido!', 'Tu cuenta ha sido creada exitosamente.');
      onRegisterSuccess();
    } else {
      Alert.alert('Registro Fallido', result.error || 'No se pudo crear el usuario');
    }
  };

  return (
    <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.card}>
        <View style={styles.header}>
          <View style={styles.iconContainer}>
            <MaterialIcons name="person-add" size={36} color="#14B8A6" />
          </View>
          <Text style={styles.title}>Crear Cuenta</Text>
          <Text style={styles.subtitle}>Únete a la plataforma de EcoParques</Text>
        </View>

        <View style={styles.inputContainer}>
          <MaterialIcons name="person" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Nombre de usuario"
            placeholderTextColor="#8FA3A9"
            autoCapitalize="words"
            value={username}
            onChangeText={setUsername}
          />
        </View>
        
        <View style={styles.inputContainer}>
          <MaterialIcons name="email" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Correo electrónico"
            placeholderTextColor="#8FA3A9"
            autoCapitalize="none"
            keyboardType="email-address"
            value={email}
            onChangeText={setEmail}
          />
        </View>
        
        <View style={styles.inputContainer}>
          <MaterialIcons name="lock" size={20} color="#8FA3A9" style={styles.icon} />
          <TextInput
            style={styles.input}
            placeholder="Contraseña"
            placeholderTextColor="#8FA3A9"
            secureTextEntry
            value={password}
            onChangeText={setPassword}
          />
        </View>
        
        {loading ? (
          <ActivityIndicator size="large" color="#14B8A6" style={{ marginTop: 20 }} />
        ) : (
          <TouchableOpacity style={styles.primaryButton} onPress={handleRegister} activeOpacity={0.8}>
            <Text style={styles.primaryButtonText}>Registrarme</Text>
          </TouchableOpacity>
        )}

        <TouchableOpacity style={styles.secondaryButton} onPress={onGoToLogin} disabled={loading} activeOpacity={0.6}>
          <Text style={styles.secondaryButtonText}>¿Ya tienes cuenta? <Text style={styles.linkText}>Inicia sesión</Text></Text>
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
}


