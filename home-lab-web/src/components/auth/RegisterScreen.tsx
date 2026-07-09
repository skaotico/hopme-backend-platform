import React, { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { Loader2 } from 'lucide-react';
import './Auth.css';

interface RegisterScreenProps {
  onGoToLogin: () => void;
}

export const RegisterScreen: React.FC<RegisterScreenProps> = ({ onGoToLogin }) => {
  const { register, login } = useAuth();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsSubmitting(true);
    
    const result = await register(username, email, password);
    if (result.success) {
      // Auto login after register
      await login(email, password);
    } else {
      setError(result.error || 'Error al registrar usuario');
      setIsSubmitting(false);
    }
  };

  return (
    <div className="auth-container animate-fade-in">
      <div className="card auth-card glass-panel">
        <h1 className="auth-title">Registro</h1>
        <p className="auth-subtitle">Crea una nueva cuenta</p>
        
        <form className="auth-form" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="username">Nombre de Usuario</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Ej. juanperez"
              required
              disabled={isSubmitting}
            />
          </div>
          <div>
            <label htmlFor="email">Correo Electrónico</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="tu@email.com"
              required
              disabled={isSubmitting}
            />
          </div>
          <div>
            <label htmlFor="password">Contraseña</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              disabled={isSubmitting}
            />
          </div>
          
          {error && <p className="error-text">{error}</p>}
          
          <button 
            type="submit" 
            className="btn btn-primary w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? <Loader2 className="animate-spin" size={20} /> : 'Registrarse'}
          </button>
        </form>
        
        <div className="auth-footer">
          ¿Ya tienes cuenta?
          <span className="auth-link" onClick={onGoToLogin}>Inicia sesión aquí</span>
        </div>
      </div>
    </div>
  );
};
