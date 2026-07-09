export interface Arbol {
  id: string;
  zona_id: string;
  especie_id?: string;
  estado_id?: string;
  codigo?: string;
  latitud?: number;
  longitud?: number;
  edad_estimada_anios?: number;
  altura_m?: number;
  ancho_copa_m?: number;
  diametro_tronco_cm?: number;
  fecha_plantacion?: string;
  fecha_registro?: string;
  observaciones?: string;
  fecha_creacion: string;
  fecha_actualizacion?: string;
}

export interface CreateArbolRequest {
  zona_id: string;
  especie_id?: string;
  estado_id?: string;
  codigo?: string;
  latitud?: number;
  longitud?: number;
  edad_estimada_anios?: number;
  altura_m?: number;
  ancho_copa_m?: number;
  diametro_tronco_cm?: number;
  fecha_plantacion?: string;
  observaciones?: string;
}
