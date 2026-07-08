import React, { useEffect, useRef, useState } from 'react';
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  RefreshControl,
  StatusBar,
  StyleSheet,
  Animated,
  ScrollView,
  Dimensions,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { useSensores, useLecturas, useAlertas } from '../hooks/useSensores';
import { Sensor, LecturaSensor } from '../dto/sensor.dto';

const { width: SCREEN_WIDTH } = Dimensions.get('window');

// ─── Colores ───────────────────────────────────────────────────────────────
const COLORS = {
  bg: '#061114',
  surface: '#0C1E22',
  border: '#1D343B',
  teal: '#14B8A6',
  tealDim: 'rgba(20,184,166,0.12)',
  tealBorder: 'rgba(20,184,166,0.25)',
  amber: '#F59E0B',
  rose: '#F43F5E',
  text: '#FFFFFF',
  muted: '#8FA3A9',
  cardBg: '#0E2026',
};

// ─── Animated Card genérica ─────────────────────────────────────────────────
function FadeSlideCard({
  children,
  delay = 0,
  style,
}: {
  children: React.ReactNode;
  delay?: number;
  style?: object;
}) {
  const anim = useRef(new Animated.Value(0)).current;
  useEffect(() => {
    Animated.timing(anim, {
      toValue: 1,
      duration: 400,
      delay,
      useNativeDriver: true,
    }).start();
  }, []);
  return (
    <Animated.View
      style={[
        {
          opacity: anim,
          transform: [
            {
              translateY: anim.interpolate({
                inputRange: [0, 1],
                outputRange: [24, 0],
              }),
            },
          ],
        },
        style,
      ]}
    >
      {children}
    </Animated.View>
  );
}

// ─── Gráfico de humedad (Victory Native XL) ─────────────────────────────────
function HumedadChart({ lecturas }: { lecturas: LecturaSensor[] }) {
  // Intentamos importar Victory Native dinámicamente para evitar error de build si
  // las deps no están instaladas todavía.
  let CartesianChart: any = null;
  let Line: any = null;
  let Area: any = null;
  let useFont: any = null;

  try {
    const vn = require('victory-native');
    CartesianChart = vn.CartesianChart;
    Line = vn.Line;
    Area = vn.Area;
    useFont = vn.useFont;
  } catch {
    // Victory Native no disponible aún
  }

  if (!CartesianChart || lecturas.length === 0) {
    return (
      <View style={styles.chartEmpty}>
        <MaterialIcons name="show-chart" size={32} color={COLORS.muted} />
        <Text style={styles.chartEmptyText}>Sin lecturas disponibles</Text>
      </View>
    );
  }

  // Preparamos datos: índice vs valor
  const chartData = lecturas.slice(-20).map((l, i) => ({
    index: i,
    valor: Math.round(l.valor * 10) / 10,
  }));

  // Rango del eje Y
  const values = chartData.map((d) => d.valor);
  const minVal = Math.max(0, Math.floor(Math.min(...values) - 5));
  const maxVal = Math.ceil(Math.max(...values) + 5);

  return (
    <View style={{ height: 220, width: '100%' }}>
      <CartesianChart
        data={chartData}
        xKey="index"
        yKeys={['valor']}
        domain={{ y: [minVal, maxVal] }}
        axisOptions={{
          tickCount: { x: 5, y: 4 },
          labelColor: COLORS.muted,
          lineColor: COLORS.border,
          labelOffset: { x: 2, y: 4 },
        }}
        chartPressState={undefined}
      >
        {({ points, chartBounds }: any) => (
          <>
            {/* Area fill */}
            <Area
              points={points.valor}
              y0={chartBounds.bottom}
              color={COLORS.teal}
              opacity={0.15}
              animate={{ type: 'timing', duration: 600 }}
              curveType="natural"
            />
            {/* Line */}
            <Line
              points={points.valor}
              color={COLORS.teal}
              strokeWidth={2.5}
              animate={{ type: 'timing', duration: 600 }}
              curveType="natural"
            />
          </>
        )}
      </CartesianChart>
    </View>
  );
}

// ─── Badge de estado sensor ─────────────────────────────────────────────────
function SensorBadge({ activo }: { activo: boolean }) {
  return (
    <View
      style={[
        styles.badge,
        activo ? styles.badgeActive : styles.badgeInactive,
      ]}
    >
      <View
        style={[
          styles.badgeDot,
          { backgroundColor: activo ? COLORS.teal : COLORS.rose },
        ]}
      />
      <Text
        style={[
          styles.badgeText,
          { color: activo ? COLORS.teal : COLORS.rose },
        ]}
      >
        {activo ? 'Activo' : 'Inactivo'}
      </Text>
    </View>
  );
}

// ─── Card de sensor ──────────────────────────────────────────────────────────
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
  return (
    <FadeSlideCard delay={index * 80} style={{ marginBottom: 12 }}>
      <TouchableOpacity
        style={[styles.sensorCard, isSelected && styles.sensorCardSelected]}
        onPress={onPress}
        activeOpacity={0.8}
      >
        <View style={styles.sensorCardLeft}>
          <View
            style={[
              styles.sensorIcon,
              isSelected && { backgroundColor: COLORS.teal },
            ]}
          >
            <MaterialIcons
              name="sensors"
              size={20}
              color={isSelected ? COLORS.bg : COLORS.teal}
            />
          </View>
          <View style={{ flex: 1, marginLeft: 12 }}>
            <Text style={styles.sensorCodigo}>{sensor.codigo || 'S/C'}</Text>
            <Text style={styles.sensorModelo}>
              {[sensor.fabricante, sensor.modelo].filter(Boolean).join(' · ') ||
                'Modelo desconocido'}
            </Text>
            {sensor.fecha_instalacion && (
              <Text style={styles.sensorFecha}>
                Instalado:{' '}
                {new Date(sensor.fecha_instalacion).toLocaleDateString('es-CL')}
              </Text>
            )}
          </View>
        </View>
        <SensorBadge activo={sensor.activo} />
      </TouchableOpacity>
    </FadeSlideCard>
  );
}

// ─── Panel de detalle de sensor ──────────────────────────────────────────────
function SensorDetailPanel({ sensorId }: { sensorId: string }) {
  const { lecturas, loading, error, refresh } = useLecturas(sensorId);
  const {
    alertas,
    loading: loadingAlertas,
    refresh: refreshAlertas,
  } = useAlertas(sensorId);

  useEffect(() => {
    refresh();
    refreshAlertas();
  }, [sensorId]);

  const lastLectura = lecturas[lecturas.length - 1];
  const alertasActivas = alertas.filter((a) => a.estado === 'activa');

  return (
    <ScrollView
      style={styles.detailPanel}
      showsVerticalScrollIndicator={false}
      contentContainerStyle={{ paddingBottom: 24 }}
    >
      {/* KPIs */}
      <FadeSlideCard delay={50}>
        <View style={styles.kpiRow}>
          <View style={[styles.kpiCard, { borderColor: COLORS.tealBorder }]}>
            <MaterialIcons name="water-drop" size={22} color={COLORS.teal} />
            <Text style={styles.kpiValue}>
              {lastLectura ? `${lastLectura.valor.toFixed(1)}%` : '—'}
            </Text>
            <Text style={styles.kpiLabel}>Último valor</Text>
          </View>
          <View style={[styles.kpiCard, { borderColor: 'rgba(245,158,11,0.3)' }]}>
            <MaterialIcons name="battery-std" size={22} color={COLORS.amber} />
            <Text style={[styles.kpiValue, { color: COLORS.amber }]}>
              {lastLectura?.bateria_porcentaje != null
                ? `${lastLectura.bateria_porcentaje}%`
                : '—'}
            </Text>
            <Text style={styles.kpiLabel}>Batería</Text>
          </View>
          <View
            style={[
              styles.kpiCard,
              { borderColor: alertasActivas.length > 0 ? 'rgba(244,63,94,0.3)' : COLORS.border },
            ]}
          >
            <MaterialIcons
              name="notifications-active"
              size={22}
              color={alertasActivas.length > 0 ? COLORS.rose : COLORS.muted}
            />
            <Text
              style={[
                styles.kpiValue,
                { color: alertasActivas.length > 0 ? COLORS.rose : COLORS.muted },
              ]}
            >
              {alertasActivas.length}
            </Text>
            <Text style={styles.kpiLabel}>Alertas</Text>
          </View>
        </View>
      </FadeSlideCard>

      {/* Gráfico humedad */}
      <FadeSlideCard delay={120}>
        <View style={styles.chartCard}>
          <View style={styles.chartHeader}>
            <MaterialIcons name="show-chart" size={18} color={COLORS.teal} />
            <Text style={styles.chartTitle}>Historial de Humedad</Text>
            <Text style={styles.chartSubtitle}>
              {lecturas.length} lectura{lecturas.length !== 1 ? 's' : ''}
            </Text>
          </View>
          {loading ? (
            <ActivityIndicator color={COLORS.teal} style={{ paddingVertical: 32 }} />
          ) : (
            <HumedadChart lecturas={lecturas} />
          )}
        </View>
      </FadeSlideCard>

      {/* Lecturas recientes */}
      <FadeSlideCard delay={200}>
        <View style={styles.sectionCard}>
          <Text style={styles.sectionTitle}>LECTURAS RECIENTES</Text>
          {loading && lecturas.length === 0 ? (
            <ActivityIndicator color={COLORS.teal} />
          ) : lecturas.length === 0 ? (
            <Text style={styles.emptyText}>Sin lecturas</Text>
          ) : (
            lecturas
              .slice()
              .reverse()
              .slice(0, 8)
              .map((l, i) => (
                <View key={l.id} style={styles.lecturaRow}>
                  <View style={styles.lecturaLeft}>
                    <Text style={styles.lecturaValor}>{l.valor.toFixed(2)}</Text>
                    <Text style={styles.lecturaUnidad}>%</Text>
                  </View>
                  <View>
                    <Text style={styles.lecturaFecha}>
                      {new Date(l.fecha_lectura).toLocaleString('es-CL', {
                        dateStyle: 'short',
                        timeStyle: 'short',
                      })}
                    </Text>
                    {l.observacion ? (
                      <Text style={styles.lecturaObs}>{l.observacion}</Text>
                    ) : null}
                  </View>
                </View>
              ))
          )}
        </View>
      </FadeSlideCard>

      {/* Alertas */}
      {alertas.length > 0 && (
        <FadeSlideCard delay={280}>
          <View style={styles.sectionCard}>
            <Text style={styles.sectionTitle}>ALERTAS</Text>
            {alertas.slice(0, 5).map((a) => (
              <View key={a.id} style={styles.alertaRow}>
                <MaterialIcons
                  name={a.estado === 'activa' ? 'warning-amber' : 'check-circle'}
                  size={18}
                  color={a.estado === 'activa' ? COLORS.rose : COLORS.teal}
                />
                <View style={{ flex: 1, marginLeft: 10 }}>
                  <Text style={styles.alertaTipo}>{a.tipo}</Text>
                  <Text style={styles.alertaDetalle}>
                    Detectado: {a.valor_detectado} · Umbral: {a.umbral}
                  </Text>
                </View>
                <View
                  style={[
                    styles.alertaEstadoBadge,
                    {
                      backgroundColor:
                        a.estado === 'activa'
                          ? 'rgba(244,63,94,0.12)'
                          : COLORS.tealDim,
                    },
                  ]}
                >
                  <Text
                    style={{
                      fontSize: 11,
                      fontWeight: '700',
                      color: a.estado === 'activa' ? COLORS.rose : COLORS.teal,
                    }}
                  >
                    {a.estado.toUpperCase()}
                  </Text>
                </View>
              </View>
            ))}
          </View>
        </FadeSlideCard>
      )}
    </ScrollView>
  );
}

// ─── Pantalla principal ─────────────────────────────────────────────────────
interface SensorListScreenProps {
  arbolId: string;
  arbolCodigo?: string;
  onBack: () => void;
}

export function SensorListScreen({
  arbolId,
  arbolCodigo,
  onBack,
}: SensorListScreenProps) {
  const { sensores, loading, error, refresh } = useSensores(arbolId);
  const [selectedSensor, setSelectedSensor] = useState<Sensor | null>(null);

  useEffect(() => {
    refresh();
  }, [refresh]);

  // Auto-select first sensor when list loads
  useEffect(() => {
    if (sensores.length > 0 && !selectedSensor) {
      setSelectedSensor(sensores[0]);
    }
  }, [sensores]);

  if (loading && sensores.length === 0) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={COLORS.teal} />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor={COLORS.bg} />

      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity style={styles.backBtn} onPress={onBack} activeOpacity={0.7}>
          <MaterialIcons name="arrow-back" size={20} color={COLORS.muted} />
        </TouchableOpacity>
        <View>
          <Text style={styles.headerTitle}>Sensores IoT</Text>
          {arbolCodigo && (
            <Text style={styles.headerSub}>Árbol: {arbolCodigo}</Text>
          )}
        </View>
        <View style={styles.headerCount}>
          <Text style={styles.headerCountText}>{sensores.length}</Text>
        </View>
      </View>

      {error && (
        <View style={styles.errorBanner}>
          <MaterialIcons name="error-outline" size={16} color={COLORS.rose} />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      <View style={styles.content}>
        {/* Lista de sensores */}
        <View style={styles.leftPanel}>
          <Text style={styles.panelLabel}>DISPOSITIVOS</Text>
          <FlatList
            data={sensores}
            keyExtractor={(s) => s.id}
            showsVerticalScrollIndicator={false}
            refreshControl={
              <RefreshControl
                refreshing={loading}
                onRefresh={refresh}
                tintColor={COLORS.teal}
              />
            }
            ListEmptyComponent={
              <View style={styles.emptyContainer}>
                <MaterialIcons name="sensors-off" size={40} color={COLORS.muted} />
                <Text style={styles.emptyTitle}>Sin sensores</Text>
                <Text style={styles.emptyDesc}>
                  No hay sensores registrados para este árbol.
                </Text>
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

        {/* Panel de detalle */}
        {selectedSensor ? (
          <View style={styles.rightPanel}>
            <Text style={styles.panelLabel}>MONITOREO</Text>
            <Text style={styles.panelSensorCode}>
              {selectedSensor.codigo || 'Sensor'}
            </Text>
            <SensorDetailPanel sensorId={selectedSensor.id} />
          </View>
        ) : (
          <View style={[styles.rightPanel, styles.center]}>
            <MaterialIcons name="touch-app" size={36} color={COLORS.muted} />
            <Text style={{ color: COLORS.muted, marginTop: 8, textAlign: 'center' }}>
              Selecciona un sensor para ver su monitoreo
            </Text>
          </View>
        )}
      </View>
    </View>
  );
}

// ─── Estilos ─────────────────────────────────────────────────────────────────
const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: COLORS.bg },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: COLORS.bg },

  // Header
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingTop: 52,
    paddingHorizontal: 20,
    paddingBottom: 16,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
    gap: 12,
  },
  backBtn: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: COLORS.surface,
    justifyContent: 'center',
    alignItems: 'center',
  },
  headerTitle: { fontSize: 18, fontWeight: '800', color: COLORS.text, letterSpacing: -0.3 },
  headerSub: { fontSize: 12, color: COLORS.muted, marginTop: 2 },
  headerCount: {
    marginLeft: 'auto',
    backgroundColor: COLORS.tealDim,
    borderWidth: 1,
    borderColor: COLORS.tealBorder,
    borderRadius: 20,
    paddingHorizontal: 10,
    paddingVertical: 4,
  },
  headerCountText: { fontSize: 13, fontWeight: '700', color: COLORS.teal },

  errorBanner: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(244,63,94,0.1)',
    paddingHorizontal: 16,
    paddingVertical: 10,
    gap: 8,
  },
  errorText: { color: COLORS.rose, fontSize: 13 },

  // Layout
  content: { flex: 1, flexDirection: 'row' },
  leftPanel: {
    width: SCREEN_WIDTH * 0.42,
    borderRightWidth: 1,
    borderRightColor: COLORS.border,
    paddingHorizontal: 12,
    paddingTop: 16,
  },
  rightPanel: {
    flex: 1,
    paddingHorizontal: 14,
    paddingTop: 16,
  },
  panelLabel: {
    fontSize: 10,
    fontWeight: '700',
    color: COLORS.muted,
    letterSpacing: 1.4,
    marginBottom: 12,
  },
  panelSensorCode: {
    fontSize: 15,
    fontWeight: '800',
    color: COLORS.text,
    marginBottom: 12,
    letterSpacing: -0.2,
  },

  // Sensor card
  sensorCard: {
    backgroundColor: COLORS.cardBg,
    borderRadius: 14,
    padding: 12,
    borderWidth: 1,
    borderColor: COLORS.border,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  sensorCardSelected: {
    borderColor: COLORS.teal,
    backgroundColor: 'rgba(14,32,38,0.9)',
  },
  sensorCardLeft: { flexDirection: 'row', alignItems: 'center', flex: 1 },
  sensorIcon: {
    width: 38,
    height: 38,
    borderRadius: 10,
    backgroundColor: COLORS.tealDim,
    justifyContent: 'center',
    alignItems: 'center',
  },
  sensorCodigo: { fontSize: 13, fontWeight: '700', color: COLORS.text },
  sensorModelo: { fontSize: 11, color: COLORS.muted, marginTop: 2 },
  sensorFecha: { fontSize: 10, color: COLORS.muted, marginTop: 2 },

  // Badge
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 20,
    borderWidth: 1,
    gap: 5,
  },
  badgeActive: {
    backgroundColor: COLORS.tealDim,
    borderColor: COLORS.tealBorder,
  },
  badgeInactive: {
    backgroundColor: 'rgba(244,63,94,0.1)',
    borderColor: 'rgba(244,63,94,0.25)',
  },
  badgeDot: { width: 6, height: 6, borderRadius: 3 },
  badgeText: { fontSize: 10, fontWeight: '700' },

  // KPIs
  kpiRow: {
    flexDirection: 'row',
    gap: 8,
    marginBottom: 14,
  },
  kpiCard: {
    flex: 1,
    backgroundColor: COLORS.cardBg,
    borderRadius: 12,
    borderWidth: 1,
    padding: 10,
    alignItems: 'center',
    gap: 4,
  },
  kpiValue: {
    fontSize: 16,
    fontWeight: '800',
    color: COLORS.text,
    letterSpacing: -0.5,
  },
  kpiLabel: { fontSize: 10, color: COLORS.muted, textAlign: 'center' },

  // Chart
  chartCard: {
    backgroundColor: COLORS.cardBg,
    borderRadius: 14,
    borderWidth: 1,
    borderColor: COLORS.border,
    padding: 14,
    marginBottom: 14,
  },
  chartHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
    gap: 6,
  },
  chartTitle: { fontSize: 13, fontWeight: '700', color: COLORS.text, flex: 1 },
  chartSubtitle: { fontSize: 11, color: COLORS.muted },
  chartEmpty: {
    height: 120,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
  },
  chartEmptyText: { fontSize: 12, color: COLORS.muted },

  // Sections
  sectionCard: {
    backgroundColor: COLORS.cardBg,
    borderRadius: 14,
    borderWidth: 1,
    borderColor: COLORS.border,
    padding: 14,
    marginBottom: 12,
  },
  sectionTitle: {
    fontSize: 10,
    fontWeight: '700',
    color: COLORS.muted,
    letterSpacing: 1.4,
    marginBottom: 10,
  },
  emptyText: { fontSize: 12, color: COLORS.muted, textAlign: 'center', paddingVertical: 12 },

  // Lecturas
  lecturaRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  lecturaLeft: { flexDirection: 'row', alignItems: 'baseline', gap: 3 },
  lecturaValor: { fontSize: 18, fontWeight: '800', color: COLORS.teal },
  lecturaUnidad: { fontSize: 11, color: COLORS.muted },
  lecturaFecha: { fontSize: 11, color: COLORS.muted, textAlign: 'right' },
  lecturaObs: { fontSize: 10, color: COLORS.muted, marginTop: 2, textAlign: 'right' },

  // Alertas
  alertaRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
    gap: 8,
  },
  alertaTipo: { fontSize: 12, fontWeight: '700', color: COLORS.text, textTransform: 'capitalize' },
  alertaDetalle: { fontSize: 10, color: COLORS.muted, marginTop: 2 },
  alertaEstadoBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 8,
  },

  // Empty states
  emptyContainer: {
    paddingVertical: 32,
    alignItems: 'center',
    gap: 8,
  },
  emptyTitle: { fontSize: 15, fontWeight: '700', color: COLORS.muted },
  emptyDesc: {
    fontSize: 12,
    color: COLORS.muted,
    textAlign: 'center',
    paddingHorizontal: 16,
  },

  detailPanel: { flex: 1 },
});
