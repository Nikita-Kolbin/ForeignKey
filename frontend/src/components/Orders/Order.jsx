import React from 'react';

const Order = ({ order, children }) => {
  return (
    <tr>
      <td>{order.id}</td>
      <td>{order.client}</td>
      <td>{order.amount}</td>
      <td>{order.date}</td>
      <td>
        <button>{order.status === 'active' ? 'Активный' : 'Закрытый'}</button>
      </td>
      <td>
        <button>Подробнее</button>
        {children} {/* Кнопка "Показать открытые/закрытые заказы" */}
      </td>
    </tr>
  );
};

export default Order;
