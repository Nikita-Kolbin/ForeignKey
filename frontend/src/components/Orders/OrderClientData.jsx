import React from 'react';

const OrderClientData = ({ clientData, orderDate }) => {
  return (
    <div>
      <h3>Данные клиента</h3>
      <p>ФИО: {clientData.fullName}</p>
      <p>Почта: {clientData.email}</p>
      <p>Способ оплаты: {clientData.paymentMethod}</p>
      <p>Доставка: {clientData.delivery}</p>
      <p>Телефон: {clientData.phone}</p>
      <p>Дата оформления заказа: {orderDate}</p>
      <p>Комментарий: {clientData.comment}</p>
    </div>
  );
};

export default OrderClientData;
