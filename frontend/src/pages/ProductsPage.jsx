import React, { useState, useEffect } from 'react';
import NavigationControlPanel from '../components/NavigationControlPanel';
import ProductDetails from '../components/Product/ProductDetails';
import AddProductForm from '../components/Product/AddProductForm';
import '../styles/ProductsPage.css';
import { useParams } from 'react-router-dom';
import { API_BASE_URL } from '../components/ApiConfig';

const ProductPage = () => {
  const { alias } = useParams();
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [products, setProducts] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [siteName, setSiteName] = useState('');
  const [siteAlias, setSiteAlias] = useState(alias || '');
  const [error, setError] = useState('');

  useEffect(() => {
    if (!siteAlias) {
      fetchSiteName();
    } else {
      fetchProducts(siteAlias);
    }
  }, [siteAlias]);

  const fetchSiteName = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/website/aliases`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK' && data.aliases.length > 0) {
        const firstAlias = data.aliases[0];
        setSiteName(firstAlias);
        setSiteAlias(firstAlias);
        fetchProducts(firstAlias);
      } else {
        setError(data.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      console.error('Ошибка при получении названия сайта:', error);
      setError('Ошибка при получении названия сайта. Попробуйте снова позже.');
    }
  };

  const fetchProducts = async (alias) => {
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

  const handleProductClick = (product) => {
    setSelectedProduct(product);
  };

  const handleCloseProductDetails = () => {
    setSelectedProduct(null);
  };

  const handleShowAddForm = () => {
    setShowAddForm(true);
  };

  const handleCloseAddForm = () => {
    setShowAddForm(false);
  };

  const filteredProducts = products.filter(product =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const renderProductRow = (product) => (
    <tr key={product.id} onClick={() => handleProductClick(product)}>
      <td>{product.name}</td>
      <td>{product.description}</td>
      <td>{product.price}</td>
      <td>
        {product.images_id && product.images_id.split(' ').map(imageId => (
          <img 
            key={imageId}
            src={`${API_BASE_URL}/image/download/${imageId}`} 
            alt={product.name} 
            className="product-image" 
            onError={(e) => e.target.src = '/path/to/placeholder/image.jpg'} 
          />
        ))}
      </td>
    </tr>
  );

  return (
    <div>
      <NavigationControlPanel />
      <div className="product-list-container">
        <h2>Список товаров для {siteName}</h2>
        <div className="search-container">
          <input
            type="text"
            placeholder="Поиск..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="search-input"
          />
          <button className="product-details-button" onClick={handleShowAddForm}>Добавить товар</button>
        </div>
        <table className="product-table">
          <thead>
            <tr>
              <th>Название</th>
              <th>Описание</th>
              <th>Цена</th>
              <th>Изображения</th>
            </tr>
          </thead>
          <tbody>
            {filteredProducts.map(renderProductRow)}
          </tbody>
        </table>
        {selectedProduct && (
          <div className="modal-wrapper" onClick={handleCloseProductDetails}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <span className="close" onClick={handleCloseProductDetails}>&times;</span>
              <ProductDetails product={selectedProduct} alias={siteAlias} />
            </div>
          </div>
        )}
        {showAddForm && (
          <div className="modal-wrapper" onClick={handleCloseAddForm}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <span className="close" onClick={handleCloseAddForm}>&times;</span>
              <AddProductForm siteAlias={siteAlias} />
            </div>
          </div>
        )}
      </div>
      {error && <div className="error-message">{error}</div>}
    </div>
  );
};

export default ProductPage;
