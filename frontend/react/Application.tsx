import React from "react";
import ReactDOM from "react-dom/client";
import {BrowserRouter, Routes, Route} from 'react-router-dom';
import About from "./pages/About";
import Home from "./pages/Home";
import Login from "./pages/Login"
import Signup from "./pages/Signup"
import Topnav from "./components/Topnav";


const root = ReactDOM.createRoot(document.querySelector("#application")!);
root.render(
    <BrowserRouter>
        <Topnav></Topnav>
        <Routes> 
            <Route index element={<Home />} />
            <Route path="/about" element={<About />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
        </Routes>
    </BrowserRouter>
);
