import { useState, useCallback } from 'react';
import { Zona, CreateZonaRequest } from '../dto/zona.dto';
import { ZonaService } from '../services/zona.service';

export function useZonas(parqueId: string | null) {
  const [zonas, setZonas] = useState<Zona[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchZonas = useCallback(async () => {
    if (!parqueId) return;
    
    setLoading(true);
    setError(null);
    try {
      const data = await ZonaService.getZonasByParqueId(parqueId);
      setZonas(data || []);
    } catch (err) {
      setError('No se pudieron cargar las zonas. Inténtalo más tarde.');
    } finally {
      setLoading(false);
    }
  }, [parqueId]);

  const addZona = async (data: Omit<CreateZonaRequest, 'ecoparque_id'>) => {
    if (!parqueId) return { success: false, error: 'No hay parque seleccionado' };
    
    try {
      const newZona = await ZonaService.createZona({ ...data, ecoparque_id: parqueId });
      setZonas(prev => [...prev, newZona]);
      return { success: true };
    } catch (err) {
      return { success: false, error: 'No se pudo crear la zona.' };
    }
  };

  return {
    zonas,
    loading,
    error,
    refresh: fetchZonas,
    addZona,
  };
}
