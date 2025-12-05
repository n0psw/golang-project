import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Button } from '../common/Button';
import './Header.css';

export const Header: React.FC = () => {
  const { user, logout, isAdmin } = useAuth();
  const navigate = useNavigate();
  const [showDropdown, setShowDropdown] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
    setShowDropdown(false);
  };

  return (
    <header className="header">
      <div className="container">
        <div className="header-content">
          <Link to="/" className="header-logo">
            Movie Reviews
          </Link>
          <nav className="header-nav">
            <Link to="/" className="nav-link">
              Главная
            </Link>
            {user && (
              <Link to="/my-reviews" className="nav-link">
                Мои отзывы
              </Link>
            )}
            {isAdmin() && (
              <Link to="/admin" className="nav-link">
                Админ
              </Link>
            )}
          </nav>
          <div className="header-auth">
            {user ? (
              <div className="user-menu">
                <button
                  className="user-button"
                  onClick={() => setShowDropdown(!showDropdown)}
                >
                  {user.username}
                </button>
                {showDropdown && (
                  <div className="user-dropdown">
                    <div className="user-info">
                      <p className="user-email">{user.email}</p>
                      <p className="user-role">{user.role}</p>
                    </div>
                    <button className="dropdown-item" onClick={handleLogout}>
                      Выйти
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <Button onClick={() => navigate('/login')} variant="primary">
                Войти
              </Button>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

