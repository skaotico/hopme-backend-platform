import { API_CATALOGO_URL, fetchWithAuth } from "./api";
import { CatalogItem, CatalogType } from "../dto/catalogo.dto";

export class CatalogoService {
  /**
   * Obtiene la lista de elementos de un catálogo específico.
   */
  static async list(type: CatalogType): Promise<CatalogItem[]> {
    try {
      const response = await fetchWithAuth(`${API_CATALOGO_URL}/catalogo/${type}`);
      if (!response.ok) {
        throw new Error(`Error al obtener catálogo ${type}`);
      }
      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error(`Error fetching catalog ${type}:`, error);
      throw error;
    }
  }

  /**
   * Obtiene un elemento por su ID de un catálogo específico.
   */
  static async getByID(type: CatalogType, id: string): Promise<CatalogItem> {
    try {
      const response = await fetchWithAuth(`${API_CATALOGO_URL}/catalogo/${type}/${id}`);
      if (!response.ok) {
        throw new Error(`Error al obtener elemento ${id} del catálogo ${type}`);
      }
      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error(`Error fetching catalog item ${type}/${id}:`, error);
      throw error;
    }
  }

  /**
   * Crea un nuevo elemento en un catálogo específico.
   */
  static async create(type: CatalogType, data: any): Promise<CatalogItem> {
    try {
      const response = await fetchWithAuth(`${API_CATALOGO_URL}/catalogo/${type}`, {
        method: "POST",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(`Error al crear elemento en catálogo ${type}`);
      }

      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error(`Error creating catalog item in ${type}:`, error);
      throw error;
    }
  }

  /**
   * Actualiza un elemento existente de un catálogo específico.
   */
  static async update(type: CatalogType, id: string, data: any): Promise<CatalogItem> {
    try {
      const response = await fetchWithAuth(`${API_CATALOGO_URL}/catalogo/${type}/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(`Error al actualizar elemento ${id} en catálogo ${type}`);
      }

      const envelope = await response.json();
      return envelope.data ?? envelope;
    } catch (error) {
      console.error(`Error updating catalog item ${type}/${id}:`, error);
      throw error;
    }
  }

  /**
   * Elimina un elemento de un catálogo específico.
   */
  static async delete(type: CatalogType, id: string): Promise<void> {
    try {
      const response = await fetchWithAuth(`${API_CATALOGO_URL}/catalogo/${type}/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(`Error al eliminar elemento ${id} del catálogo ${type}`);
      }
    } catch (error) {
      console.error(`Error deleting catalog item ${type}/${id}:`, error);
      throw error;
    }
  }
}
