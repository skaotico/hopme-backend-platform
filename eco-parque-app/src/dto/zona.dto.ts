export interface Zona {
  id: string;
  ecoparque_id: string;
  nombre: string;
  descripcion?: string;
  area_m2?: number;
  fecha_creacion: string;
  fecha_actualizacion?: string;
}

export interface CreateZonaRequest {
  ecoparque_id: string;
  nombre: string;
  descripcion?: string;
  area_m2?: number;
}
