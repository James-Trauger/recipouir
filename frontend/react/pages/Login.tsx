import React, { FormEvent, FormEventHandler } from 'react';

interface LoginElements extends HTMLFormControlsCollection {
  username: HTMLInputElement;
  password: HTMLInputElement;
}

interface LoginForm extends HTMLFormElement {
  readonly elements: LoginElements;
}

export default function Login() {
  
  const handleSubmit = (event: FormEvent<LoginForm>) => {
    event.preventDefault();
    const elements = event.currentTarget.elements;
    const data = {
      username: elements.username.value,
      password: elements.password.value,
    };
    console.log(data);
    // send post request
  };
  
  return (
    <form onSubmit={handleSubmit} className="login">
        <h2>Login</h2>
        <label>Username:</label>
        <input
            id="username"
            type="text"
            required
        />
        <label>Password</label>
        <input
            id="password"
            type="text"
            required
        />
        <button type="submit">Submit</button>
    </form>
  );
}
