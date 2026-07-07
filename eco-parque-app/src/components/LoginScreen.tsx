import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, Alert, ActivityIndicator, KeyboardAvoidingView, Platform, StatusBar } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { styles } from './LoginScreen.styles';

interface LoginScreenProps {
  onLoginSuccess: () => void;
  onGoToRegister: () => void;
  login: (email: string, pass: string) => Promise<{ success: boolean; error?: string }>;
  loading: boolean;
}

export function LoginScreen({ onLoginSuccess, onGoToRegister, login, loading }: LoginScreenProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert('Oops', 'Por favor ingresa email y contraseña');
      return;
    }
    const result = await login(email, password);
    if (result.success) {
      onLoginSuccess();
    } else {
      Alert.alert('Login Fallido', result.error || 'Credenciales incorrectas');
    }
  };

  return (
    <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#061114" />
      <View style={styles.card}>
        <View style={styles.header}>
          <View style={styles.iconContainer}>
             <MaterialIcons name="eco" size={42} color="#14B8A6" />
          </View>
          <Text style={styles.title}>EcoParque</Text>
          <Text style={styles.subtitle}>Ingresa a tu cuenta para continuar</Text>
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
          <TouchableOpacity style={styles.primaryButton} onPress={handleLogin} activeOpacity={0.8}>
            <Text style={styles.primaryButtonText}>Iniciar Sesión</Text>
          </TouchableOpacity>
        )}

        <TouchableOpacity style={styles.secondaryButton} onPress={onGoToRegister} disabled={loading} activeOpacity={0.6}>
          <Text style={styles.secondaryButtonText}>¿No tienes cuenta? <Text style={styles.linkText}>Regístrate</Text></Text>
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
}


