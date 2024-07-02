import React, { useState } from 'react';
import { Link, Navigate } from 'react-router-dom';
import '../styles/LoginPage.css'; // Подключаем файл со стилями
import { API_BASE_URL } from '../components/ApiConfig';

const LoginPage = () => {
  const [loggedIn, setLoggedIn] = useState(false);
  const [email, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleLogin = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/admin/sign-in`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setLoggedIn(true);
        localStorage.setItem('token', data.token);
        localStorage.setItem('login', email); // Сохраняем имя пользователя в localStorage
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка входа:', error);
      setError('Ошибка входа. Попробуйте снова позже.');
    }
  };

  if (loggedIn) {
    return <Navigate to="/profile" />;
  }

  return (
    <div className="login-container">
      <div className="login-content">
        <div className="login-form">
          <div className="auth-options">
            <a href="#" className="auth-link active">Вход</a>
            <a href="/register" className="auth-link">Регистрация</a>
          </div>
          <button className="vk-button">Продолжить с Vk</button>
          <p className="separator">Или</p>
          <input
            type="text"
            placeholder="Логин"
            className="login-input"
            value={email}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="password"
            placeholder="Пароль"
            className="login-input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button onClick={handleLogin} className="login-form-button">Войти</button>
          {error && <div className="error-message">{error}</div>}
        </div>
        <div className="image-container">
          <img src="/public/Баннер.png" alt="Login" className="side-image" />
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
