import React, { useState } from 'react';

const EditProfileForm = ({ userData, onSaveChanges }) => {
  const [editedData, setEditedData] = useState({ ...userData });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEditedData({ ...editedData, [name]: value });
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      setEditedData({ ...editedData, photo: imageUrl });
    }
  };

  const handleSave = () => {
    onSaveChanges(editedData);
  };

  return (
    <div className="edit-profile-form">
      <div className="photo-section">
        <label className="image-upload-label">
          Изображение:
          <div className="image-upload-container">
            <input type="file" accept="image/*" onChange={handleImageChange} />
            <span>Вставьте изображение</span>
          </div>
        </label>
      </div>
      <div className="info-section">
        <div>
          <label>

            <input type="text" name="first_name" value={editedData.first_name} onChange={handleInputChange} className="personal-input" placeholder="Имя"/>
          </label>
          <label>

            <input type="text" name="last_name" value={editedData.last_name} onChange={handleInputChange} className="personal-input" placeholder="Фамилия"/>
          </label>
        </div>
        <div>
          <label>

            <input type="text" name="father_name" value={editedData.father_name} onChange={handleInputChange} className="personal-input" placeholder="Отчество"/>
          </label>
          <label>

            <input type="text" name="city" value={editedData.city} onChange={handleInputChange} className="personal-input" placeholder="Город"/>
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
