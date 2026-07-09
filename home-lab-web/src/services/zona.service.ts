import { API_ZONA_URL, fetchWithAuth } from "./api";
import { Zona, CreateZonaRequest, UpdateZonaRequest } from "../dto/zona.dto";

export class ZonaService {
  /**
   * Obtiene la lista de zonas asociadas a un ecoparque.
   */
  static async getZonasByParqueId(ecoparqueId: string): Promise<Zona[]> {
    try {
      const response = await fetchWithAuth(
        `${API_ZONA_URL}/zonas?ecoparque_id=${ecoparqueId}`,
      );
      if (!response.ok) {
        throw new Error("Error al obtener zonas");
      }
      const envelope = await response.json();
      console.log("Envelope:", envelope.data);
      return envelope.data ?? envelope;
    } catch (error) {
      console.error("Error fetching zonas:", error);
      throw error;
    }
  }

  /**
   * Crea una nueva zona para un ecoparque.
   */
  static async createZona(data: CreateZonaRequest): Promise<Zona> {
    try {
      const response = await fetchWithAuth(`${API_ZONA_URL}/zonas`, {
        method: "POST",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("Error al crear la zona");
      }

      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error("Error creating zona:", error);
      throw error;
    }
  }

  /**
   * Actualiza una zona existente por su ID.
   */
  static async updateZona(id: string, data: UpdateZonaRequest): Promise<Zona> {
    try {
      const response = await fetchWithAuth(`${API_ZONA_URL}/zonas/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("Error al actualizar la zona");
      }

      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error("Error updating zona:", error);
      throw error;
    }
  }

  /**
   * Elimina una zona por su ID.
   */
  static async deleteZona(id: string): Promise<void> {
    try {
      const response = await fetchWithAuth(`${API_ZONA_URL}/zonas/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error("Error al eliminar la zona");
      }
    } catch (error) {
      console.error("Error deleting zona:", error);
      throw error;
    }
  }
}
