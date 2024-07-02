import React, { useState } from 'react';
import { API_BASE_URL } from '../ApiConfig';

const CustomerSignInForm = ({ alias, onSuccess }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`${API_BASE_URL}/customer/sign-in`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ alias, email, password })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        localStorage.setItem('customerToken', data.token);
        onSuccess();
      } else {
        setError(data.error || 'Ошибка при входе');
      }
    } catch (error) {
      console.error('Ошибка при входе:', error);
      setError('Ошибка при входе. Попробуйте снова позже.');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h3>Вход</h3>
      <input 
        type="email" 
        placeholder="Email" 
        value={email} 
        onChange={(e) => setEmail(e.target.value)} 
      />
      <input 
        type="password" 
        placeholder="Пароль" 
        value={password} 
        onChange={(e) => setPassword(e.target.value)} 
      />
      {error && <div className="error-message">{error}</div>}
      <button type="submit">Войти</button>
    </form>
  );
};

export default CustomerSignInForm;
