import React from 'react';

const Footer = ({ content, backgroundColor, width, height, fontSize, onClick }) => {
  return (
    <footer style={{ backgroundColor, width, height, fontSize }} onClick={onClick}>
      <p>{content}</p>
    </footer>
  );
};

export default Footer;
