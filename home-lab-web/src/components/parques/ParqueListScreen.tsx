import React, { useState, useEffect } from 'react';
import { Plus, MapPin, Loader2 } from 'lucide-react';
import { ParqueService } from '../../services/parque.service';
import { Parque } from '../../dto/parque.dto';

export const ParqueListScreen = () => {
  const [parques, setParques] = useState<Parque[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadParques();
  }, []);

  const loadParques = async () => {
    try {
      setLoading(true);
      const data = await ParqueService.getParques();
      setParques(data);
    } catch (err: any) {
      setError(err.message || 'Error al cargar los parques');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="screen-container animate-fade-in">
      <div className="screen-header">
        <h1>Parques Eco-Lógicos</h1>
        <button className="btn btn-primary">
          <Plus size={20} />
          Nuevo Parque
        </button>
      </div>

      {error && (
        <div className="card glass-panel" style={{ borderColor: 'var(--danger)', marginBottom: '1rem' }}>
          <p className="error-text" style={{ margin: 0 }}>{error}</p>
        </div>
      )}

      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '3rem' }}>
          <Loader2 className="animate-spin" size={32} color="var(--accent-primary)" />
        </div>
      ) : parques.length === 0 ? (
        <div className="card glass-panel" style={{ textAlign: 'center', padding: '4rem 2rem' }}>
          <MapPin size={48} color="var(--text-secondary)" style={{ margin: '0 auto 1rem', opacity: 0.5 }} />
          <h3 style={{ color: 'var(--text-secondary)' }}>No hay parques registrados</h3>
          <p>Comienza añadiendo un nuevo parque al sistema.</p>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
          {parques.map((parque) => (
            <div key={parque.id} className="card clickable glass-panel">
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1rem' }}>
                <div style={{ backgroundColor: 'var(--accent-light)', padding: '0.75rem', borderRadius: 'var(--radius-md)', color: 'var(--accent-primary)' }}>
                  <MapPin size={24} />
                </div>
                <div>
                  <h3 style={{ margin: 0, fontSize: '1.25rem' }}>{parque.nombre}</h3>
                  <p style={{ margin: 0, fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
                    {parque.direccion || 'Sin dirección'}
                  </p>
                </div>
              </div>
              
              <div style={{ display: 'flex', gap: '0.5rem', marginTop: '1.5rem' }}>
                <button className="btn btn-secondary w-full" style={{ fontSize: '0.875rem', padding: '0.5rem' }}>
                  Ver Zonas
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
