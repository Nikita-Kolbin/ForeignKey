import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import '../styles/Charts.css';
import { API_BASE_URL } from '../ApiConfig';

const Chart = () => {
  const [data, setData] = useState([]);
  const [alias, setAlias] = useState('');
  const [timeFrame, setTimeFrame] = useState('month'); // По умолчанию отображение по месяцам

  useEffect(() => {
    const fetchAlias = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/website/aliases`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK' && data.aliases.length > 0) {
          setAlias(data.aliases[0]); // Используем первый alias из списка
        } else {
          console.error('Ошибка при получении сайтов:', data.error);
        }
      } catch (error) {
        console.error('Ошибка при получении сайтов:', error);
      }
    };

    fetchAlias();
  }, []);

  useEffect(() => {
    if (alias) {
      const fetchOrdersByAlias = async () => {
        try {
          const response = await fetch(`${API_BASE_URL}/order/get-completed/${alias}`, {
            method: 'GET',
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          });
          const result = await response.json();

          if (response.ok && result.status === 'OK') {
            const completedOrders = result.orders.filter(order => order.status === 6);
            const processedData = processOrders(completedOrders, timeFrame);
            setData(processedData);
          } else {
            console.error(`Ошибка получения заказов для alias ${alias}:`, result.error);
          }
        } catch (error) {
          console.error(`Ошибка запроса заказов для alias ${alias}:`, error);
        }
      };

      fetchOrdersByAlias();
    }
  }, [alias, timeFrame]);

  const processOrders = (orders, timeFrame) => {
    const formatOptions = {
      month: { month: 'short' },
      day: { day: '2-digit', month: 'short' },
      hour: { hour: '2-digit', day: '2-digit', month: 'short' },
      minute: { minute: '2-digit', hour: '2-digit', day: '2-digit', month: 'short' }
    };

    const processedData = orders.map(order => {
      const totalAmount = order.order_items.reduce((sum, item) => sum + item.product.price * item.count, 0);
      const date = new Date(order.date_time);
      return { date: date.toLocaleDateString('ru-RU', formatOptions[timeFrame]), доход: totalAmount };
    });

    const groupedData = processedData.reduce((acc, curr) => {
      const existing = acc.find(item => item.name === curr.date);
      if (existing) {
        existing.доход += curr.доход;
      } else {
        acc.push({ name: curr.date, доход: curr.доход });
      }
      return acc;
    }, []);

    return groupedData;
  };

  const handleTimeFrameChange = (event) => {
    setTimeFrame(event.target.value);
  };

  return (
    <div className="chart-container">
      <h3>Динамика выручки</h3>
      <div>
        <label htmlFor="timeFrameSelector">Выберите временной интервал:</label>
        <select id="timeFrameSelector" value={timeFrame} onChange={handleTimeFrameChange}>
          <option value="month">Месяц</option>
          <option value="day">День</option>
          <option value="hour">Час</option>
          <option value="minute">Минута</option>
        </select>
      </div>
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
