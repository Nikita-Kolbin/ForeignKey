// src/components/Header.js
import React from 'react';
import { Link } from 'react-router-dom';
import '../styles/Header.css';

const Header = () => {
    const handleLoginClick = () => {
        window.location.href = '/login';
    };

    return (
        <header className="header">
            <div className="logo">
                <Link to="/">MAESTRO</Link>
            </div>
            <div className="login-container-button">
                <button className="login-button" onClick={handleLoginClick}>личный кабинет</button>
            </div>
        </header>
    );
};

export default Header;
