import React from 'react';

const OrderProductsTable = ({ products }) => {
  return (
    <table>
      <thead>
        <tr>
          <th>Фото</th>
          <th>ID</th>
          <th>Наименование</th>
          <th>Цена</th>
          <th>Кол-во товара</th>
        </tr>
      </thead>
      <tbody>
        {products.map(product => (
          <tr key={product.id}>
            <td>Фото</td>
            <td>{product.id}</td>
            <td>{product.name}</td>
            <td>{product.price}</td>
            <td>{product.quantity}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default OrderProductsTable;
