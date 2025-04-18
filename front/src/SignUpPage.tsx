// src/SignUpPage.tsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signUp } from './api.ts';

interface SignUpPageProps {
  onSignUpSuccess: () => void;
}

const SignUpPage: React.FC<SignUpPageProps> = ({ onSignUpSuccess }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const [error, setError] = useState<string | null>(null);
  
  const navigate = useNavigate();

  const handleSignUp = async () => {
    try {
      setError(null);
      await signUp(username, password, email);
      onSignUpSuccess(); // authorized = true
      navigate('/tasks');
    } catch (err: any) {
      setError(err.message || 'Sign up failed');
    }
  };

  return (
    <div style={{ maxWidth: 300, margin: '50px auto' }}>
      <h2>Sign Up</h2>
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
        <label>Email</label>
        <input
          type="email"
          style={{ width: '100%' }}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
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

      <button style={{ marginTop: 20 }} onClick={handleSignUp}>
        Create Account
      </button>
    </div>
  );
};

export default SignUpPage;
