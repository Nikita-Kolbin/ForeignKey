import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { API_BASE_URL } from '../components/ApiConfig';
import NavigationControlPanel from '../components/NavigationControlPanel'; // Импортируем шапку навигации
import '../styles/ConstructorPage.css'; // Импортируем стили

const ConstructorPage = () => {
  const [aliases, setAliases] = useState([]);
  const [newAlias, setNewAlias] = useState('');
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [aliasToDelete, setAliasToDelete] = useState('');
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
    if (!/^[a-zA-Z0-9]+$/.test(newAlias)) {
      setError('Имя сайта должно содержать только английские буквы и цифры.');
      return;
    }

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
        setIsModalOpen(false);
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при создании нового сайта:', error);
      setError('Ошибка при создании нового сайта. Попробуйте снова позже.');
    }
  };

  const handleDeleteWebsite = (alias) => {
    setAliasToDelete(alias);
    setIsDeleteModalOpen(true);
  };

  const confirmDeleteWebsite = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/website/delete/${aliasToDelete}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        fetchAliases();
        setIsDeleteModalOpen(false);
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при удалении сайта:', error);
      setError('Ошибка при удалении сайта. Попробуйте снова позже.');
      setIsDeleteModalOpen(false);
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
      <NavigationControlPanel /> {/* Вставляем шапку навигации */}
      <h2>Страница конструктора</h2>
      <div className="constructor-content">
        {aliases.length > 0 ? (
          <div className="website-container">
            {aliases.map((alias, index) => (
              <div className="website-card" key={index}>
                <img src="/Фон-сайт.png" alt={alias} className="website-thumbnail" />
                <div className="website-actions">
                  <button className="action-button edit-button" onClick={() => handleSelectWebsite(alias)}>
                    <img src="/Конструктор.svg" alt="Редактировать" />
                  </button>
                  <button className="action-button delete-button" onClick={() => handleDeleteWebsite(alias)}>
                    <img src="/Удалить.png" alt="Удалить" />
                  </button>
                </div>
                <div className="website-info">
                  <h3>{alias}</h3>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <button className="add-website-button" onClick={() => setIsModalOpen(true)}>
              <span className="plus-icon">+</span>
            </button>
          </div>
        )}
        {isModalOpen && (
          <div className="modal-overlay">
            <div className="modal">
              <h3>Создать новый сайт</h3>
              <input
                type="text"
                value={newAlias}
                onChange={(e) => setNewAlias(e.target.value)}
                placeholder="Введите имя нового сайта"
                className="new-alias-input"
              />
              <button onClick={createWebsite} className="create-website-button">Создать</button>
              <button onClick={() => setIsModalOpen(false)} className="cancel-button">Отмена</button>
              {error && <div className="error-message">{error}</div>}
            </div>
          </div>
        )}
        {isDeleteModalOpen && (
          <div className="modal-overlay">
            <div className="modal">
              <h3>Вы правда хотите удалить сайт {aliasToDelete}?</h3>
              <button onClick={confirmDeleteWebsite} className="confirm-delete-button">Удалить</button>
              <button onClick={() => setIsDeleteModalOpen(false)} className="cancel-button">Отмена</button>
              {error && <div className="error-message">{error}</div>}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ConstructorPage;
