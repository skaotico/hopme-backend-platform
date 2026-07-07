import { useState, useCallback } from 'react';
import { Arbol, CreateArbolRequest } from '../dto/arbol.dto';
import { ArbolService } from '../services/arbol.service';

export function useArboles(zonaId: string | null) {
  const [arboles, setArboles] = useState<Arbol[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchArboles = useCallback(async () => {
    if (!zonaId) return;
    
    setLoading(true);
    setError(null);
    try {
      const data = await ArbolService.getArbolesByZonaId(zonaId);
      setArboles(data || []);
    } catch (err) {
      setError('No se pudieron cargar los árboles. Inténtalo más tarde.');
    } finally {
      setLoading(false);
    }
  }, [zonaId]);

  const addArbol = async (data: Omit<CreateArbolRequest, 'zona_id'>) => {
    if (!zonaId) return { success: false, error: 'No hay zona seleccionada' };
    
    try {
      const newArbol = await ArbolService.createArbol({ ...data, zona_id: zonaId });
      setArboles(prev => [...prev, newArbol]);
      return { success: true };
    } catch (err) {
      return { success: false, error: 'No se pudo registrar el árbol.' };
    }
  };

  const updateArbol = async (id: string, data: any) => {
    try {
      const updated = await ArbolService.updateArbol(id, data);
      setArboles(prev => prev.map(a => a.id === id ? updated : a));
      return { success: true };
    } catch (err) {
      return { success: false, error: 'No se pudo actualizar el árbol.' };
    }
  };

  const removeArbol = async (id: string) => {
    try {
      await ArbolService.deleteArbol(id);
      setArboles(prev => prev.filter(a => a.id !== id));
      return { success: true };
    } catch (err) {
      return { success: false, error: 'No se pudo eliminar el árbol.' };
    }
  };

  return {
    arboles,
    loading,
    error,
    refresh: fetchArboles,
    addArbol,
    updateArbol,
    removeArbol,
  };
}
