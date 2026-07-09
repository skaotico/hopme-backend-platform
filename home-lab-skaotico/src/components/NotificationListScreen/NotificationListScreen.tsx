import React from 'react';
import { View, Text, FlatList, TouchableOpacity, StatusBar } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { styles } from './NotificationListScreen.styles';

interface NotificationListScreenProps {
  onBack: () => void;
}

// Interfaz temporal para notificaciones mockeadas
interface Notification {
  id: string;
  title: string;
  body: string;
  date: string;
  read: boolean;
}

export function NotificationListScreen({ onBack }: NotificationListScreenProps) {
  // Datos temporales (mock) para mostrar cómo se verán. 
  // En el futuro puedes reemplazar esto por un hook de notificaciones reales.
  const mockNotifications: Notification[] = [
    {
      id: '1',
      title: '¡Bienvenido al sistema!',
      body: 'Gracias por utilizar la plataforma Eco-Parque. Aquí recibirás alertas de tus sensores y recordatorios de riego.',
      date: 'Justo ahora',
      read: false,
    },
    {
      id: '2',
      title: 'Sensor desconectado',
      body: 'El sensor de humedad en el Árbol "Roble Viejo" no está enviando datos desde hace 2 horas.',
      date: 'Hace 2 horas',
      read: false,
    },
    {
      id: '3',
      title: 'Mantenimiento programado',
      body: 'Mañana a las 10:00 AM habrá un corte programado para actualización de servidores.',
      date: 'Ayer',
      read: true,
    }
  ];

  const renderItem = ({ item }: { item: Notification }) => (
    <TouchableOpacity style={[styles.card, item.read ? { opacity: 0.7 } : null]} activeOpacity={0.8}>
      <View style={styles.cardHeader}>
        <View style={styles.titleContainer}>
          {!item.read && <View style={styles.unreadIndicator} />}
          <Text style={styles.cardTitle} numberOfLines={1}>{item.title}</Text>
        </View>
        <Text style={styles.timeText}>{item.date}</Text>
      </View>
      <Text style={styles.cardBody}>{item.body}</Text>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#0A191E" />
      <View style={styles.header}>
        <TouchableOpacity style={styles.backButton} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={24} color="#8FA3A9" />
        </TouchableOpacity>
        <Text style={styles.title}>Notificaciones</Text>
      </View>

      <FlatList
        data={mockNotifications}
        keyExtractor={item => item.id}
        renderItem={renderItem}
        contentContainerStyle={styles.listContent}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIconCircle}>
              <MaterialIcons name="notifications-none" size={48} color="#14B8A6" />
            </View>
            <Text style={styles.emptyText}>No tienes notificaciones</Text>
            <Text style={styles.emptySubtext}>Aquí aparecerán tus alertas y mensajes importantes.</Text>
          </View>
        }
      />
    </View>
  );
}
