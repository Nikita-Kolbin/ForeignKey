import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../ApiConfig';
import '../styles/MessageSender.css';

const MessageSender = () => {
  const [messengerNotification, setMessengerNotification] = useState(false);
  const [emailNotification, setEmailNotification] = useState(false);

  useEffect(() => {
    // Fetch initial states from localStorage
    const initialMessengerStatus = localStorage.getItem('telegramNotification') === 'true';
    const initialEmailStatus = localStorage.getItem('emailNotification') === 'true';
    
    setMessengerNotification(initialMessengerStatus);
    setEmailNotification(initialEmailStatus);
  }, []);

  const handleMessengerToggle = async () => {
    const newStatus = !messengerNotification;
    setMessengerNotification(newStatus);
    try {
      const response = await fetch(`${API_BASE_URL}/admin/set-telegram-notification`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ notification: newStatus ? 1 : 0 })
      });

      if (!response.ok) {
        throw new Error('Ошибка обновления статуса уведомлений в мессенджерах');
      }

      localStorage.setItem('telegramNotification', newStatus);
    } catch (error) {
      console.error('Ошибка обновления статуса уведомлений в мессенджерах:', error);
      setMessengerNotification(!newStatus); // Revert the state in case of an error
    }
  };

  const handleEmailToggle = async () => {
    const newStatus = !emailNotification;
    setEmailNotification(newStatus);
    try {
      const response = await fetch(`${API_BASE_URL}/admin/set-email-notification`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ notification: newStatus ? 1 : 0 })
      });

      if (!response.ok) {
        throw new Error('Ошибка обновления статуса уведомлений по эл. почте');
      }

      localStorage.setItem('emailNotification', newStatus);
    } catch (error) {
      console.error('Ошибка обновления статуса уведомлений по эл. почте:', error);
      setEmailNotification(!newStatus); // Revert the state in case of an error
    }
  };

  return (
    <div className="message-sender">
      <h3>Рассылки</h3>
      <div className="toggle-container">
        <label>
          Уведомления по эл. почте
          <div className="toggle-switch">
            <input type="checkbox" checked={emailNotification} onChange={handleEmailToggle} />
            <span className="slider round"></span>
          </div>
        </label>
      </div>
      <div className="toggle-container">
        <label>
          Уведомления в мессенджерах
          <div className="toggle-switch">
            <input type="checkbox" checked={messengerNotification} onChange={handleMessengerToggle} />
            <span className="slider round"></span>
          </div>
        </label>
      </div>
    </div>
  );
};

export default MessageSender;
