import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import '../styles/Charts.css';

const Chart = () => {
  const data = [
    { name: 'Ноябрь', доход: 45 },
    { name: 'Декабрь', доход: 30 },
    { name: 'Январь', доход: 60 },
    { name: 'Февраль', доход: 50 },
    { name: 'Март', доход: 40 },
    { name: 'Апрель', доход: 70 },
  ];

  return (
    <div className="chart-container">
      <h3>Динамика роста заказов</h3>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="name" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="доход" stroke="#8884d8" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default Chart;
