import React from 'react';

const About = () => {
  return (
    <div>
      <h1>About Page</h1>
      <button onClick={() => window.location.href = '/'}>Go to Home</button>
    </div>
  );
};

export default About;
