import React from 'react';

const WelcomeScreen = ({ content, backgroundColor, width, height, fontSize, onClick }) => {
  return (
    <section className="welcome-screen" style={{ backgroundColor, width, height, fontSize }} onClick={onClick}>
      <h1>{content}</h1>
    </section>
  );
};

export default WelcomeScreen;
