import React, { useState } from 'react';
import '../styles/MessageSender.css';

const MessageSender = () => {
  const [messengerNotification, setMessengerNotification] = useState(false);
  const [emailNotification, setEmailNotification] = useState(false);

  const handleMessengerToggle = () => {
    setMessengerNotification(!messengerNotification);
  };

  const handleEmailToggle = () => {
    setEmailNotification(!emailNotification);
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
