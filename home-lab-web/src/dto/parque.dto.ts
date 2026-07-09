export interface Parque {
  id: string;
  nombre: string;
  descripcion: string;
  direccion: string;
  [key: string]: any; // por si hay más campos en el struct
}
