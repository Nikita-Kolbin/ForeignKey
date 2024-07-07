import React from 'react';

const OrderClientData = ({ clientData, orderDate }) => {
  return (
    <div>
      <div className='information-client'>
        <p className='personal-info-client'>ФИО: <span>{clientData.fullName}</span></p>
        <p className='personal-info-client'>Доставка: <span>{clientData.delivery}</span></p>
      </div>
      <div className='information-client'>
        <p className='personal-info-client'>Почта: <span>{clientData.email}</span></p>
        <p className='personal-info-client'>Способ оплаты: <span>{clientData.paymentMethod}</span></p>
      </div>
      <div className='information-client'>
        <p className='personal-info-client'>Телефон: <span>{clientData.phone}</span></p>
        <p className='personal-info-client'>Имя в телеграмм: <span>{clientData.telegram}</span></p>
      </div>
      <p className='personal-info-client'>Дата оформления заказа: <span>{orderDate}</span></p>
    </div>
  );
};

export default OrderClientData;
