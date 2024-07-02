import React, { useState, useEffect } from 'react';
import NavigationControlPanel from '../components/NavigationControlPanel';
import EditProfileForm from '../components/PersonalDashboardComponent/EditProfileForm';
import MessageSender from '../components/PersonalDashboardComponent/MessageSender';
import Chart from '../components/PersonalDashboardComponent/Chart';
import { API_BASE_URL } from '../components/ApiConfig';
import '../styles/PersonalDashboard.css';

const PersonalDashboard = () => {
  const [userData, setUserData] = useState({
    first_name: 'John',
    last_name: 'Doe',
    father_name: 'Smith',
    city: 'Example City',
    image_id: null,
    photo: 'https://via.placeholder.com/150',
  });
  const [editMode, setEditMode] = useState(false);

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/admin/get-profile`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });

        const result = await response.json();
        if (result.status === 'OK' && result.profile) {
          setUserData(result.profile);
          if (result.profile.image_id) {
            const imageResponse = await fetch(`${API_BASE_URL}/image/download/${result.profile.image_id}`, {
              method: 'GET',
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              }
            });
            if (imageResponse.ok) {
              const imageBlob = await imageResponse.blob();
              const imageUrl = URL.createObjectURL(imageBlob);
              setUserData(prevData => ({ ...prevData, photo: imageUrl }));
            } else {
              console.error('Ошибка загрузки изображения');
            }
          }
        } else {
          console.error('Ошибка получения профиля', result.error);
        }
      } catch (error) {
        console.error('Ошибка получения профиля', error);
      }
    };

    fetchUserProfile();
  }, []);

  const handleSaveChanges = (updatedUserData) => {
    setUserData(updatedUserData);
    setEditMode(false);
  };

  return (
      <div>
        <NavigationControlPanel />
        <div className="personal-dashboard">
          <div className="top-section">
            <div className="dashboard-section personal-info">
              {editMode ? (
                <EditProfileForm userData={userData} onSaveChanges={handleSaveChanges} />
              ) : (
                <div className="user-info">
                  <div>
                    <img src={userData.photo} alt="User" className="profile-photo" />
                  </div>
                  <div className="info-section">
                    <div className="profile-header">
                      <h3>Персональная информация</h3>
                      <button onClick={() => setEditMode(true)} className="edit-button">
                        <img src="/public/Изменить.png" alt="Edit" className="edit-icon" />
                      </button>
                    </div>
                    <div className="personal-info">
                      <div>
                        <p>Имя: {userData.first_name}</p>
                        <p>Фамилия: {userData.last_name}</p>
                      </div>
                      <div>
                        <p>Отчество: {userData.father_name}</p>
                        <p>Город: {userData.city}</p>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
            <div className="dashboard-section message-sender">
              <MessageSender />
            </div>
          </div>
          <div className="dashboard-section chart-section">
            <Chart />
          </div>
        </div>
      </div>
  );
};

export default PersonalDashboard;
