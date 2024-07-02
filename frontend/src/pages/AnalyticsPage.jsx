// AnalyticsPage.js

import React from 'react';
import NavigationControlPanel from '../components/NavigationControlPanel'; // Импортируем шапку навигации
import RevenueChart from '../components/Analitics/RevenueChart';
import OrdersChart from '../components/Analitics/OrdersChart';
import RevenueHistogram from '../components/Analitics/RevenueHistogram';
import "../styles/AnalyticsPage.css";

const AnalyticsPage = () => {
  return (
    <div>
      <NavigationControlPanel /> {/* Вставляем шапку навигации */}
      <div className="analytics-page">
        <h2>Аналитика</h2>
        <div className="charts-container">
          <div className="chart">
            <h3>График выручки</h3>
            <RevenueChart />
          </div>
          <div className="chart">
            <h3>График заказов</h3>
            <OrdersChart />
          </div>
          <div className="chart">
            <h3>Гистограмма динамики выручки</h3>
            <RevenueHistogram />
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsPage;
