import { useState, useEffect } from 'react';
import { ParqueService } from '../services/parque.service';
import { Parque } from '../dto/parque.dto';

export function useParques() {
  const [parques, setParques] = useState<Parque[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchParques();
  }, []);

  const fetchParques = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await ParqueService.getParques();
      if (result.success !== false) { // Asumimos exito por defecto si no es explícitamente false
        setParques(result.data || result || []);
      } else {
        setError(result.error?.message || 'Error fetching parques');
      }
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const addParque = async (parque: Omit<Parque, 'id'>) => {
    try {
      const result = await ParqueService.createParque(parque);
      if (result.success !== false) {
        await fetchParques();
        return { success: true };
      }
      return { success: false, error: result.error?.message };
    } catch (e: any) {
      return { success: false, error: e.message };
    }
  };

  const removeParque = async (id: string) => {
    try {
      const result = await ParqueService.deleteParque(id);
      if (result.success !== false) {
        await fetchParques();
        return { success: true };
      }
      return { success: false, error: result.error?.message };
    } catch (e: any) {
      return { success: false, error: e.message };
    }
  };

  return {
    parques,
    loading,
    error,
    refresh: fetchParques,
    addParque,
    removeParque,
  };
}
