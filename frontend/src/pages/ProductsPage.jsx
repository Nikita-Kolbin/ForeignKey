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
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(5);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [productToDelete, setProductToDelete] = useState(null);

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

  const handlePageChange = (direction) => {
    if (direction === 'next') {
      setCurrentPage((prevPage) => prevPage + 1);
    } else {
      setCurrentPage((prevPage) => prevPage - 1);
    }
  };

  const handleItemsPerPageChange = (event) => {
    setItemsPerPage(parseInt(event.target.value));
    setCurrentPage(1);
  };

  const handleDeleteProduct = (product, event) => {
    event.stopPropagation();
    setProductToDelete(product);
    setIsDeleteModalOpen(true);
  };

  const confirmDeleteProduct = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/product/delete/${productToDelete.id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (response.ok) {
        setProducts((prevProducts) => prevProducts.filter((p) => p.id !== productToDelete.id));
        setIsDeleteModalOpen(false);
      } else {
        const data = await response.json();
        setError(data.error || 'Не удалось удалить продукт');
      }
    } catch (error) {
      console.error('Ошибка при удалении продукта:', error);
      setError('Ошибка при удалении продукта. Попробуйте снова позже.');
    }
  };

  const handleToggleProductVisibility = async (product, event) => {
    event.stopPropagation();
    try {
      const response = await fetch(`${API_BASE_URL}/product/set-active`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          alias: siteAlias,
          product_id: product.id,
          active: product.active ? 0 : 1
        })
      });

      if (response.ok) {
        fetchProducts(siteAlias); // Reload products to reflect the change
      } else {
        console.error('Failed to update product visibility');
      }
    } catch (error) {
      console.error('Error updating product visibility:', error);
    }
  };

  const filteredProducts = products.filter(product =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const paginatedProducts = filteredProducts.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

  const renderProductRow = (product) => {
    const imageIds = product.images_id ? product.images_id.split(' ') : [];
    const firstImageId = imageIds[0];

    return (
      <tr key={product.id} onClick={() => handleProductClick(product)}>
        <td>
          {firstImageId && (
            <img 
              src={`${API_BASE_URL}/image/download/${firstImageId}`} 
              alt={product.name} 
              className="product-image" 
              onError={(e) => e.target.src = '/path/to/placeholder/image.jpg'} 
            />
          )}
        </td>
        <td>{product.id}</td>
        <td>{product.name}</td>
        <td>{product.price} руб</td>
        <td>
          <button onClick={(event) => handleToggleProductVisibility(product, event)}>
            {product.active 
              ? <img src="/Скрыть.svg" alt="Скрыть" /> 
              : <img src="/Показать.svg" alt="Показать" />
            }
          </button>
          <button onClick={(event) => handleDeleteProduct(product, event)}>
            <img src="/Удалить.svg" alt="Удалить" />
          </button>
        </td>
      </tr>
    );
  };

  return (
    <div>
      <NavigationControlPanel />
      <div className="search-container">
        <div>
          <button>Каталог</button>
          <button className="product-details-button" onClick={handleShowAddForm}>Добавить</button>
        </div>
      </div>
      <div className="product-list-container">
        {error && <p className="error">{error}</p>}
        <table className="product-table">
          <thead>
            <tr>
              <th>Вид</th>
              <th>ID</th>
              <th>Наименование</th>
              <th>Цена</th>
              <th>
                <input
                  type="text"
                  placeholder="Поиск..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="search-input"
                />
              </th>
            </tr>
          </thead>
          <tbody>
            {paginatedProducts.map(renderProductRow)}
          </tbody>
        </table>

        <div className="pagination">
          <button onClick={() => handlePageChange('prev')} disabled={currentPage === 1}>
            <img src="/Назад.svg" alt="назад" />
          </button>
          <select value={itemsPerPage} onChange={handleItemsPerPageChange}>
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={15}>15</option>
          </select>
          <button onClick={() => handlePageChange('next')} disabled={currentPage === Math.ceil(filteredProducts.length / itemsPerPage)}>
            <img src="/Вперед.svg" alt="Вперед" />
          </button>
        </div>
        
        {showAddForm && (
          <div className="modal-wrapper" onClick={handleCloseAddForm}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <span className="close" onClick={handleCloseAddForm}>&times;</span>
              <AddProductForm siteAlias={siteAlias} />
            </div>
          </div>
        )}
        {isDeleteModalOpen && (
          <div className="modal-overlay" onClick={() => setIsDeleteModalOpen(false)}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <h3>Вы правда хотите удалить товар {productToDelete?.name}?</h3>
              <button onClick={confirmDeleteProduct} className="confirm-delete-button">Удалить</button>
              <button onClick={() => setIsDeleteModalOpen(false)} className="cancel-button">Отмена</button>
              {error && <div className="error-message">{error}</div>}
            </div>
          </div>
        )}
      </div>
      {error && <div className="error-message">{error}</div>}
    </div>
  );
};

export default ProductPage;
