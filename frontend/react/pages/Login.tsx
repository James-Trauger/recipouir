import React from 'react';

export default function Login() {
  return (
    <div className="login">
        <h2>Login</h2>
        <label>Username:</label>
        <input
            type="text"
            required
        />
        <label>Password</label>
        <input
            type="text"
            required
        />
    </div>
  )
}
