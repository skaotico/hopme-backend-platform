import { useState, useCallback } from "react";
import { CatalogItem, CatalogType } from "../dto/catalogo.dto";
import { CatalogoService } from "../services/catalogo.service";

export function useCatalogo(type: CatalogType) {
  const [items, setItems] = useState<CatalogItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchItems = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await CatalogoService.list(type);
      setItems(data || []);
    } catch (err) {
      setError(`No se pudieron cargar los elementos de ${type}. Inténtalo más tarde.`);
    } finally {
      setLoading(false);
    }
  }, [type]);

  const addItem = async (data: any) => {
    try {
      const newItem = await CatalogoService.create(type, data);
      setItems((prev) => [...prev, newItem]);
      return { success: true };
    } catch (err) {
      return { success: false, error: "No se pudo crear el registro." };
    }
  };

  const updateItem = async (id: string, data: any) => {
    try {
      const updated = await CatalogoService.update(type, id, data);
      setItems((prev) => prev.map((item) => (item.id === id ? updated : item)));
      return { success: true };
    } catch (err) {
      return { success: false, error: "No se pudo actualizar el registro." };
    }
  };

  const removeItem = async (id: string) => {
    try {
      await CatalogoService.delete(type, id);
      setItems((prev) => prev.filter((item) => item.id !== id));
      return { success: true };
    } catch (err) {
      return { success: false, error: "No se pudo eliminar el registro." };
    }
  };

  return {
    items,
    loading,
    error,
    refresh: fetchItems,
    addItem,
    updateItem,
    removeItem,
  };
}
