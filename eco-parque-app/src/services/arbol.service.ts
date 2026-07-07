import { API_FLORA_URL, fetchWithAuth } from "./api";
import { Arbol, CreateArbolRequest } from "../dto/arbol.dto";

export class ArbolService {
  /**
   * Obtiene la lista de árboles asociados a una zona.
   */
  static async getArbolesByZonaId(zonaId: string): Promise<Arbol[]> {
    try {
      const response = await fetchWithAuth(
        `${API_FLORA_URL}/arboles?zona_id=${zonaId}`,
      );
      if (!response.ok) {
        throw new Error("Error al obtener árboles");
      }
      const data = await response.json();
      console.log("árboles obtenidos:", data);

      return data;
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
      const response = await fetchWithAuth(`${API_FLORA_URL}/arboles`, {
        method: "POST",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("Error al crear el árbol");
      }

      return await response.json();
    } catch (error) {
      console.error("Error creating árbol:", error);
      throw error;
    }
  }
}
