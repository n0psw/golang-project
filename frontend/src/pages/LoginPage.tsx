import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { LoginForm } from '../components/Auth/LoginForm';
import './AuthPage.css';

export const LoginPage: React.FC = () => {
  const navigate = useNavigate();

  const handleSuccess = () => {
    navigate('/');
  };

  return (
    <div className="auth-page">
      <div className="auth-container">
        <h1 className="auth-title">Вход</h1>
        <LoginForm onSuccess={handleSuccess} />
        <p className="auth-link">
          Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
        </p>
      </div>
    </div>
  );
};

