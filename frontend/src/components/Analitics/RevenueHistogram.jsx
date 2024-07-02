import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const data = [
  { name: 'Янв', revenue: 4000 },
  { name: 'Фев', revenue: 3000 },
  { name: 'Мар', revenue: 5000 },
  { name: 'Апр', revenue: 6000 },
  { name: 'Май', revenue: 8000 },
  { name: 'Июн', revenue: 7000 },
  { name: 'Июл', revenue: 9000 },
];

const RevenueHistogram = () => {
  return (
    <ResponsiveContainer width={800} height={300}>
      <BarChart data={data}>
        <XAxis dataKey="name" />
        <YAxis />
        <CartesianGrid strokeDasharray="3 3" />
        <Tooltip />
        <Legend />
        <Bar dataKey="revenue" fill="#8884d8" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default RevenueHistogram;
