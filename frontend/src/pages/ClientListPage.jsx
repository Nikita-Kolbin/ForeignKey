import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../components/ApiConfig';
import NavigationControlPanel from '../components/NavigationControlPanel';
import ClientDetails from '../components/ClientDashboard/ClientDetails';
import '../styles/ClientListPage.css';

const ClientList = () => {
  const [selectedClient, setSelectedClient] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [clients, setClients] = useState([]);
  const [alias, setAlias] = useState('');

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
    if (!alias) return;

    const fetchClients = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/customer/get-by-alias/${alias}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          const clientsWithData = data.customers.map(client => ({
            id: client.id,
            email: client.email,
            // Добавьте другие поля, если они есть в ответе API
          }));
          setClients(clientsWithData);
        } else {
          console.error('Ошибка при получении клиентов:', data.error);
        }
      } catch (error) {
        console.error('Ошибка при получении клиентов:', error);
      }
    };

    fetchClients();
  }, [alias]);

  const handleClientClick = (client) => {
    setSelectedClient(client);
  };

  const handleCloseClientDetails = () => {
    setSelectedClient(null);
  };

  const renderClientRow = (client) => (
    <tr key={client.id} onClick={() => handleClientClick(client)}>
      <td>{client.email}</td>
      {/* Добавьте отображение других полей, если они есть */}
    </tr>
  );

  // Фильтрация клиентов по поисковому запросу
  const filteredClients = clients.filter((client) =>
    client.email.toLowerCase().includes(searchTerm.toLowerCase())
    // Добавьте другие условия фильтрации, если необходимо
  );

  return (
    <div>
      <NavigationControlPanel />
      <div className="client-list-container">
        <h2>Список клиентов</h2>
        <input 
          type="text" 
          placeholder="Поиск..." 
          className="search-input"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <table className="client-table">
          <thead>
            <tr>
              <th>Почта</th>
              {/* Добавьте другие заголовки столбцов, если необходимо */}
            </tr>
          </thead>
          <tbody>
            {filteredClients.map(renderClientRow)}
          </tbody>
        </table>
        {selectedClient && (
          <div className="modal-wrapper" onClick={handleCloseClientDetails}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <span className="close" onClick={handleCloseClientDetails}>&times;</span>
              <ClientDetails client={selectedClient} />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ClientList;
