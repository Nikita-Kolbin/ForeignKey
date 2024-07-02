import React, { useState } from 'react';
import { API_BASE_URL } from '../ApiConfig';

const CustomerSignUpForm = ({ alias, onSuccess }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`${API_BASE_URL}/customer/sign-up`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ alias, email, password })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        onSuccess();
      } else {
        setError(data.error || 'Ошибка при регистрации');
      }
    } catch (error) {
      console.error('Ошибка при регистрации:', error);
      setError('Ошибка при регистрации. Попробуйте снова позже.');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h3>Регистрация</h3>
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
      <button type="submit">Зарегистрироваться</button>
    </form>
  );
};

export default CustomerSignUpForm;
