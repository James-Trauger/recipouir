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
    const creds = {
      Uname: elements.username.value,
      Pass: elements.password.value,
    };
    console.log(creds); //TODO
    // send post request
    // receive jwt token if login is successful
  };
  
  return (
    <div className="login">
      <form onSubmit={handleSubmit}>
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
      <a href="/signup">Signup</a>
    </div>
  );
}
