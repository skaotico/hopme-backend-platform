import { API_ARBOL_URL, fetchWithAuth } from "./api";
import { Arbol, CreateArbolRequest } from "../dto/arbol.dto";

export class ArbolService {
  /**
   * Obtiene la lista de árboles asociados a una zona.
   */
  static async getArbolesByZonaId(zonaId: string): Promise<Arbol[]> {
    try {
      const response = await fetchWithAuth(
        `${API_ARBOL_URL}/arboles?zona_id=${zonaId}`,
      );
      if (!response.ok) {
        throw new Error("Error al obtener árboles");
      }
      const envelope = await response.json();
      console.log("árboles obtenidos:", envelope);

      return envelope.data ?? envelope;
    } catch (error) {
      console.error("Error fetching árboles:", error);
      throw error;
    }
  }

  /**
   * Crea un nuevo árbol para una zona.
   */
  static async createArbol(data: CreateArbolRequest): Promise<Arbol> {
    try {
      const response = await fetchWithAuth(`${API_ARBOL_URL}/arboles`, {
        method: "POST",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("Error al crear el árbol");
      }

      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error("Error creating árbol:", error);
      throw error;
    }
  }

  /**
   * Actualiza los datos de un árbol existente.
   */
  static async updateArbol(id: string, data: any): Promise<Arbol> {
    try {
      const response = await fetchWithAuth(`${API_ARBOL_URL}/arboles/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("Error al actualizar el árbol");
      }

      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error(`Error updating árbol ${id}:`, error);
      throw error;
    }
  }

  /**
   * Elimina un árbol.
   */
  static async deleteArbol(id: string): Promise<void> {
    try {
      const response = await fetchWithAuth(`${API_ARBOL_URL}/arboles/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error("Error al eliminar el árbol");
      }
    } catch (error) {
      console.error(`Error deleting árbol ${id}:`, error);
      throw error;
    }
  }
}
