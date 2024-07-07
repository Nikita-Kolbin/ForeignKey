import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { API_BASE_URL } from '../ApiConfig';


const OrderProductsTable = ({ products }) => {
  const [images, setImages] = useState({});
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchImages = async () => {
      const newImages = {};
      for (const product of products) {
        const imageIds = product.images_id.split(' '); // Разделяем идентификаторы по пробелам
        const firstImageId = imageIds[0]; // Берем только первый идентификатор

        try {
          const encodedImageId = encodeURIComponent(firstImageId);
          const response = await axios.get(`${API_BASE_URL}/image/download/${encodedImageId}`, {
            responseType: 'blob'
          });
          const imageUrl = URL.createObjectURL(response.data);
          newImages[product.id] = imageUrl;
        } catch (error) {
          setError(`Error fetching image for product ${product.id}: ${error.response?.data?.error || error.message}`);
        }
      }
      setImages(newImages);
    };

    fetchImages();
  }, [products]);

  return (
    <div>
      {error && <p className="order-product-error">{error}</p>}
      <table className="order-product-table">
        <thead>
          <tr>
            <th>Фото</th>
            <th>ID</th>
            <th>Наименование</th>
            <th>Цена</th>
            <th>Кол-во товара</th>
          </tr>
        </thead>
        <tbody>
          {products.map(product => (
            <tr key={product.id}>
              <td>
              <div className="order-product-image-container">
                {images[product.id] ? (
                  <img src={images[product.id]} alt={product.name} className="order-product-image" />
                ) : (
                  'Загрузка...'
                )}
              </div>
            </td>

              <td>{product.id}</td>
              <td>{product.name}</td>
              <td>{product.price} руб</td>
              <td>{product.quantity}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default OrderProductsTable;
