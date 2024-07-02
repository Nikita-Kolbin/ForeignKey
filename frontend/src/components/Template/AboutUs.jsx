import React from 'react';

const AboutUs = ({ content, backgroundColor, width, height, fontSize, onClick }) => {
  return (
    <section className="about-us" style={{ backgroundColor, width, height, fontSize }} onClick={onClick}>
      <h2>{content}</h2>
    </section>
  );
};

export default AboutUs;
