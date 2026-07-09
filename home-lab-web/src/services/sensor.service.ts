import { fetchWithAuth, API_SENSOR_URL } from './api';
import {
  Sensor,
  LecturaSensor,
  Alerta,
  CreateSensorRequest,
  UpdateSensorRequest,
  CreateLecturaRequest,
  CreateAlertaRequest,
  UpdateAlertaStatusRequest,
} from '../dto/sensor.dto';


export class SensorService {
  // ── Sensores ──────────────────────────────────────────────────────────────

  static async getSensoresByArbolId(arbolId: string): Promise<Sensor[]> {
    try {
      const response = await fetchWithAuth(
        `${API_SENSOR_URL}/sensores?arbol_id=${arbolId}`,
      );
      if (!response.ok) throw new Error('Error al obtener sensores');
      const data = await response.json();
      return Array.isArray(data) ? data : data.data ?? [];
    } catch (error) {
      console.error('Error fetching sensores:', error);
      throw error;
    }
  }

  static async getSensorById(id: string): Promise<Sensor> {
    const response = await fetchWithAuth(`${API_SENSOR_URL}/sensores/${id}`);
    if (!response.ok) throw new Error('Error al obtener sensor');
    const data = await response.json();
    return data.data ?? data;
  }

  static async createSensor(payload: CreateSensorRequest): Promise<Sensor> {
    const response = await fetchWithAuth(`${API_SENSOR_URL}/sensores`, {
      method: 'POST',
      body: JSON.stringify(payload),
    });
    if (!response.ok) throw new Error('Error al crear sensor');
    const data = await response.json();
    return data.data ?? data;
  }

  static async updateSensor(id: string, payload: UpdateSensorRequest): Promise<Sensor> {
    const response = await fetchWithAuth(`${API_SENSOR_URL}/sensores/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
    if (!response.ok) throw new Error('Error al actualizar sensor');
    const data = await response.json();
    return data.data ?? data;
  }

  static async deleteSensor(id: string): Promise<void> {
    const response = await fetchWithAuth(`${API_SENSOR_URL}/sensores/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Error al eliminar sensor');
  }

  // ── Lecturas ──────────────────────────────────────────────────────────────

  static async getLecturasBySensorId(sensorId: string): Promise<LecturaSensor[]> {
    try {
      const response = await fetchWithAuth(
        `${API_SENSOR_URL}/sensores/${sensorId}/lecturas`,
      );
      if (!response.ok) throw new Error('Error al obtener lecturas');
      const data = await response.json();
      return Array.isArray(data) ? data : data.data ?? [];
    } catch (error) {
      console.error('Error fetching lecturas:', error);
      throw error;
    }
  }

  static async createLectura(
    sensorId: string,
    payload: CreateLecturaRequest,
  ): Promise<LecturaSensor> {
    const response = await fetchWithAuth(
      `${API_SENSOR_URL}/sensores/${sensorId}/lecturas`,
      {
        method: 'POST',
        body: JSON.stringify(payload),
      },
    );
    if (!response.ok) throw new Error('Error al registrar lectura');
    const data = await response.json();
    return data.data ?? data;
  }

  // ── Alertas ───────────────────────────────────────────────────────────────

  static async getAlertas(filters?: {
    sensor_id?: string;
    estado?: string;
  }): Promise<Alerta[]> {
    const params = new URLSearchParams();
    if (filters?.sensor_id) params.append('sensor_id', filters.sensor_id);
    if (filters?.estado) params.append('estado', filters.estado);
    const qs = params.toString() ? `?${params.toString()}` : '';

    const response = await fetchWithAuth(`${API_SENSOR_URL}/alertas${qs}`);
    if (!response.ok) throw new Error('Error al obtener alertas');
    const data = await response.json();
    return Array.isArray(data) ? data : data.data ?? [];
  }

  static async createAlerta(
    sensorId: string,
    payload: CreateAlertaRequest,
  ): Promise<Alerta> {
    const response = await fetchWithAuth(
      `${API_SENSOR_URL}/sensores/${sensorId}/alertas`,
      {
        method: 'POST',
        body: JSON.stringify(payload),
      },
    );
    if (!response.ok) throw new Error('Error al registrar alerta');
    const data = await response.json();
    return data.data ?? data;
  }

  static async updateAlertaEstado(
    alertaId: string,
    payload: UpdateAlertaStatusRequest,
  ): Promise<Alerta> {
    const response = await fetchWithAuth(
      `${API_SENSOR_URL}/alertas/${alertaId}/estado`,
      {
        method: 'PUT',
        body: JSON.stringify(payload),
      },
    );
    if (!response.ok) throw new Error('Error al actualizar alerta');
    const data = await response.json();
    return data.data ?? data;
  }
}
