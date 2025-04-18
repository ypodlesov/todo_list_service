// src/LoginPage.tsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signIn } from './api.ts';

interface LoginPageProps {
  onLoginSuccess: () => void;
}

const LoginPage: React.FC<LoginPageProps> = ({ onLoginSuccess }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);

  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      setError(null);
      // Отправляем запрос POST /sign_in
      await signIn(username, password);
      // Сервер поставил cookie, теперь говорим App-у: мы авторизовались
      onLoginSuccess();
      // А для красоты перейдём на /tasks
      navigate('/tasks');
    } catch (err: any) {
      setError(err.message || 'Sign in failed');
    }
  };

  return (
    <div style={{ maxWidth: 300, margin: '50px auto' }}>
      <h2>Login</h2>
      {error && <div style={{ color: 'red' }}>{error}</div>}

      <div style={{ marginTop: 10 }}>
        <label>Username</label>
        <input
          type="text"
          style={{ width: '100%' }}
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>
      <div style={{ marginTop: 10 }}>
        <label>Password</label>
        <input
          type="password"
          style={{ width: '100%' }}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>

      <button style={{ marginTop: 20 }} onClick={handleLogin}>
        Sign In
      </button>

      <div style={{ marginTop: 10 }}>
        or <a onClick={() => navigate('/signup')}>Sign up</a>
      </div>
    </div>
  );
};

export default LoginPage;
