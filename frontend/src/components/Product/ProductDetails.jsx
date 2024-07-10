import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../ApiConfig';

const ProductDetails = ({ alias }) => {
  const [products, setProducts] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/product/get-all-by-alias/${alias}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          setProducts(data.products);
        } else {
          setError(data.error || 'Ошибка при получении товаров');
        }
      } catch (error) {
        console.error('Ошибка при получении товаров:', error);
        setError('Ошибка при получении товаров. Попробуйте снова позже.');
      }
    };

    fetchProducts();
  }, [alias]);

  const handleStatusToggle = async (product) => {
    try {
      const response = await fetch(`${API_BASE_URL}/product/set-active`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          alias: alias,
          product_id: product.id,
          active: product.active ? 0 : 1
        })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setProducts(prevProducts => prevProducts.map(p =>
          p.id === product.id ? { ...p, active: !p.active } : p
        ));
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при изменении статуса товара:', error);
      setError('Ошибка при изменении статуса товара. Попробуйте снова позже.');
    }
  };

  const fetchImage = async (id) => {
    try {
      const response = await fetch(`${API_BASE_URL}/image/download/${id}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      if (response.ok) {
        const imageBlob = await response.blob();
        return URL.createObjectURL(imageBlob);
      } else {
        setError('Ошибка загрузки изображения');
        return '';
      }
    } catch (error) {
      console.error('Ошибка загрузки изображения:', error);
      setError('Ошибка загрузки изображения');
      return '';
    }
  };

  const renderImages = (imagesId) => {
    const imageIdsArray = imagesId.split(' ');
    return imageIdsArray.map((id, index) => {
      const [imageUrl, setImageUrl] = useState('');

      useEffect(() => {
        const getImageUrl = async () => {
          const url = await fetchImage(id);
          setImageUrl(url);
        };

        getImageUrl();
      }, [id]);

      return imageUrl && (
        <div key={index} className="image-container">
          <img src={imageUrl} alt={`Product ${index}`} />
        </div>
      );
    });
  };

  return (
    <div className="product-list-container">
      {error && <div className="error-message">{error}</div>}
      {products.map((product) => (
        <div key={product.id} className="product-details-container">
          <h3>{product.name}</h3>
          <p><strong>Описание:</strong> {product.description}</p>
          <p><strong>Цена:</strong> {product.price}</p>
          {product.images_id && renderImages(product.images_id)}
          <button onClick={() => handleStatusToggle(product)}>
            {product.active ? 'Скрыть товар' : 'Показать товар'}
          </button>
        </div>
      ))}
    </div>
  );
};

export default ProductDetails;
