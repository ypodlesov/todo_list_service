// src/index.tsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css'; // Можно убрать или оставить
import { BrowserRouter } from 'react-router-dom';
import App from './App.tsx';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
    <BrowserRouter>
      <App />
    </BrowserRouter>
);
