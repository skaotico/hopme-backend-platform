const AsyncStorage = {
  getItem: async (key: string) => localStorage.getItem(key),
  setItem: async (key: string, value: string) => localStorage.setItem(key, value),
  removeItem: async (key: string) => localStorage.removeItem(key),
};
export const API_AUTH_URL = "http://192.168.1.27:9090/api/v1";
export const API_PARQUE_URL = "http://192.168.1.27:9091/api/v1";
export const API_ZONA_URL = "http://192.168.1.27:9092/api/v1";
export const API_FLORA_URL = "http://192.168.1.27:9093/api/v1";
export const API_CATALOGO_URL = "http://192.168.1.27:9094/api/v1";
export const API_ARBOL_URL = "http://192.168.1.27:9095/api/v1";
export const API_SENSOR_URL = "http://192.168.1.27:9096/api/v1";



export const TOKEN_KEY = "@eco_parque_access_token";

/**
 * Realiza una petición `fetch` inyectando automáticamente el token JWT
 * si existe en AsyncStorage.
 */
export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  console.log('MOCK FETCH:', options.method || 'GET', url);
  
  const mockResponse = (data: any, status = 200) => {
    return {
      ok: status >= 200 && status < 300,
      status,
      json: async () => data
    } as Response;
  };

  if (url.includes('/auth/login') || url.includes('/auth/register')) {
    return mockResponse({ token: 'mock-token', user: { id: '1', email: 'mock@test.com' } });
  }

  // Handle catalogo endpoints statically
  if (url.includes('catalogo/especies')) {
    return mockResponse([{ id: '1', nombre_comun: 'Roble', nombre_cientifico: 'Quercus', familia: 'Fagaceae', requerimientos_hidricos: 'Medio', requerimientos_luz: 'Alta' }]);
  }
  if (url.includes('catalogo/estados')) {
    return mockResponse([{ id: '1', codigo: 'SANO', nombre: 'Sano', descripcion: 'Buen estado', severidad: 0 }]);
  }
  
  const pathParts = url.split('/').filter(Boolean);
  const entityName = pathParts.length > 3 ? pathParts[3] : 'data'; // e.g. "parques", "zonas"
  const storageKey = `@mock_${entityName}`;
  
  const getStored = async () => {
    const data = await AsyncStorage.getItem(storageKey);
    return data ? JSON.parse(data) : [];
  };
  
  const method = options.method || 'GET';
  const isListReq = url.includes('?') || !url.split('/').pop()?.match(/^[0-9a-fA-F-]+$/); // rough guess if it's a list or detail

  if (method === 'GET') {
    const data = await getStored();
    
    // Simulate query params filtering for child elements
    if (url.includes('parque_id=')) {
      const pId = url.split('parque_id=')[1].split('&')[0];
      return mockResponse(data.filter((d: any) => d.parque_id === pId || d.ecoparque_id === pId));
    }
    if (url.includes('zona_id=')) {
      const zId = url.split('zona_id=')[1].split('&')[0];
      return mockResponse(data.filter((d: any) => d.zona_id === zId));
    }
    if (url.includes('arbol_id=')) {
      const aId = url.split('arbol_id=')[1].split('&')[0];
      return mockResponse(data.filter((d: any) => d.arbol_id === aId));
    }
    if (url.includes('sensor_id=')) {
      const sId = url.split('sensor_id=')[1].split('&')[0];
      return mockResponse(data.filter((d: any) => d.sensor_id === sId));
    }

    return mockResponse(data);
  }
  
  if (method === 'POST') {
    const body = options.body ? JSON.parse(options.body as string) : {};
    const data = await getStored();
    const newItem = { id: Math.random().toString(36).substr(2, 9), fecha_creacion: new Date().toISOString(), ...body };
    data.push(newItem);
    await AsyncStorage.setItem(storageKey, JSON.stringify(data));
    return mockResponse(newItem, 201);
  }
  
  if (method === 'PUT' || method === 'PATCH') {
    const id = pathParts[pathParts.length - 1];
    const body = options.body ? JSON.parse(options.body as string) : {};
    const data = await getStored();
    const idx = data.findIndex((d: any) => d.id === id);
    if (idx !== -1) {
      data[idx] = { ...data[idx], ...body };
      await AsyncStorage.setItem(storageKey, JSON.stringify(data));
      return mockResponse(data[idx]);
    }
    return mockResponse({ error: 'Not found' }, 404);
  }
  
  if (method === 'DELETE') {
    const id = pathParts[pathParts.length - 1];
    let data = await getStored();
    data = data.filter((d: any) => d.id !== id);
    await AsyncStorage.setItem(storageKey, JSON.stringify(data));
    return mockResponse({}, 204);
  }
  
  return mockResponse([]);
}
