import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const data = [
  { name: 'Янв', orders: 20 },
  { name: 'Фев', orders: 25 },
  { name: 'Мар', orders: 30 },
  { name: 'Апр', orders: 35 },
  { name: 'Май', orders: 40 },
  { name: 'Июн', orders: 45 },
  { name: 'Июл', orders: 50 },
];

const OrdersChart = () => {
  return (
    <ResponsiveContainer width={800} height={300}>
      <LineChart data={data}>
        <XAxis dataKey="name" />
        <YAxis />
        <CartesianGrid strokeDasharray="3 3" />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="orders" stroke="#82ca9d" />
      </LineChart>
    </ResponsiveContainer>
  );
};

export default OrdersChart;
