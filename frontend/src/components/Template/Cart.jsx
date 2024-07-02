import React from 'react';

const Cart = ({ cartItems, onAdd, onRemove, onMakeOrder }) => {
  const getTotalItems = () => {
    return cartItems.reduce((total, item) => total + item.quantity, 0);
  };

  const getTotalPrice = () => {
    return cartItems.reduce((total, item) => total + item.price * item.quantity, 0);
  };

  return (
    <div className="cart">
      <h2>Корзина</h2>
      <div>Всего товаров: {getTotalItems()}</div>
      <div>Общая стоимость: {getTotalPrice().toFixed(2)}</div>
      <table className="cart-table">
        <thead>
          <tr>
            <th>Название</th>
            <th>Количество</th>
            <th>Цена</th>
            <th>Действие</th>
          </tr>
        </thead>
        <tbody>
          {cartItems.map(item => (
            <tr key={item.id}>
              <td>{item.name}</td>
              <td>{item.quantity}</td>
              <td>{(item.price * item.quantity).toFixed(2)}</td>
              <td>
                <button onClick={() => onAdd(item)}>+</button>
                <button onClick={() => onRemove(item)}>-</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <button onClick={onMakeOrder}>Сделать заказ</button>
    </div>
  );
};

export default Cart;
