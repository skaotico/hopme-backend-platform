export interface EspecieArbol {
  id: string;
  nombre_comun?: string;
  nombre_cientifico: string;
  familia?: string;
  descripcion?: string;
  activo: boolean;
}

export interface EstadoArbol {
  id: string;
  codigo?: string;
  nombre?: string;
  descripcion?: string;
}

export interface EstadoEstanque {
  id: string;
  codigo?: string;
  nombre?: string;
}

export interface EstadoAgua {
  id: string;
  codigo?: string;
  nombre?: string;
}

export interface TipoSensor {
  id: string;
  codigo?: string;
  nombre?: string;
  unidad_medida?: string;
  descripcion?: string;
}

export type CatalogItem = EspecieArbol | EstadoArbol | EstadoEstanque | EstadoAgua | TipoSensor;

export type CatalogType = 'especies' | 'estados-arbol' | 'estados-estanque' | 'estados-agua' | 'tipos-sensor';

export interface CatalogConfig {
  type: CatalogType;
  displayName: string;
  icon: string;
  fields: {
    name: string;
    label: string;
    type: 'text' | 'boolean' | 'textarea';
    required?: boolean;
  }[];
}

export const CATALOGS_CONFIG: Record<CatalogType, CatalogConfig> = {
  'especies': {
    type: 'especies',
    displayName: 'Especies de Árbol',
    icon: 'nature',
    fields: [
      { name: 'nombre_cientifico', label: 'Nombre Científico', type: 'text', required: true },
      { name: 'nombre_comun', label: 'Nombre Común', type: 'text' },
      { name: 'familia', label: 'Familia', type: 'text' },
      { name: 'descripcion', label: 'Descripción', type: 'textarea' },
      { name: 'activo', label: 'Activo', type: 'boolean' }
    ]
  },
  'estados-arbol': {
    type: 'estados-arbol',
    displayName: 'Estados de Árbol',
    icon: 'health-and-safety',
    fields: [
      { name: 'codigo', label: 'Código', type: 'text', required: true },
      { name: 'nombre', label: 'Nombre', type: 'text', required: true },
      { name: 'descripcion', label: 'Descripción', type: 'textarea' }
    ]
  },
  'estados-estanque': {
    type: 'estados-estanque',
    displayName: 'Estados de Estanque',
    icon: 'opacity',
    fields: [
      { name: 'codigo', label: 'Código', type: 'text', required: true },
      { name: 'nombre', label: 'Nombre', type: 'text', required: true }
    ]
  },
  'estados-agua': {
    type: 'estados-agua',
    displayName: 'Estados de Agua',
    icon: 'water',
    fields: [
      { name: 'codigo', label: 'Código', type: 'text', required: true },
      { name: 'nombre', label: 'Nombre', type: 'text', required: true }
    ]
  },
  'tipos-sensor': {
    type: 'tipos-sensor',
    displayName: 'Tipos de Sensor',
    icon: 'settings-input-component',
    fields: [
      { name: 'codigo', label: 'Código', type: 'text', required: true },
      { name: 'nombre', label: 'Nombre', type: 'text', required: true },
      { name: 'unidad_medida', label: 'Unidad de Medida', type: 'text' },
      { name: 'descripcion', label: 'Descripción', type: 'textarea' }
    ]
  }
};
