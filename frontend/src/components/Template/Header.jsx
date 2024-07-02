import React from 'react';

const Header = ({ content, backgroundColor, width, height, fontSize, onClick }) => {
  return (
    <header
      className="header"
      style={{ backgroundColor, width, height, fontSize }}
      onClick={onClick}
    >
      {content}
    </header>
  );
};

export default Header;
