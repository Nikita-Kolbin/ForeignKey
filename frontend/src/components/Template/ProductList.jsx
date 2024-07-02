import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../ApiConfig';


const ProductList = ({ alias, onAddToCart }) => {
  const [products, setProducts] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/product/get-by-alias/${alias}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          setProducts(data.products);
        } else {
          setError(data.error || 'Неизвестная ошибка');
        }
      } catch (error) {
        console.error('Ошибка при получении списка товаров:', error);
        setError('Ошибка при получении списка товаров. Попробуйте снова позже.');
      }
    };

    fetchProducts();
  }, [alias]);

  const renderProductRow = (product) => (
    <tr key={product.id}>
      <td>{product.name}</td>
      <td>{product.category}</td>
      <td>{product.price}</td>
      <td>
        {product.images_id && product.images_id.split(' ').map(imageId => (
          <img 
            key={imageId}
            src={`${API_BASE_URL}/image/download/${imageId}`} 
            alt={product.name} 
            className="product-image" 
            onError={(e) => e.target.src = '/path/to/placeholder/image.jpg'} // Замените на путь к изображению-заполнителю
          />
        ))}
      </td>
      <td>
        <button onClick={() => onAddToCart(product)}>Добавить в корзину</button>
      </td>
    </tr>
  );

  return (
    <div className="product-list">
      <h2>Список товаров</h2>
      {error && <div className="error-message">{error}</div>}
      <table className="product-table">
        <thead>
          <tr>
            <th>Название</th>
            <th>Категория</th>
            <th>Цена</th>
            <th>Изображения</th>
            <th>Действие</th>
          </tr>
        </thead>
        <tbody>
          {products.map(renderProductRow)}
        </tbody>
      </table>
    </div>
  );
};

export default ProductList;
