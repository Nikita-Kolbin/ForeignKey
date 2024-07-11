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
  const [deliveryType, setDeliveryType] = useState('курьер');
  const [fullName, setFullName] = useState('');
  const [paymentType, setPaymentType] = useState('наличные');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [telegramName, setTelegramName] = useState('');
  const [comment, setComment] = useState('');
  const [emailNotification, setEmailNotification] = useState(false);
  const [telegramNotification, setTelegramNotification] = useState(false);

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
            .filter(item => item.count > 0)
            .map(item => ({
              id: item.product.id,
              name: item.product.name,
              price: item.product.price,
              quantity: item.count
            }));
          setCartItems(items);
          setCustomerLoggedIn(true);
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

  const fetchCustomerProfile = async () => {
    const token = localStorage.getItem('customerToken');
    if (token) {
      try {
        const response = await fetch(`${API_BASE_URL}/customer/get-profile`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          const profile = data.profile;
          setDeliveryType(profile.delivery_type || 'курьер');
          let fullName = '';
          if (profile.last_name) fullName += profile.last_name;
          if (profile.first_name) fullName += fullName ? ` ${profile.first_name}` : profile.first_name;
          if (profile.father_name) fullName += fullName ? ` ${profile.father_name}` : profile.father_name;
          setFullName(fullName);
          setPaymentType(profile.payment_type || 'наличные');
          setPhoneNumber(profile.phone || '');
          setTelegramName(profile.telegram || '');
          setEmailNotification(profile.email_notification || false);
          setTelegramNotification(profile.telegram_notification || false);
        } else {
          console.error('Ошибка при получении профиля клиента:', data.error);
        }
      } catch (error) {
        console.error('Ошибка при получении профиля клиента:', error);
      }
    }
  };

  useEffect(() => {
    if (customerLoggedIn) {
      fetchCustomerProfile();
    }
  }, [customerLoggedIn]);

  const handleSignInSuccess = () => {
    setCustomerLoggedIn(true);
    setShowSignIn(false);
    fetchCartItems();
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
        fetchCartItems();
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
        fetchCartItems();
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

    if (cartItems.length === 0) {
      alert('Корзина пуста. Добавьте товары в корзину перед созданием заказа.');
      return;
    }

    try {
      const profileResponse = await fetch(`${API_BASE_URL}/customer/update-profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          delivery_type: deliveryType,
          father_name: fullName.split(' ')[2] || '',
          first_name: fullName.split(' ')[1] || '',
          last_name: fullName.split(' ')[0] || '',
          payment_type: paymentType,
          phone: phoneNumber,
          telegram: telegramName,
          comment: comment // Include comment in the profile update request
        })
      });

      const profileData = await profileResponse.json();

      if (!profileResponse.ok || profileData.status !== 'OK') {
        console.error('Ошибка при обновлении профиля:', profileData.error);
        alert(`Ошибка при обновлении профиля: ${profileData.error}`);
        return;
      }

      const orderResponse = await fetch(`${API_BASE_URL}/order/make`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json', // Ensure content type is set
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          comment: comment // Include comment in the order creation request
        })
      });

      const orderData = await orderResponse.json();

      if (orderResponse.ok && orderData.status === 'OK') {
        setCartItems([]);
        alert('Заказ успешно создан');
      } else {
        console.error('Ошибка при создании заказа:', orderData.error);
        alert(`Ошибка при создании заказа: ${orderData.error}`);
      }
    } catch (error) {
      console.error('Ошибка при создании заказа:', error);
      alert('Ошибка при создании заказа. Попробуйте снова позже.');
    }
  };

  const handleEmailNotificationChange = async () => {
    const newStatus = !emailNotification;
    setEmailNotification(newStatus);
    const token = localStorage.getItem('customerToken');
    try {
      const response = await fetch(`${API_BASE_URL}/customer/set-email-notification`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ notification: newStatus ? 1 : 0 })
      });

      const data = await response.json();

      if (!response.ok || data.status !== 'OK') {
        console.error('Ошибка при изменении статуса email уведомлений:', data.error);
        setEmailNotification(!newStatus); // Revert state on error
      }
    } catch (error) {
      console.error('Ошибка при изменении статуса email уведомлений:', error);
      setEmailNotification(!newStatus); // Revert state on error
    }
  };

  const handleTelegramNotificationChange = async () => {
    const newStatus = !telegramNotification;
    setTelegramNotification(newStatus);
    const token = localStorage.getItem('customerToken');
    try {
      const response = await fetch(`${API_BASE_URL}/customer/set-telegram-notification`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ notification: newStatus ? 1 : 0 })
      });

      const data = await response.json();

      if (!response.ok || data.status !== 'OK') {
        console.error('Ошибка при изменении статуса telegram уведомлений:', data.error);
        setTelegramNotification(!newStatus); // Revert state on error
      }
    } catch (error) {
      console.error('Ошибка при изменении статуса telegram уведомлений:', error);
      setTelegramNotification(!newStatus); // Revert state on error
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
      />
      {customerLoggedIn && (
        <div>
          <h2>Оформление заказа</h2>
          <label>
            Тип доставки:
            <select value={deliveryType} onChange={(e) => setDeliveryType(e.target.value)}>
              <option value="курьер">Курьер</option>
              <option value="самовывоз">Самовывоз</option>
            </select>
          </label>
          <label>
            ФИО:
            <input
              type="text"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
            />
          </label>
          <label>
            Тип оплаты:
            <select value={paymentType} onChange={(e) => setPaymentType(e.target.value)}>
              <option value="наличные">Наличные</option>
              <option value="безнал">Безнал</option>
            </select>
          </label>
          <label>
            Номер телефона:
            <input
              type="text"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
            />
          </label>
          <label>
            Имя в телеграмм:
            <input
              type="text"
              value={telegramName}
              onChange={(e) => setTelegramName(e.target.value)}
            />
          </label>
          <label>
            Комментарий:
            <input
              type="text"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
          </label>
          <label>
            Уведомления по email:
            <input
              type="checkbox"
              checked={emailNotification}
              onChange={handleEmailNotificationChange}
            />
          </label>
          <label>
            Уведомления по Telegram:
            <input
              type="checkbox"
              checked={telegramNotification}
              onChange={handleTelegramNotificationChange}
            />
          </label>
          <button onClick={handleMakeOrder}>Сделать заказ</button>
        </div>
      )}
    </div>
  );
};

export default Template2;
