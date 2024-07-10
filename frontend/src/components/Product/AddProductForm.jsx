import React, { useState } from 'react';
import { API_BASE_URL } from '../ApiConfig';

const AddProductForm = ({ siteAlias }) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [imageIds, setImageIds] = useState([]);
  const [images, setImages] = useState([]);
  const [error, setError] = useState('');

  const handleImageChange = async (e) => {
    const files = e.target.files;
    const newImageIds = [];
    const newImages = [];

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const reader = new FileReader();

      reader.readAsArrayBuffer(file);
      reader.onloadend = async () => {
        const byteArray = new Uint8Array(reader.result);
        try {
          const response = await fetch(`${API_BASE_URL}/image/upload`, {
            method: 'POST',
            headers: {
              'Content-Type': 'image/jpeg',
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: byteArray
          });

          const result = await response.json();
          if (result.id) {
            newImageIds.push(result.id);
            setImageIds(prevIds => [...prevIds, result.id]);

            // Fetch the image to display
            const imgResponse = await fetch(`${API_BASE_URL}/image/download/${result.id}`, {
              method: 'GET',
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              }
            });
            if (imgResponse.ok) {
              const imageBlob = await imgResponse.blob();
              const imageObjectURL = URL.createObjectURL(imageBlob);
              newImages.push({ id: result.id, url: imageObjectURL });
              setImages(prevImages => [...prevImages, { id: result.id, url: imageObjectURL }]);
            }
          } else {
            setError(result.error || 'Ошибка загрузки изображения');
          }
        } catch (error) {
          setError('Ошибка загрузки изображения');
          console.error('Ошибка загрузки изображения:', error);
        }
      };
      reader.onerror = (error) => {
        setError('Ошибка чтения файла изображения');
        console.error('Ошибка чтения файла изображения:', error);
      };
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (imageIds.length === 0) {
      setError('Необходимо загрузить хотя бы одно изображение');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/product/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          alias: siteAlias,
          product_info: {
            name,
            description,
            price: parseFloat(price),
            images_id: imageIds.join(' ')
          }
        })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        window.location.reload();
      } else {
        setError(data.error || 'Ошибка при добавлении товара');
      }
    } catch (error) {
      console.error('Ошибка при добавлении товара:', error);
      setError('Ошибка при добавлении товара. Попробуйте снова позже.');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="add-product-form">
      <h3>Добавить товар</h3>
      <div className="form-group">
        <label>Название:</label>
        <input 
          type="text" 
          placeholder="Название" 
          value={name} 
          onChange={(e) => setName(e.target.value)} 
        />
      </div>
      <div className="form-group">
        <label>Описание:</label>
        <input 
          type="text" 
          placeholder="Описание" 
          value={description} 
          onChange={(e) => setDescription(e.target.value)} 
        />
      </div>
      <div className="form-group">
        <label>Цена:</label>
        <input 
          type="text" 
          placeholder="Цена" 
          value={price} 
          onChange={(e) => setPrice(e.target.value)} 
        />
      </div>
      <div className="form-group">
        <label>Изображения:</label>
        <div className="image-upload-area">
          <input 
            type="file" 
            accept="image/*" 
            multiple
            onChange={handleImageChange} 
          />
          <div className="uploaded-images">
            {images.map((image, index) => (
              <div key={index} className="uploaded-image">
                <img src={image.url} alt={`Uploaded ${image.id}`} />
                <span>ID: {image.id}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
      {error && <div className="error-message">{error}</div>}
      <button type="submit" className="product-details-button">Добавить</button>
    </form>
  );
};

export default AddProductForm;
