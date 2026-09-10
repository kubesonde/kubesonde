import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import './styles/global.scss';
import App from './hoc/App';
import reportWebVitals from './reportWebVitals';
import { applyTheme, getInitialTheme } from './utils/theme';

// Set the theme before first paint to avoid a flash of the wrong palette.
applyTheme(getInitialTheme());

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
      <BrowserRouter>
    <App />
      </BrowserRouter>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
