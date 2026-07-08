import { useState, useCallback } from 'react';
import { Sensor, LecturaSensor, Alerta } from '../dto/sensor.dto';
import { SensorService } from '../services/sensor.service';

export function useSensores(arbolId: string | null) {
  const [sensores, setSensores] = useState<Sensor[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchSensores = useCallback(async () => {
    if (!arbolId) return;
    setLoading(true);
    setError(null);
    try {
      const data = await SensorService.getSensoresByArbolId(arbolId);
      setSensores(data || []);
    } catch {
      setError('No se pudieron cargar los sensores.');
    } finally {
      setLoading(false);
    }
  }, [arbolId]);

  return { sensores, loading, error, refresh: fetchSensores };
}

export function useLecturas(sensorId: string | null) {
  const [lecturas, setLecturas] = useState<LecturaSensor[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchLecturas = useCallback(async () => {
    if (!sensorId) return;
    setLoading(true);
    setError(null);
    try {
      const data = await SensorService.getLecturasBySensorId(sensorId);
      // Sort ascending by date
      const sorted = [...data].sort(
        (a, b) =>
          new Date(a.fecha_lectura).getTime() - new Date(b.fecha_lectura).getTime(),
      );
      setLecturas(sorted);
    } catch {
      setError('No se pudieron cargar las lecturas.');
    } finally {
      setLoading(false);
    }
  }, [sensorId]);

  return { lecturas, loading, error, refresh: fetchLecturas };
}

export function useAlertas(sensorId: string | null) {
  const [alertas, setAlertas] = useState<Alerta[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAlertas = useCallback(async () => {
    if (!sensorId) return;
    setLoading(true);
    setError(null);
    try {
      const data = await SensorService.getAlertas({ sensor_id: sensorId });
      setAlertas(data || []);
    } catch {
      setError('No se pudieron cargar las alertas.');
    } finally {
      setLoading(false);
    }
  }, [sensorId]);

  return { alertas, loading, error, refresh: fetchAlertas };
}
