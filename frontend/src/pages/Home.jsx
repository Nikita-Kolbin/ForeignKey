import React, { useEffect } from 'react';
import Header from '../components/Header';
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import AOS from 'aos';
import 'aos/dist/aos.css';
import '../styles/Home.css';

const Home = () => {
  useEffect(() => {
    AOS.init({
      duration: 600, // Установите длительность анимации в 1000 миллисекунд (1 секунда)
      delay: 200, // Задержка перед началом анимации в 200 миллисекунд
      easing: 'ease-in-out', // Используйте функцию сглаживания ease-in-out
      offset: 200, // Начните анимацию, когда элемент находится в 200 пикселях от нижней границы окна браузера
    });
    const smoothScroll = (event) => {
      event.preventDefault();
      const delta = Math.max(-1, Math.min(1, (event.wheelDelta || -event.detail)));
      document.documentElement.scrollTop -= delta * 30; // Измените значение 30 для регулировки скорости прокрутки
    };

    window.addEventListener('mousewheel', smoothScroll, { passive: false });
    window.addEventListener('DOMMouseScroll', smoothScroll, { passive: false });

    return () => {
      window.removeEventListener('mousewheel', smoothScroll);
      window.removeEventListener('DOMMouseScroll', smoothScroll);
    };
  }, []);

  return (
      
        <div className="home">
          <Header />
          <div className="block block-1" data-aos="fade-up">
            <div className="section_text">
              <h1>Создайте сайт и начните продавать уже сегодня!</h1>
              <p>Превратите свое хобби в прибыльное дело с нашей CRM системой для индивидуальных мастеров</p>
              <button className="cta-button"><a href="/register">Начать сейчас</a></button>
            </div>
            <div className='section-img'>
              <img src="/Home1.png" alt="абоба" />
            </div>
          </div>
          <div className="block block-2" data-aos="fade-up">
            <h2 className='advertising'>Все что нужно для создания интернет магазина!</h2>
            <div className="feature feature-1">
              <h2>Готовые дизайны</h2>
              <p>В Maestro уже готовы несколько красивых вариантов дизайна для сайта. Поэтому вам остается только наполнить своей информацией готовый шаблон и опубликовать сайт.</p>
              <img src='/валик.png'/>
              <button className="cta-button"><a href="/register">Выбрать шаблон</a></button>
            </div>
            <div className="feature feature-2">
              <h2>Работа с заявками</h2>
              <p>Все заявки автоматически попадают в единую панель, где можно менять статусы заказов и вести продажи. Встроенная CRM считает прибыль по каждому клиенту.</p>
            </div>
            <div className="feature feature-3">
              <p>Удобное пространство посетителя сайта для просмотра выбранных товаров и их оплаты.</p>
              <h2>Корзина</h2>
            </div>
            <div className="feature feature-4">
              <p>Ваши клиенты могут оплачивать товары на сайте удобным для них способом.</p>
              <h2>Оплата</h2>
            </div>
          </div>
          <div className="block block-3" data-aos="fade-up">
            <h2>Готовый сайт за три шага</h2>
            <div className="steps">
              <div className="step" data-step="1">
                <h3>Соберите основу</h3>
                <p>Зарегистрируйтесь на платформе Maestro — конструктор сайтов. Придумайте название сайта и его наполнение.</p>
              </div>
              <div className="step" data-step="2">
                <h3>Выберите оформление</h3>
                <p>Структура сайта уже готова. Выберите один из готовых шаблонов и наполните его своими текстами и изображениями. Настройте шрифты и цвета.</p>
              </div>
              <div className="step" data-step="3">
                <h3>Настройте магазин</h3>
                <p>Добавьте описание, фотографии и цены своих товаров и услуг. Подключите прием онлайн-оплаты и отслеживайте заказы в личном кабинете.</p>
              </div>
            </div>
          </div>
          <div className="block block-4" data-aos="fade-up">
            <h2>Попробуйте сейчас</h2>
            <p>Ваш бизнес начинается с сайта — запустите его за 15 минут, без программистов и дизайнеров</p>
            <button className="cta-button"><a href="/register">Создать сайт бесплатно</a></button>
          </div>
          <div className="footer">
            <p>MAESTRO</p>
            <div className="social-links">
              <a href="#"><img src="/vk_icon.png" alt="VK" /></a>
              <a href="#"><img src="/telegram_icon.png" alt="Telegram" /></a>
            </div>
          </div>
        </div>
      
  );
};

export default Home;
