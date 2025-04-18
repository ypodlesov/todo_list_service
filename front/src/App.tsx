// src/App.tsx
import React, { useEffect, useState } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';

import LoginPage from './LoginPage.tsx';
import TasksPage from './TasksPage.tsx';
import SignUpPage from './SignUpPage.tsx'; // Если хотите отдельную регистрацию
import { checkAuth } from './api.ts';

/**
 * Вариант с двумя состояниями:
 * - loadingAuth: мы ещё не знаем, авторизован ли юзер (идёт запрос /get_tasks или др.)
 * - authorized: как только мы узнаем (true/false), перерисовываем UI
 */
const App: React.FC = () => {
  const [authorized, setAuthorized] = useState(false);
  const [loadingAuth, setLoadingAuth] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const isAuth = await checkAuth();  // /get_tasks → 200 OK или 401 Unauthorized
        setAuthorized(isAuth);
      } catch (err) {
        // Если упадёт — точно не авторизован
        setAuthorized(false);
      } finally {
        setLoadingAuth(false); // закончили проверять
      }
    })();
  }, []);

  // Пока не знаем, авторизован ли пользователь — показываем "загрузку"
  if (loadingAuth) {
    return <div>Loading...</div>;
  }

  // Как только loadingAuth = false, рендерим маршруты
  return (
    <Routes>
      {/* 
        "/" — корневой путь.
        Если авторизован → /tasks
        Иначе → /login
      */}
      <Route
        path="/"
        element={
          authorized ? <Navigate to="/tasks" /> : <Navigate to="/login" />
        }
      />

      {/*
        "/login" 
        Если уже авторизован -> редирект на /tasks
        Иначе показываем LoginPage, куда пробрасываем onLoginSuccess
      */}
      <Route
        path="/login"
        element={
          authorized
            ? <Navigate to="/tasks" />
            : <LoginPage onLoginSuccess={() => setAuthorized(true)} />
        }
      />

      {/*
        "/signup" — отдельная страница регистрации (если нужна).
        Если уже авторизован -> редирект на /tasks
        Иначе показываем SignUpPage.
      */}
      <Route
        path="/signup"
        element={
          authorized
            ? <Navigate to="/tasks" />
            : <SignUpPage onSignUpSuccess={() => setAuthorized(true)} />
        }
      />

      {/* 
        "/tasks"
        Если не авторизован -> /login
        Иначе -> TasksPage
      */}
      <Route
        path="/tasks"
        element={
          authorized
            ? <TasksPage onLogout={() => setAuthorized(false)} />
            : <Navigate to="/login" />
        }
      />

      {/* 
        Любой другой путь -> на "/"
      */}
      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  );
};

export default App;
