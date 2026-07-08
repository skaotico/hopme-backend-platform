export interface Sensor {
  id: string;
  arbol_id?: string;
  estanque_id?: string;
  tipo_sensor_id?: string;
  codigo?: string;
  modelo?: string;
  fabricante?: string;
  fecha_instalacion?: string;
  activo: boolean;
}

export interface CreateSensorRequest {
  arbol_id?: string;
  estanque_id?: string;
  tipo_sensor_id?: string;
  codigo?: string;
  modelo?: string;
  fabricante?: string;
  fecha_instalacion?: string;
}

export interface UpdateSensorRequest {
  arbol_id?: string;
  estanque_id?: string;
  tipo_sensor_id?: string;
  codigo?: string;
  modelo?: string;
  fabricante?: string;
  fecha_instalacion?: string;
  activo?: boolean;
}

export interface LecturaSensor {
  id: string;
  sensor_id: string;
  valor: number;
  fecha_lectura: string;
  bateria_porcentaje?: number;
  observacion?: string;
}

export interface CreateLecturaRequest {
  valor: number;
  fecha_lectura?: string;
  bateria_porcentaje?: number;
  observacion?: string;
}

export interface Alerta {
  id: string;
  sensor_id: string;
  tipo: string;
  estado: string;
  umbral: number;
  valor_detectado: number;
  fecha_alerta: string;
  observacion?: string;
}

export interface CreateAlertaRequest {
  tipo: string;
  umbral: number;
  valor_detectado: number;
  fecha_alerta?: string;
  observacion?: string;
}

export interface UpdateAlertaStatusRequest {
  estado: string;
  observacion?: string;
}
