import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../ApiConfig'; // Импортируйте базовый URL

const EditProfileForm = ({ userData, onSaveChanges }) => {
  const [editedData, setEditedData] = useState(userData);
  const [selectedImageId, setSelectedImageId] = useState(userData.image_id || null);
  const [imageSrc, setImageSrc] = useState(null);

  useEffect(() => {
    if (selectedImageId) {
      fetchImage(selectedImageId);
    }
  }, [selectedImageId]);

  const fetchImage = async (id) => {
    try {
      const response = await fetch(`${API_BASE_URL}/image/download/${id}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}` // Добавьте токен авторизации
        }
      });

      if (response.ok) {
        const blob = await response.blob();
        const imageUrl = URL.createObjectURL(blob);
        setImageSrc(imageUrl);
      } else {
        const result = await response.json();
        console.error('Ошибка получения изображения', result.error);
      }
    } catch (error) {
      console.error('Ошибка получения изображения', error);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEditedData({ ...editedData, [name]: value });
  };

  const handleImageChange = async (e) => {
    const imageFile = e.target.files[0];
    if (imageFile) {
      const reader = new FileReader();
      reader.readAsArrayBuffer(imageFile);
      reader.onloadend = async () => {
        const byteArray = new Uint8Array(reader.result);
        try {
          const response = await fetch(`${API_BASE_URL}/image/upload`, {
            method: 'POST',
            headers: {
              'Content-Type': 'image/jpeg',
              'Authorization': `Bearer ${localStorage.getItem('token')}` // Добавьте токен авторизации
            },
            body: byteArray
          });

          const result = await response.json();
          if (result.id) {
            setSelectedImageId(result.id);
            setEditedData({ ...editedData, image_id: result.id });
            fetchImage(result.id);
          } else {
            console.error('Ошибка загрузки изображения', result.error);
          }
        } catch (error) {
          console.error('Ошибка загрузки изображения', error);
        }
      };
      reader.onerror = (error) => {
        console.error('Ошибка чтения файла изображения', error);
      };
    }
  };

  const handleSave = async () => {
    const { first_name, last_name, father_name, city, image_id, telegram } = editedData;
    const updatedUserData = { first_name, last_name, father_name, city, image_id: selectedImageId || image_id, telegram };
    try {
      const response = await fetch(`${API_BASE_URL}/admin/update-profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}` // Добавьте токен авторизации
        },
        body: JSON.stringify(updatedUserData)
      });

      const result = await response.json();
      if (result.status === 'OK') {
        onSaveChanges(updatedUserData);
        window.location.reload(); // Обновить страницу
      } else {
        console.error('Ошибка обновления профиля', result.error);
      }
    } catch (error) {
      console.error('Ошибка обновления профиля', error);
    }
  };

  return (
    <div className="edit-profile-form">
      <div className="photo-section">
        <label className="image-upload-label">
          Изображение:
          <div className="image-upload-container">
            <input type="file" accept="image/*" onChange={handleImageChange} className="view-image"/>
            {imageSrc ? (
              <img src={imageSrc} alt="Загруженное изображение" className="uploaded-image" />
            ) : (
              <span>Вставьте изображение</span>
            )}
          </div>
        </label>
      </div>
      <div className="info-section">
        <div>
          <label>
            <input type="text" name="first_name" value={editedData.first_name} onChange={handleInputChange} className="personal-input" placeholder="Имя" />
          </label>
          <label>
            <input type="text" name="last_name" value={editedData.last_name} onChange={handleInputChange} className="personal-input" placeholder="Фамилия" />
          </label>
        </div>
        <div>
          <label>
            <input type="text" name="father_name" value={editedData.father_name} onChange={handleInputChange} className="personal-input" placeholder="Отчество" />
          </label>
          <label>
            <input type="text" name="city" value={editedData.city} onChange={handleInputChange} className="personal-input" placeholder="Город" />
          </label>
        </div>
        <div>
          <label>
            <input type="text" name="telegram" value={editedData.telegram} onChange={handleInputChange} className="personal-input" placeholder="Telegram" />
          </label>
        </div>
        <div className="button-container">
          <button onClick={handleSave}>Сохранить</button>
        </div>
      </div>
    </div>
  );
};

export default EditProfileForm;
