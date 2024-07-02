import React from 'react';

const Products = ({ content, backgroundColor, width, height, fontSize, onClick }) => {
  return (
    <section className="products" style={{ backgroundColor, width, height, fontSize }} onClick={onClick}>
      <h2>{content}</h2>
    </section>
  );
};

export default Products;
