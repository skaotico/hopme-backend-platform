import React, { useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Menu, X, LogOut, TreePine, Map, Bell, Database } from 'lucide-react';
import './Layout.css';

export const Layout = () => {
  const { user, logout } = useAuth();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  const navItems = [
    { name: 'Parques', path: '/', icon: <TreePine size={20} /> },
    { name: 'Notificaciones', path: '/notifications', icon: <Bell size={20} /> },
    { name: 'Catálogos', path: '/catalogos/especies', icon: <Database size={20} /> },
  ];

  return (
    <div className="layout-container">
      {/* Mobile Header */}
      <header className="mobile-header">
        <button onClick={toggleSidebar} className="menu-btn">
          <Menu size={24} />
        </button>
        <span className="mobile-title">EcoParque Web</span>
      </header>

      {/* Sidebar Overlay */}
      {isSidebarOpen && (
        <div className="sidebar-overlay" onClick={toggleSidebar} />
      )}

      {/* Sidebar */}
      <aside className={`sidebar ${isSidebarOpen ? 'open' : ''}`}>
        <div className="sidebar-header">
          <h2 className="sidebar-title">EcoParque</h2>
          <button className="close-btn d-md-none" onClick={toggleSidebar}>
            <X size={24} />
          </button>
        </div>

        <div className="user-info">
          <p className="user-email">{user?.email}</p>
        </div>

        <nav className="sidebar-nav">
          {navItems.map((item) => (
            <NavLink
              key={item.name}
              to={item.path}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
              onClick={() => setIsSidebarOpen(false)}
            >
              {item.icon}
              <span>{item.name}</span>
            </NavLink>
          ))}
        </nav>

        <div className="sidebar-footer">
          <button className="btn btn-secondary w-full logout-btn" onClick={logout}>
            <LogOut size={20} />
            <span>Cerrar Sesión</span>
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
};
