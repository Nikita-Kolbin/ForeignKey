import React, { useState } from 'react';
import Header from './Header'; // импортируем компонент шапки
import WelcomeScreen from './WelcomeScreen'; // импортируем компонент приветственного экрана
import Products from './Products'; // импортируем компонент каталога товаров
import AboutUs from './AboutUs'; // импортируем компонент раздела "О нас"
import Footer from './Footer'; // импортируем компонент футера
import EditableBlock from '../Constructor/EditableBlock';
import EditorPanel from '../Constructor/EditorPanel';

const Template1 = () => {
  const [blocks, setBlocks] = useState([
    { id: 1, content: 'Блок 1', backgroundColor: '#1a1a1a' },
    { id: 2, content: 'Блок 2', backgroundColor: '#1a1a1a' },
  ]);

  const [selectedBlock, setSelectedBlock] = useState(null);

  const handleBlockClick = (block) => {
    setSelectedBlock(block);
  };

  const handleContentChange = (newContent) => {
    setBlocks((prevBlocks) =>
      prevBlocks.map((block) =>
        block.id === selectedBlock.id
          ? { ...block, content: newContent }
          : block
      )
    );
  };

  const handleColorChange = (newColor) => {
    setBlocks((prevBlocks) =>
      prevBlocks.map((block) =>
        block.id === selectedBlock.id
          ? { ...block, backgroundColor: newColor }
          : block
      )
    );
  };

  return (
    <div className="template1">
      <Header /> {/* Вставляем компонент шапки */}
      <WelcomeScreen /> {/* Вставляем компонент приветственного экрана */}
      <Products /> {/* Вставляем компонент каталога товаров */}
      <AboutUs /> {/* Вставляем компонент раздела "О нас" */}
      <Footer /> {/* Вставляем компонент футера */}
      <div className="blocks-container">
        {blocks.map((block) => (
          <EditableBlock
            key={block.id}
            block={block}
            onClick={handleBlockClick}
          />
        ))}
      </div>
      <EditorPanel
        selectedBlock={selectedBlock}
        onChangeContent={handleContentChange}
        onChangeColor={handleColorChange}
      />
    </div>
  );
};

export default Template1;
