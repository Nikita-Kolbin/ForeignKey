import React, { useState } from 'react';
import { API_BASE_URL } from '../ApiConfig';

const ProductDetails = ({ product, alias }) => {
  const [error, setError] = useState('');
  const [active, setActive] = useState(product.active);

  const handleStatusToggle = async () => {
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
          active: active ? 0 : 1
        })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setActive(prevActive => !prevActive);
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при изменении статуса товара:', error);
      setError('Ошибка при изменении статуса товара. Попробуйте снова позже.');
    }
  };

  return (
    <div className="product-details-container">
      <h3>{product.name}</h3>
      <p><strong>Описание:</strong> {product.description}</p>
      <p><strong>Цена:</strong> {product.price}</p>
      <p><strong>ID изображения:</strong> {product.image_id}</p>
      <button onClick={handleStatusToggle}>
        {active ? 'Скрыть товар' : 'Показать товар'}
      </button>
      {error && <div className="error-message">{error}</div>}
    </div>
  );
};

export default ProductDetails;
