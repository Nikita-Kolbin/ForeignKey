import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../components/ApiConfig';
import "../styles/OrdersPage.css";
import NavigationControlPanel from '../components/NavigationControlPanel';
import OrderProductsTable from '../components/Orders/OrderProductsTable';
import OrderClientData from '../components/Orders/OrderClientData';

const OrdersPage = () => {
  const [showClosedOrders, setShowClosedOrders] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [orders, setOrders] = useState([]);
  const [completedOrders, setCompletedOrders] = useState([]);
  const [alias, setAlias] = useState('');
  const [error, setError] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(5);

  const orderStatuses = [
    { value: 0, label: 'Ожидает подтверждения' },
    { value: 1, label: 'Взят в работу' },
    { value: 2, label: 'В работе' },
    { value: 3, label: 'Сделан' },
    { value: 4, label: 'Отправлен' },
    { value: 5, label: 'Доставлен' },
    { value: 6, label: 'Завершен' },
    { value: 7, label: 'Нестандартная ситуация' },
  ];

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
    if (!alias) return;

    const fetchOrders = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/order/get-by-alias/${alias}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          const ordersWithClientData = await Promise.all(data.orders.map(async (order) => {
            const clientData = await fetchClientData(order.customer_id);
            return {
              id: order.id,
              client: `${clientData.fullName.split(' ')[0]}.${clientData.fullName.split(' ')[1][0]}.${clientData.fullName.split(' ')[2][0]}`,
              amount: calculateTotalPrice(order.order_items),
              date: formatDate(order.date_time),
              status: order.status,
              comment: order.comment, // Добавляем комментарий к заказу
              products: order.order_items.map(item => ({
                id: item.product.id,
                name: item.product.name,
                price: item.product.price,
                quantity: item.count
              })),
              clientData
            };
          }));
          setOrders(ordersWithClientData);
        } else {
          setError(data.error || 'Ошибка при получении заказов');
        }
      } catch (error) {
        setError('Ошибка при получении заказов');
      }
    };

    const fetchCompletedOrders = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/order/get-completed/${alias}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const data = await response.json();

        if (response.ok && data.status === 'OK') {
          const ordersWithClientData = await Promise.all(data.orders.map(async (order) => {
            const clientData = await fetchClientData(order.customer_id);
            return {
              id: order.id,
              client: `${clientData.fullName.split(' ')[0]}.${clientData.fullName.split(' ')[1][0]}.${clientData.fullName.split(' ')[2][0]}`,
              amount: calculateTotalPrice(order.order_items),
              date: formatDate(order.date_time),
              status: order.status,
              comment: order.comment, // Добавляем комментарий к заказу
              products: order.order_items.map(item => ({
                id: item.product.id,
                name: item.product.name,
                price: item.product.price,
                quantity: item.count
              })),
              clientData
            };
          }));
          setCompletedOrders(ordersWithClientData);
        } else {
          setError(data.error || 'Ошибка при получении завершенных заказов');
        }
      } catch (error) {
        setError('Ошибка при получении завершенных заказов');
      }
    };

    fetchOrders();
    if (showClosedOrders) {
      fetchCompletedOrders();
    }
  }, [alias, showClosedOrders]);

  const fetchClientData = async (customerId) => {
    try {
      const response = await fetch(`${API_BASE_URL}/customer/get-by-alias/${alias}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
  
      const data = await response.json();
  
      if (response.ok && data.status === 'OK') {
        const customer = data.customers.find(customer => customer.id === customerId);
        if (customer) {
          return {
            fullName: `${customer.last_name} ${customer.first_name} ${customer.father_name}`,
            email: customer.email,
            paymentMethod: customer.payment_type,
            delivery: customer.delivery_type,
            phone: customer.phone,
            telegram: customer.telegram,
            comment: customer.comment // Убедитесь, что комментарий включен
          };
        } else {
          console.error('Клиент не найден');
          return {};
        }
      } else {
        console.error('Ошибка при получении данных клиента:', data.error);
        return {};
      }
    } catch (error) {
      console.error('Ошибка при получении данных клиента:', error);
      return {};
    }
  };
  

  const handleOrderStatusChange = async (orderId, newStatus) => {
    try {
      const response = await fetch(`${API_BASE_URL}/order/set-status`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ order_id: orderId, status: newStatus })
      });

      const data = await response.json();

      if (response.ok && data.status === 'OK') {
        setOrders(prevOrders => prevOrders.map(order =>
          order.id === orderId ? { ...order, status: newStatus } : order
        ));
        setCompletedOrders(prevOrders => prevOrders.map(order =>
          order.id === orderId ? { ...order, status: newStatus } : order
        ));
        setSelectedOrder(prevOrder => prevOrder ? { ...prevOrder, status: newStatus } : null);
      } else {
        setError(data.error || 'Ошибка при изменении статуса заказа');
      }
    } catch (error) {
      setError('Ошибка при изменении статуса заказа');
    }
  };

  const toggleShowClosedOrders = () => {
    setShowClosedOrders(!showClosedOrders);
  };

  const filteredOrders = showClosedOrders ? completedOrders : orders.filter(order => order.status !== 6);

  const handleOrderDetailsClick = (order) => {
    setSelectedOrder(order);
  };

  const calculateTotalPrice = (products) => {
    return Math.floor(products.reduce((total, product) => total + parseFloat(product.price) * product.quantity, 0));
  };

  const calculateTotalQuantity = (products) => {
    return products.reduce((total, product) => total + product.quantity, 0);
  };

  const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'numeric', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('ru-RU', options);
  };

  // Pagination logic
  const handlePageChange = (direction) => {
    if (direction === 'prev' && currentPage > 1) {
      setCurrentPage(currentPage - 1);
    } else if (direction === 'next' && currentPage < Math.ceil(filteredOrders.length / itemsPerPage)) {
      setCurrentPage(currentPage + 1);
    }
  };

  const handleItemsPerPageChange = (event) => {
    setItemsPerPage(parseInt(event.target.value));
    setCurrentPage(1); // Reset to first page when items per page changes
  };

  const paginatedOrders = filteredOrders.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

  return (
    <div>
      <NavigationControlPanel />
      <div className="orders-page">
        {error && <p className="error">{error}</p>}
        <table className="orders-table">
          <thead>
            <tr>
              <th>Заказ</th>
              <th>Клиент</th>
              <th>Сумма</th>
              <th>Оформлен</th>
              <th>
                <button onClick={toggleShowClosedOrders}>
                  {showClosedOrders ? 'Активные' : 'Закрытые'}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {paginatedOrders.map(order => (
              <tr key={order.id}>
                <td>{order.id}</td>
                <td>{order.client}</td>
                <td>{calculateTotalPrice(order.products)}</td>
                <td>{order.date}</td>
                <td>
                  <button onClick={() => handleOrderDetailsClick(order)}>Подробнее</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        <div className="pagination">
          <button onClick={() => handlePageChange('prev')} disabled={currentPage === 1}>
            <img src="/Назад.svg" alt="назад" />
          </button>
          <select value={itemsPerPage} onChange={handleItemsPerPageChange}>
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={15}>15</option>
          </select>
          <button onClick={() => handlePageChange('next')} disabled={currentPage === Math.ceil(filteredOrders.length / itemsPerPage)}>
            <img src="/Вперед.svg" alt="Вперед" />
          </button>
        </div>
      </div>
      {selectedOrder && (
        <div className="modal-wrapper" onClick={() => setSelectedOrder(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <span className="close" onClick={() => setSelectedOrder(null)}>&times;</span>
            <h2>Заказ номер: {selectedOrder.id}</h2>
            <OrderProductsTable products={selectedOrder.products} />
            <p>Итого: ${calculateTotalPrice(selectedOrder.products).toFixed(2)}</p>
            <p>Всего товаров: {calculateTotalQuantity(selectedOrder.products)}</p>
            <p>Оформлен: {selectedOrder.date}</p>
            <OrderClientData clientData={selectedOrder.clientData} orderDate={selectedOrder.date} />
            <p>Комментарий: {selectedOrder.comment}</p> {/* Отображение комментария */}
            <div>
              <label htmlFor="order-status">Статус заказа:</label>
              <select
                id="order-status"
                value={selectedOrder.status}
                onChange={(e) => handleOrderStatusChange(selectedOrder.id, parseInt(e.target.value))}
              >
                {orderStatuses.map(status => (
                  <option key={status.value} value={status.value}>
                    {status.label}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrdersPage;
