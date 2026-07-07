import { API_ZONA_URL, fetchWithAuth } from './api';
import { Zona, CreateZonaRequest } from '../dto/zona.dto';

export class ZonaService {
  /**
   * Obtiene la lista de zonas asociadas a un ecoparque.
   */
  static async getZonasByParqueId(ecoparqueId: string): Promise<Zona[]> {
    try {
      const response = await fetchWithAuth(`${API_ZONA_URL}/zonas?ecoparque_id=${ecoparqueId}`);
      if (!response.ok) {
        throw new Error('Error al obtener zonas');
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching zonas:', error);
      throw error;
    }
  }

  /**
   * Crea una nueva zona para un ecoparque.
   */
  static async createZona(data: CreateZonaRequest): Promise<Zona> {
    try {
      const response = await fetchWithAuth(`${API_ZONA_URL}/zonas`, {
        method: 'POST',
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error('Error al crear la zona');
      }

      return await response.json();
    } catch (error) {
      console.error('Error creating zona:', error);
      throw error;
    }
  }
}
