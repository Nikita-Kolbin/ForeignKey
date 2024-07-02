import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import CustomerSignInForm from '../Template/CustomerSignInForm';
import CustomerSignUpForm from '../Template/CustomerSignUpForm';
import ProductList from '../Template/ProductList';
import Cart from '../Template/Cart';
import { API_BASE_URL } from '../ApiConfig';

const Template2 = () => {
  const { alias } = useParams();
  const [styles, setStyles] = useState({ backgroundColor: '#ffffff', font: 'Arial' });
  const [showSignIn, setShowSignIn] = useState(false);
  const [showSignUp, setShowSignUp] = useState(false);
  const [customerLoggedIn, setCustomerLoggedIn] = useState(false);
  const [cartItems, setCartItems] = useState([]);

  useEffect(() => {
    const fetchStyles = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/website/get-style/${alias}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('customerToken')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          setStyles({
            backgroundColor: data.background_color,
            font: data.font,
          });
        }
      } catch (error) {
        console.error('Ошибка при получении стилей сайта:', error);
      }
    };

    fetchStyles();
  }, [alias]);

  const fetchCartItems = async () => {
    const token = localStorage.getItem('customerToken');
    if (token) {
      try {
        const response = await fetch(`${API_BASE_URL}/cart/get`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          const items = data.cart_items
            .filter(item => item.count > 0) // Фильтруем товары с количеством больше нуля
            .map(item => ({
              id: item.product.id,
              name: item.product.name,
              price: item.product.price,
              quantity: item.count
            }));
          setCartItems(items);
          setCustomerLoggedIn(true); // Устанавливаем состояние входа в систему, если корзина успешно загружена
        } else {
          console.error('Ошибка при получении товаров в корзине:', data.error);
        }
      } catch (error) {
        console.error('Ошибка при получении товаров в корзине:', error);
      }
    }
  };

  useEffect(() => {
    fetchCartItems();
  }, []);

  const handleSignInSuccess = () => {
    setCustomerLoggedIn(true);
    setShowSignIn(false);
    fetchCartItems(); // Загружаем корзину после успешного входа
  };

  const handleSignUpSuccess = () => {
    setShowSignIn(true);
    setShowSignUp(false);
  };

  const handleAddToCart = async (product) => {
    const token = localStorage.getItem('customerToken');
    try {
      const response = await fetch(`${API_BASE_URL}/cart/add`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ count: 1, product_id: product.id })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        fetchCartItems(); // Обновляем корзину после добавления товара
      } else {
        console.error('Ошибка при добавлении товара в корзину:', data.error);
      }
    } catch (error) {
      console.error('Ошибка при добавлении товара в корзину:', error);
    }
  };

  const handleChangeCartItemCount = async (product, newCount) => {
    const token = localStorage.getItem('customerToken');
    try {
      const response = await fetch(`${API_BASE_URL}/cart/change-count`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ new_count: newCount, product_id: product.id })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        fetchCartItems(); // Обновляем корзину после изменения количества товара
      } else {
        console.error('Ошибка при изменении количества товара в корзине:', data.error);
      }
    } catch (error) {
      console.error('Ошибка при изменении количества товара в корзине:', error);
    }
  };

  const handleRemoveFromCart = async (product) => {
    const newCount = product.quantity - 1;
    handleChangeCartItemCount(product, newCount);
  };

  const handleMakeOrder = async () => {
    const token = localStorage.getItem('customerToken');

    // Проверяем, что корзина не пуста
    if (cartItems.length === 0) {
      alert('Корзина пуста. Добавьте товары в корзину перед созданием заказа.');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/order/make`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setCartItems([]); // Обнуляем корзину после успешного создания заказа
        alert('Заказ успешно создан');
      } else {
        console.error('Ошибка при создании заказа:', data.error);
        alert(`Ошибка при создании заказа: ${data.error}`);
      }
    } catch (error) {
      console.error('Ошибка при создании заказа:', error);
      alert('Ошибка при создании заказа. Попробуйте снова позже.');
    }
  };

  return (
    <div
      className="template"
      style={{ backgroundColor: styles.backgroundColor, fontFamily: styles.font }}
    >
      <h1>Шаблон сайта: {alias}</h1>
      {!customerLoggedIn && (
        <div>
          <button onClick={() => setShowSignIn(true)}>Войти</button>
          <button onClick={() => setShowSignUp(true)}>Регистрация</button>
        </div>
      )}
      {showSignIn && (
        <CustomerSignInForm alias={alias} onSuccess={handleSignInSuccess} />
      )}
      {showSignUp && (
        <CustomerSignUpForm alias={alias} onSuccess={handleSignUpSuccess} />
      )}
      <ProductList alias={alias} onAddToCart={handleAddToCart} />
      <Cart 
        cartItems={cartItems} 
        onAdd={handleAddToCart} 
        onRemove={handleRemoveFromCart} 
        onMakeOrder={handleMakeOrder} 
      />
    </div>
  );
};

export default Template2;
