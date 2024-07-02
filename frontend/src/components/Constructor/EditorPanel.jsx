import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './EditorPanel.css';

const EditorPanel = ({ alias, onSave }) => {
  const [backgroundColor, setBackgroundColor] = useState('#ffffff');
  const [font, setFont] = useState('Arial');
  const [availableFonts, setAvailableFonts] = useState([
    'Arial', 'Roboto', 'Open Sans', 'Lato', 'Montserrat'
  ]);
  const navigate = useNavigate();

  useEffect(() => {
    // Dynamically load Google Fonts
    const link = document.createElement('link');
    link.href = 'https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&family=Open+Sans:wght@400;700&family=Lato:wght@400;700&family=Montserrat:wght@400;700&display=swap';
    link.rel = 'stylesheet';
    document.head.appendChild(link);

    return () => {
      document.head.removeChild(link);
    };
  }, []);

  const handleSave = () => {
    const styleData = {
      alias,
      background_color: backgroundColor,
      font,
    };
    onSave(styleData);
  };

  const handleNavigate = () => {
    window.open(`${window.location.origin}/${alias}`, '_blank');
  };

  return (
    <div className="editor-panel">
      <h3>Редактирование стилей сайта</h3>
      <label>
        Цвет фона:
        <input
          type="color"
          value={backgroundColor}
          onChange={(e) => setBackgroundColor(e.target.value)}
        />
      </label>
      <label>
        Шрифт:
        <select value={font} onChange={(e) => setFont(e.target.value)}>
          {availableFonts.map((fontName) => (
            <option key={fontName} value={fontName}>{fontName}</option>
          ))}
        </select>
      </label>
      <div className="editor-panel-buttons">
        <button onClick={handleSave}>Сохранить</button>
        <button onClick={handleNavigate}>Перейти на вебсайт</button>
      </div>
    </div>
  );
};

export default EditorPanel;
