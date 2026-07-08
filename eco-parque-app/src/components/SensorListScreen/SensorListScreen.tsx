import React, { useEffect, useState } from 'react';
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  RefreshControl,
  StatusBar,
  StyleSheet,
  ScrollView,
  Dimensions,
} from 'react-native';
import { Feather } from '@expo/vector-icons';
import { LinearGradient } from 'expo-linear-gradient';
import { useSensores, useLecturas, useAlertas } from '../../hooks/useSensores';
import { Sensor, LecturaSensor } from '../../dto/sensor.dto';
import { LineChart } from 'react-native-chart-kit';
import { PressableScale, FadeSlideCard } from '../ui/AnimatedCards';

const { width: SCREEN_WIDTH } = Dimensions.get('window');

// ─── Colores Vibrantes ──────────────────────────────────────────────────────────
const COLORS = {
  bg: '#040d10',
  surface: '#0d2227',
  border: '#173b42',
  teal: '#2DD4BF',
  tealDim: 'rgba(45, 212, 191, 0.15)',
  tealBorder: 'rgba(45, 212, 191, 0.3)',
  amber: '#FBBF24',
  rose: '#FB7185',
  text: '#F8FAFC',
  muted: '#94A3B8',
  cardBg: '#0f262b',
  gradientSelected: ['#0f3535', '#082121'] as const,
};

type RangoTiempo = 'hoy' | 'semana' | 'mes';

// ─── Gráfico de humedad (react-native-chart-kit) ────────────────────────────
function HumedadChart({ lecturas }: { lecturas: LecturaSensor[] }) {
  if (lecturas.length === 0) {
    return (
      <View style={styles.chartEmpty}>
        <Feather name="activity" size={32} color={COLORS.muted} />
        <Text style={styles.chartEmptyText}>Sin lecturas en este periodo</Text>
      </View>
    );
  }

  // Reducir la cantidad de puntos si son muchos para no colapsar el gráfico
  const maxPoints = 20;
  const step = Math.max(1, Math.floor(lecturas.length / maxPoints));
  const sampledLecturas = lecturas.filter((_, i) => i % step === 0).slice(-maxPoints);

  const chartData = sampledLecturas.map(l => Math.round(l.valor * 10) / 10);
  const labels = sampledLecturas.map((_, i) => (i % 4 === 0 ? String(i) : ''));
  const data = chartData.length > 0 ? chartData : [0];

  return (
    <View style={{ width: '100%', alignItems: 'center' }}>
      <LineChart
        data={{
          labels: labels,
          datasets: [{ data: data }]
        }}
        width={SCREEN_WIDTH - 60}
        height={200}
        withDots={true}
        withInnerLines={false}
        withOuterLines={false}
        withVerticalLines={false}
        chartConfig={{
          backgroundColor: 'transparent',
          backgroundGradientFrom: COLORS.cardBg,
          backgroundGradientTo: COLORS.cardBg,
          backgroundGradientFromOpacity: 0,
          backgroundGradientToOpacity: 0,
          decimalPlaces: 1,
          color: (opacity = 1) => `rgba(45, 212, 191, ${opacity})`,
          labelColor: (opacity = 1) => `rgba(148, 163, 184, ${opacity})`,
          style: { borderRadius: 16 },
          propsForDots: { r: "4", strokeWidth: "2", stroke: COLORS.surface }
        }}
        bezier
        style={{ marginVertical: 8, borderRadius: 16, paddingRight: 16 }}
      />
      
      {/* Explicación del gráfico */}
      <View style={styles.chartLegendContainer}>
        <View style={styles.chartLegendRow}>
          <View style={[styles.legendDot, { backgroundColor: COLORS.teal }]} />
          <Text style={styles.chartLegendText}>Eje Y: Porcentaje de Humedad (%)</Text>
        </View>
        <View style={styles.chartLegendRow}>
          <View style={[styles.legendDot, { backgroundColor: COLORS.muted }]} />
          <Text style={styles.chartLegendText}>Eje X: Secuencia temporal de lecturas</Text>
        </View>
      </View>
    </View>
  );
}

// ─── Badge de estado sensor ─────────────────────────────────────────────────
function SensorBadge({ activo }: { activo: boolean }) {
  return (
    <View style={[styles.badge, activo ? styles.badgeActive : styles.badgeInactive]}>
      <View style={[styles.badgeDot, { backgroundColor: activo ? COLORS.teal : COLORS.rose }]} />
      <Text style={[styles.badgeText, { color: activo ? COLORS.teal : COLORS.rose }]}>
        {activo ? 'Online' : 'Offline'}
      </Text>
    </View>
  );
}

// ─── Card de sensor (Horizontal Carousel Item) ───────────────────────────────
function SensorCard({
  sensor,
  index,
  onPress,
  isSelected,
}: {
  sensor: Sensor;
  index: number;
  onPress: () => void;
  isSelected: boolean;
}) {
  const innerContent = (
    <>
      <View style={styles.sensorCardHeader}>
        <View style={[styles.sensorIcon, isSelected && { backgroundColor: COLORS.teal }]}>
          <Feather name="cpu" size={20} color={isSelected ? COLORS.bg : COLORS.teal} />
        </View>
        <SensorBadge activo={sensor.activo} />
      </View>
      <View style={{ marginTop: 12 }}>
        <Text style={styles.sensorCodigo} numberOfLines={1}>{sensor.codigo || 'S/C'}</Text>
        <Text style={styles.sensorModelo} numberOfLines={1}>
          {[sensor.fabricante, sensor.modelo].filter(Boolean).join(' · ') || 'Modelo desc.'}
        </Text>
      </View>
    </>
  );

  return (
    <FadeSlideCard delay={index * 80} style={{ marginRight: 12 }}>
      <PressableScale onPress={onPress}>
        {isSelected ? (
          <LinearGradient
            colors={COLORS.gradientSelected}
            start={{ x: 0, y: 0 }}
            end={{ x: 1, y: 1 }}
            style={[styles.sensorCard, styles.sensorCardSelected]}
          >
            {innerContent}
          </LinearGradient>
        ) : (
          <View style={styles.sensorCard}>{innerContent}</View>
        )}
      </PressableScale>
    </FadeSlideCard>
  );
}

// ─── Panel de detalle de sensor ──────────────────────────────────────────────
function SensorDetailPanel({ sensorId, codigo }: { sensorId: string, codigo: string }) {
  const { lecturas, loading, error, refresh } = useLecturas(sensorId);
  const { alertas, loading: loadingAlertas, refresh: refreshAlertas } = useAlertas(sensorId);
  
  const [rango, setRango] = useState<RangoTiempo>('semana');

  useEffect(() => {
    refresh();
    refreshAlertas();
  }, [sensorId]);

  const lastLectura = lecturas[lecturas.length - 1];
  const alertasActivas = alertas.filter((a) => a.estado === 'activa');

  // Filtrar lecturas según el rango
  const now = new Date();
  const lecturasFiltradas = lecturas.filter((l) => {
    const d = new Date(l.fecha_lectura);
    const diffDays = Math.abs(now.getTime() - d.getTime()) / (1000 * 60 * 60 * 24);
    if (rango === 'hoy') return diffDays <= 1;
    if (rango === 'semana') return diffDays <= 7;
    if (rango === 'mes') return diffDays <= 30;
    return true;
  });

  return (
    <ScrollView
      style={styles.detailPanel}
      showsVerticalScrollIndicator={false}
      contentContainerStyle={{ paddingBottom: 32 }}
    >
      <FadeSlideCard delay={50}>
        <View style={styles.panelHeader}>
          <Text style={styles.panelLabel}>MONITOREO EN VIVO</Text>
          <Text style={styles.panelSensorCode}>{codigo}</Text>
        </View>
      </FadeSlideCard>

      {/* KPIs */}
      <FadeSlideCard delay={100}>
        <View style={styles.kpiRow}>
          <View style={[styles.kpiCard, { borderColor: COLORS.tealBorder }]}>
            <Feather name="droplet" size={24} color={COLORS.teal} />
            <Text style={styles.kpiValue}>
              {lastLectura ? `${lastLectura.valor.toFixed(1)}%` : '—'}
            </Text>
            <Text style={styles.kpiLabel}>Humedad</Text>
          </View>
          <View style={[styles.kpiCard, { borderColor: 'rgba(251,191,36,0.3)' }]}>
            <Feather name="battery-charging" size={24} color={COLORS.amber} />
            <Text style={[styles.kpiValue, { color: COLORS.amber }]}>
              {lastLectura?.bateria_porcentaje != null ? `${lastLectura.bateria_porcentaje}%` : '—'}
            </Text>
            <Text style={styles.kpiLabel}>Batería</Text>
          </View>
          <View
            style={[
              styles.kpiCard,
              { borderColor: alertasActivas.length > 0 ? 'rgba(251,113,133,0.3)' : COLORS.border },
            ]}
          >
            <Feather
              name="bell"
              size={24}
              color={alertasActivas.length > 0 ? COLORS.rose : COLORS.muted}
            />
            <Text style={[styles.kpiValue, { color: alertasActivas.length > 0 ? COLORS.rose : COLORS.text }]}>
              {alertasActivas.length}
            </Text>
            <Text style={styles.kpiLabel}>Alertas</Text>
          </View>
        </View>
      </FadeSlideCard>

      {/* Gráfico humedad */}
      <FadeSlideCard delay={180}>
        <View style={styles.chartCard}>
          <View style={styles.chartHeaderRow}>
            <View style={styles.chartHeader}>
              <Feather name="bar-chart-2" size={18} color={COLORS.teal} />
              <Text style={styles.chartTitle}>Historial de Humedad</Text>
            </View>
            
            {/* Filtros */}
            <View style={styles.filterRow}>
              {(['hoy', 'semana', 'mes'] as RangoTiempo[]).map((r) => (
                <TouchableOpacity
                  key={r}
                  style={[styles.filterBtn, rango === r && styles.filterBtnActive]}
                  onPress={() => setRango(r)}
                >
                  <Text style={[styles.filterText, rango === r && styles.filterTextActive]}>
                    {r.charAt(0).toUpperCase() + r.slice(1)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {loading && lecturas.length === 0 ? (
            <ActivityIndicator color={COLORS.teal} style={{ paddingVertical: 32 }} />
          ) : (
            <HumedadChart lecturas={lecturasFiltradas} />
          )}
        </View>
      </FadeSlideCard>

      {/* Lecturas recientes */}
      <FadeSlideCard delay={250}>
        <View style={styles.sectionCard}>
          <Text style={styles.sectionTitle}>ÚLTIMAS LECTURAS</Text>
          {loading && lecturas.length === 0 ? (
            <ActivityIndicator color={COLORS.teal} />
          ) : lecturas.length === 0 ? (
            <Text style={styles.emptyText}>No hay datos aún.</Text>
          ) : (
            lecturas.slice().reverse().slice(0, 5).map((l) => (
                <View key={l.id} style={styles.lecturaRow}>
                  <View style={styles.lecturaLeft}>
                    <Text style={styles.lecturaValor}>{l.valor.toFixed(1)}</Text>
                    <Text style={styles.lecturaUnidad}>%</Text>
                  </View>
                  <View>
                    <Text style={styles.lecturaFecha}>
                      {new Date(l.fecha_lectura).toLocaleString('es-CL', {
                        dateStyle: 'short', timeStyle: 'short',
                      })}
                    </Text>
                    {l.observacion ? <Text style={styles.lecturaObs}>{l.observacion}</Text> : null}
                  </View>
                </View>
              ))
          )}
        </View>
      </FadeSlideCard>
    </ScrollView>
  );
}

// ─── Pantalla principal ─────────────────────────────────────────────────────
interface SensorListScreenProps {
  arbolId: string;
  arbolCodigo?: string;
  onBack: () => void;
}

export function SensorListScreen({ arbolId, arbolCodigo, onBack }: SensorListScreenProps) {
  const { sensores, loading, error, refresh } = useSensores(arbolId);
  const [selectedSensor, setSelectedSensor] = useState<Sensor | null>(null);

  useEffect(() => { refresh(); }, [refresh]);
  useEffect(() => {
    if (sensores.length > 0 && !selectedSensor) {
      setSelectedSensor(sensores[0]);
    }
  }, [sensores]);

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor={COLORS.bg} />

      <View style={styles.header}>
        <TouchableOpacity style={styles.backBtn} onPress={onBack} activeOpacity={0.7}>
          <Feather name="chevron-left" size={24} color={COLORS.muted} />
        </TouchableOpacity>
        <View style={{ flex: 1 }}>
          <Text style={styles.headerTitle}>Dispositivos</Text>
          {arbolCodigo && <Text style={styles.headerSub}>Árbol: {arbolCodigo}</Text>}
        </View>
        <View style={styles.headerCount}>
          <Text style={styles.headerCountText}>{sensores.length} online</Text>
        </View>
      </View>

      {error && (
        <View style={styles.errorBanner}>
          <Feather name="alert-circle" size={16} color={COLORS.rose} />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      {loading && sensores.length === 0 ? (
        <View style={styles.center}>
          <ActivityIndicator size="large" color={COLORS.teal} />
        </View>
      ) : (
        <View style={styles.content}>
          <View style={styles.carouselContainer}>
            <FlatList
              horizontal
              data={sensores}
              keyExtractor={(s) => s.id}
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={styles.carouselContent}
              refreshControl={<RefreshControl refreshing={loading} onRefresh={refresh} tintColor={COLORS.teal} />}
              ListEmptyComponent={
                <View style={styles.emptyContainer}>
                  <Feather name="radio" size={40} color={COLORS.muted} />
                  <Text style={styles.emptyTitle}>Sin sensores</Text>
                  <Text style={styles.emptyDesc}>No hay dispositivos vinculados a este árbol.</Text>
                </View>
              }
              renderItem={({ item, index }) => (
                <SensorCard
                  sensor={item}
                  index={index}
                  isSelected={selectedSensor?.id === item.id}
                  onPress={() => setSelectedSensor(item)}
                />
              )}
            />
          </View>

          <View style={styles.detailContainer}>
            {selectedSensor ? (
              <SensorDetailPanel sensorId={selectedSensor.id} codigo={selectedSensor.codigo || 'S/C'} />
            ) : null}
          </View>
        </View>
      )}
    </View>
  );
}

// ─── Estilos ─────────────────────────────────────────────────────────────────
const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: COLORS.bg },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: COLORS.bg },

  header: {
    flexDirection: 'row', alignItems: 'center', paddingTop: 56, paddingHorizontal: 20,
    paddingBottom: 16, backgroundColor: COLORS.surface, borderBottomWidth: 1,
    borderBottomColor: COLORS.border, gap: 16,
  },
  backBtn: {
    width: 40, height: 40, borderRadius: 12, backgroundColor: COLORS.cardBg,
    justifyContent: 'center', alignItems: 'center', borderWidth: 1, borderColor: COLORS.border,
  },
  headerTitle: { fontSize: 20, fontWeight: '800', color: COLORS.text, letterSpacing: -0.5 },
  headerSub: { fontSize: 13, color: COLORS.teal, marginTop: 2, fontWeight: '600' },
  headerCount: {
    backgroundColor: COLORS.tealDim, borderWidth: 1, borderColor: COLORS.tealBorder,
    borderRadius: 20, paddingHorizontal: 12, paddingVertical: 6,
  },
  headerCountText: { fontSize: 12, fontWeight: '700', color: COLORS.teal },

  errorBanner: {
    flexDirection: 'row', alignItems: 'center', backgroundColor: 'rgba(244,63,94,0.1)',
    paddingHorizontal: 20, paddingVertical: 12, gap: 10,
  },
  errorText: { color: COLORS.rose, fontSize: 13, fontWeight: '600' },

  content: { flex: 1, flexDirection: 'column' },
  carouselContainer: {
    paddingVertical: 16, borderBottomWidth: 1, borderBottomColor: COLORS.border,
    backgroundColor: COLORS.surface,
  },
  carouselContent: { paddingHorizontal: 20 },

  detailContainer: { flex: 1, paddingHorizontal: 16, paddingTop: 20 },

  panelHeader: { marginBottom: 16, paddingHorizontal: 4 },
  panelLabel: { fontSize: 11, fontWeight: '800', color: COLORS.teal, letterSpacing: 1.5, marginBottom: 4 },
  panelSensorCode: { fontSize: 24, fontWeight: '800', color: COLORS.text, letterSpacing: -0.5 },

  sensorCard: {
    backgroundColor: COLORS.cardBg, borderRadius: 16, padding: 16, borderWidth: 1,
    borderColor: COLORS.border, width: 160, height: 110, shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 }, shadowOpacity: 0.2, shadowRadius: 8, elevation: 4,
  },
  sensorCardSelected: { borderColor: COLORS.teal, borderWidth: 1.5, shadowColor: COLORS.teal, shadowOpacity: 0.3, shadowRadius: 12 },
  sensorCardHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'flex-start' },
  sensorIcon: { width: 36, height: 36, borderRadius: 10, backgroundColor: COLORS.tealDim, justifyContent: 'center', alignItems: 'center' },
  sensorCodigo: { fontSize: 15, fontWeight: '800', color: COLORS.text, letterSpacing: -0.3 },
  sensorModelo: { fontSize: 12, color: COLORS.muted, marginTop: 4 },

  badge: { flexDirection: 'row', alignItems: 'center', paddingHorizontal: 8, paddingVertical: 4, borderRadius: 20, borderWidth: 1, gap: 6 },
  badgeActive: { backgroundColor: COLORS.tealDim, borderColor: COLORS.tealBorder },
  badgeInactive: { backgroundColor: 'rgba(244,63,94,0.1)', borderColor: 'rgba(244,63,94,0.25)' },
  badgeDot: { width: 6, height: 6, borderRadius: 3 },
  badgeText: { fontSize: 10, fontWeight: '800' },

  kpiRow: { flexDirection: 'row', gap: 12, marginBottom: 16 },
  kpiCard: {
    flex: 1, backgroundColor: COLORS.cardBg, borderRadius: 16, borderWidth: 1, padding: 16,
    alignItems: 'center', gap: 8, shadowColor: '#000', shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1, shadowRadius: 4, elevation: 2,
  },
  kpiValue: { fontSize: 20, fontWeight: '800', color: COLORS.text, letterSpacing: -0.5 },
  kpiLabel: { fontSize: 11, color: COLORS.muted, fontWeight: '600' },

  chartCard: { backgroundColor: COLORS.cardBg, borderRadius: 16, borderWidth: 1, borderColor: COLORS.border, padding: 16, marginBottom: 16 },
  chartHeaderRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 },
  chartHeader: { flexDirection: 'row', alignItems: 'center', gap: 8 },
  chartTitle: { fontSize: 15, fontWeight: '700', color: COLORS.text },
  
  filterRow: { flexDirection: 'row', gap: 4, backgroundColor: COLORS.surface, borderRadius: 8, padding: 4, borderWidth: 1, borderColor: COLORS.border },
  filterBtn: { paddingHorizontal: 8, paddingVertical: 4, borderRadius: 6 },
  filterBtnActive: { backgroundColor: COLORS.tealDim },
  filterText: { fontSize: 10, color: COLORS.muted, fontWeight: '600' },
  filterTextActive: { color: COLORS.teal, fontWeight: '800' },

  chartLegendContainer: { marginTop: 12, paddingHorizontal: 8, gap: 6, width: '100%', alignItems: 'flex-start' },
  chartLegendRow: { flexDirection: 'row', alignItems: 'center', gap: 8 },
  legendDot: { width: 8, height: 8, borderRadius: 4 },
  chartLegendText: { fontSize: 11, color: COLORS.muted, fontWeight: '500' },

  chartEmpty: { height: 160, justifyContent: 'center', alignItems: 'center', gap: 12 },
  chartEmptyText: { fontSize: 13, color: COLORS.muted },

  sectionCard: { backgroundColor: COLORS.cardBg, borderRadius: 16, borderWidth: 1, borderColor: COLORS.border, padding: 16, marginBottom: 16 },
  sectionTitle: { fontSize: 11, fontWeight: '800', color: COLORS.muted, letterSpacing: 1.5, marginBottom: 16 },
  emptyText: { fontSize: 13, color: COLORS.muted, textAlign: 'center', paddingVertical: 16 },

  lecturaRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', paddingVertical: 12, borderBottomWidth: 1, borderBottomColor: COLORS.border },
  lecturaLeft: { flexDirection: 'row', alignItems: 'baseline', gap: 4 },
  lecturaValor: { fontSize: 18, fontWeight: '800', color: COLORS.teal },
  lecturaUnidad: { fontSize: 12, color: COLORS.muted, fontWeight: '700' },
  lecturaFecha: { fontSize: 12, color: COLORS.text, textAlign: 'right', fontWeight: '500' },
  lecturaObs: { fontSize: 11, color: COLORS.muted, marginTop: 4, textAlign: 'right' },

  emptyContainer: { paddingVertical: 40, alignItems: 'center', gap: 12 },
  emptyTitle: { fontSize: 16, fontWeight: '700', color: COLORS.muted },
  emptyDesc: { fontSize: 13, color: COLORS.muted, textAlign: 'center' },

  detailPanel: { flex: 1 },
});
