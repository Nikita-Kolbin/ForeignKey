import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { API_BASE_URL } from '../ApiConfig';

const OrdersChart = () => {
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
      const date = new Date(order.date_time);
      return { date: date.toLocaleDateString('ru-RU', formatOptions[timeFrame]), orders: 1 };
    });

    const groupedData = processedData.reduce((acc, curr) => {
      const existing = acc.find(item => item.name === curr.date);
      if (existing) {
        existing.orders += curr.orders;
      } else {
        acc.push({ name: curr.date, orders: curr.orders });
      }
      return acc;
    }, []);

    return groupedData;
  };

  const handleTimeFrameChange = (event) => {
    setTimeFrame(event.target.value);
  };

  return (
    <div>
      <div>
        <label htmlFor="timeFrameSelector">Выберите временной интервал:</label>
        <select id="timeFrameSelector" value={timeFrame} onChange={handleTimeFrameChange}>
          <option value="month">Месяц</option>
          <option value="day">День</option>
          <option value="hour">Час</option>
          <option value="minute">Минута</option>
        </select>
      </div>
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
    </div>
  );
};

export default OrdersChart;
