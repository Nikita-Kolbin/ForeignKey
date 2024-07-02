import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { API_BASE_URL } from '../components/ApiConfig';

const ConstructorPage = () => {
  const [aliases, setAliases] = useState([]);
  const [newAlias, setNewAlias] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const fetchAliases = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/website/aliases`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setAliases(data.aliases);
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при получении списка сайтов:', error);
      setError('Ошибка при получении списка сайтов. Попробуйте снова позже.');
    }
  };

  const createWebsite = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/website/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ alias: newAlias })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setNewAlias('');
        fetchAliases();
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при создании нового сайта:', error);
      setError('Ошибка при создании нового сайта. Попробуйте снова позже.');
    }
  };

  const deleteWebsite = async (alias) => {
    try {
      const response = await fetch(`${API_BASE_URL}/website/delete/${alias}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        fetchAliases();
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при удалении сайта:', error);
      setError('Ошибка при удалении сайта. Попробуйте снова позже.');
    }
  };

  const handleSelectWebsite = (alias) => {
    navigate(`/constructor/${alias}`);
  };

  const handleOpenWebsite = (alias) => {
    window.open(`${window.location.origin}/${alias}`, '_blank');
  };

  useEffect(() => {
    fetchAliases();
  }, []);

  return (
    <div className="constructor-page">
      <h2>Страница конструктора</h2>
      <div className="constructor-content">
        {aliases.length > 0 ? (
          <div>
            <h3>Список созданных сайтов:</h3>
            <ul>
              {aliases.map((alias, index) => (
                <li key={index}>
                  <button onClick={() => handleSelectWebsite(alias)}>Редактировать</button>
                  <button onClick={() => handleOpenWebsite(alias)}>Открыть</button>
                  <button onClick={() => deleteWebsite(alias)}>Удалить</button>
                </li>
              ))}
            </ul>
          </div>
        ) : (
          <p>У вас пока нет созданных сайтов</p>
        )}
        <input
          type="text"
          value={newAlias}
          onChange={(e) => setNewAlias(e.target.value)}
          placeholder="Введите имя нового сайта"
        />
        <button onClick={createWebsite}>Создать новый сайт</button>
        {error && <div className="error-message">{error}</div>}
      </div>
      <Link to="/profile" className="back-link">Назад</Link>
    </div>
  );
};

export default ConstructorPage;
