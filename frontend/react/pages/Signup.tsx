import React, { FormEvent, FormEventHandler } from 'react';
import "./Login"

interface SignupElements extends HTMLFormControlsCollection {
    username: HTMLInputElement;
    password: HTMLInputElement;
    confirmPassword: HTMLInputElement;
}

interface SignupForm extends HTMLFormElement {
  readonly elements: SignupElements;
}

export default function Login() {
  
  const handleSubmit = (event: FormEvent<SignupForm>) => {
    event.preventDefault();
    const elements = event.currentTarget.elements;
    const signup = {
      Uname: elements.username.value,
      Pass: elements.password.value,
      confrimPass: elements.confirmPassword.value
    };
    if (signup.Pass != signup.confrimPass) {
        console.log("passwords are not the same!")
    } else {
        console.log(signup); //TODO
    }
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
        <label>Confirm Password</label>
        <input
            id="confirmPassword"
            type="text"
            required
        />
        <button type="submit">Submit</button>
    </form>
    <a href="/login">Login</a>
    </div>
  );
}
