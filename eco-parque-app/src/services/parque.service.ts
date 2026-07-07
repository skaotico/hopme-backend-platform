import { fetchWithAuth, API_PARQUE_URL } from './api';
import { Parque } from '../dto/parque.dto';

export const ParqueService = {
  /**
   * Obtiene la lista de parques.
   */
  async getParques() {
    const response = await fetchWithAuth(`${API_PARQUE_URL}/parques`, {
      method: 'GET',
    });
    return response.json();
  },

  /**
   * Obtiene un parque por ID.
   */
  async getParqueById(id: string) {
    const response = await fetchWithAuth(`${API_PARQUE_URL}/parques/${id}`, {
      method: 'GET',
    });
    return response.json();
  },

  /**
   * Crea un nuevo parque.
   */
  async createParque(parque: Omit<Parque, 'id'>) {
    const response = await fetchWithAuth(`${API_PARQUE_URL}/parques`, {
      method: 'POST',
      body: JSON.stringify(parque),
    });
    return response.json();
  },

  /**
   * Actualiza un parque.
   */
  async updateParque(id: string, parque: Partial<Parque>) {
    const response = await fetchWithAuth(`${API_PARQUE_URL}/parques/${id}`, {
      method: 'PUT',
      body: JSON.stringify(parque),
    });
    return response.json();
  },

  /**
   * Elimina un parque.
   */
  async deleteParque(id: string) {
    const response = await fetchWithAuth(`${API_PARQUE_URL}/parques/${id}`, {
      method: 'DELETE',
    });
    return response.json();
  },
};
