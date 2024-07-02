import React, { useState } from 'react';
import { Navigate } from 'react-router-dom';
import '../styles/RegisterPage.css'; // Подключаем файл со стилями
import { API_BASE_URL } from '../components/ApiConfig';

const RegisterPage = () => {
  const [registered, setRegistered] = useState(false);
  const [email, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleRegister = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/admin/sign-up`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setRegistered(true);
        localStorage.setItem('token', data.token);
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка регистрации:', error);
      setError('Ошибка регистрации. Попробуйте снова позже.');
    }
  };

  if (registered) {
    return <Navigate to="/profile" />;
  }

  return (
    <div className="register-container">
      <div className="register-content">
        <div className="register-form">
          <div className="auth-options">
            <a href="/login" className="auth-link">Вход</a>
            <a href="#" className="auth-link active">Регистрация</a>
          </div>
          <button className="vk-button">Продолжить с Vk</button>
          <p className="separator">Или</p>
          <input
            type="text"
            placeholder="Логин"
            className="register-input"
            value={email}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="password"
            placeholder="Пароль"
            className="register-input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button onClick={handleRegister} className="register-button">Войти</button>
          {error && <div className="error-message">{error}</div>}
        </div>
        <div className="image-container">
          <img src="/public/Баннер.png" alt="Registration" className="side-image" />
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
