import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../components/ApiConfig';
import "../styles/ClientListPage.css"; // Стили для страницы клиентов
import NavigationControlPanel from '../components/NavigationControlPanel';
import useTitle from '../components/customTitle';

const ClientListPage = () => {
  const [showClosedOrders, setShowClosedOrders] = useState(false);
  const [clients, setClients] = useState([]);
  const [alias, setAlias] = useState('');
  const [error, setError] = useState('');

  useTitle('Клиенты')

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
          setAlias(data.aliases[0]);
        } else {
          setError(data.error || 'Ошибка при получении сайтов');
        }
      } catch (error) {
        setError('Ошибка при получении сайтов');
      }
    };

    fetchAlias();
  }, []);

  useEffect(() => {
    const fetchClients = async () => {
      if (alias) {
        try {
          const response = await fetch(`${API_BASE_URL}/customer/get-by-alias/${alias}`, {
            method: 'GET',
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          });

          const data = await response.json();

          if (response.ok && data.status === 'OK') {
            setClients(data.customers);
          } else {
            setError(data.error || 'Ошибка при получении клиентов');
          }
        } catch (error) {
          setError('Ошибка при получении клиентов');
        }
      }
    };

    fetchClients();
  }, [alias]);

  const renderClientsTable = () => (
    <table className="clients-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Фамилия</th>
          <th>Имя</th>
          <th>Отчество</th>
          <th>Email</th>
          <th>Телефон</th>
          <th>Telegram</th>
          <th>Способ оплаты</th>
          <th>Способ доставки</th>
        </tr>
      </thead>
      <tbody>
        {clients.map(client => (
          <tr key={client.id}>
            <td>{client.id}</td>
            <td>{client.last_name}</td>
            <td>{client.first_name}</td>
            <td>{client.father_name}</td>
            <td>{client.email}</td>
            <td>{client.phone}</td>
            <td>{client.telegram}</td>
            <td>{client.payment_type}</td>
            <td>{client.delivery_type}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );

  return (
    <div className="clients-page">
      <NavigationControlPanel
        showClosedOrders={showClosedOrders}
        setShowClosedOrders={setShowClosedOrders}
      />
      {error && <div className="error">{error}</div>}
      <div className="clients-section">
        {renderClientsTable()}
      </div>
    </div>
  );
};

export default ClientListPage;
