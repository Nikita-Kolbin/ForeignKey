import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const data = [
  { name: 'Янв', revenue: 4000 },
  { name: 'Фев', revenue: 3000 },
  { name: 'Мар', revenue: 5000 },
  { name: 'Апр', revenue: 6000 },
  { name: 'Май', revenue: 8000 },
  { name: 'Июн', revenue: 7000 },
  { name: 'Июл', revenue: 9000 },
];

const RevenueChart = () => {
  return (
    <ResponsiveContainer width={800} height={300}>
      <LineChart data={data}>
        <XAxis dataKey="name" />
        <YAxis />
        <CartesianGrid strokeDasharray="3 3" />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="revenue" stroke="#8884d8" />
      </LineChart>
    </ResponsiveContainer>
  );
};

export default RevenueChart;
