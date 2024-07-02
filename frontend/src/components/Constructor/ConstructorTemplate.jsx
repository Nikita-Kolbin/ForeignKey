import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import EditorPanel from './EditorPanel';
import CustomerSignInForm from '../Template/CustomerSignInForm';
import CustomerSignUpForm from '../Template/CustomerSignUpForm';
import ProductList from '../Template/ProductList';
import Cart from '../Template/Cart'; // Импортируем новый компонент
import { API_BASE_URL } from '../ApiConfig';

const ConstructorTemplate = () => {
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
            'Authorization': `Bearer ${localStorage.getItem('token')}`
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

  const handleSaveStyles = async (styleData) => {
    try {
      const response = await fetch(`${API_BASE_URL}/website/set-style`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(styleData)
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setStyles({
          backgroundColor: styleData.background_color,
          font: styleData.font,
        });
      }
    } catch (error) {
      console.error('Ошибка при сохранении стилей сайта:', error);
    }
  };

  const handleSignInSuccess = () => {
    setCustomerLoggedIn(true);
    setShowSignIn(false);
  };

  const handleSignUpSuccess = () => {
    setShowSignIn(true);
    setShowSignUp(false);
  };

  const handleAddToCart = (product) => {
    setCartItems(prevItems => {
      const existingItem = prevItems.find(item => item.id === product.id);
      if (existingItem) {
        return prevItems.map(item =>
          item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
        );
      } else {
        return [...prevItems, { ...product, quantity: 1 }];
      }
    });
  };

  const handleRemoveFromCart = (product) => {
    setCartItems(prevItems => {
      const existingItem = prevItems.find(item => item.id === product.id);
      if (existingItem.quantity === 1) {
        return prevItems.filter(item => item.id !== product.id);
      } else {
        return prevItems.map(item =>
          item.id === product.id ? { ...item, quantity: item.quantity - 1 } : item
        );
      }
    });
  };

  const handleNavigate = () => {
    window.open(`${window.location.origin}/${alias}`, '_blank');
  };

  return (
    <div
      className="template"
      style={{ backgroundColor: styles.backgroundColor, fontFamily: styles.font }}
    >
      <h1>Конструктор сайта: {alias}</h1>
      <EditorPanel alias={alias} onSave={handleSaveStyles} onNavigate={handleNavigate} />
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
      {/* Выводим список товаров */}
      <ProductList alias={alias} onAddToCart={handleAddToCart} />
      {/* Выводим корзину */}
      <Cart cartItems={cartItems} onAdd={handleAddToCart} onRemove={handleRemoveFromCart} />
    </div>
  );
};

export default ConstructorTemplate;
