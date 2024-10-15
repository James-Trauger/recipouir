import React from "react";
import ReactDOM from "react-dom/client";
import {BrowserRouter, Routes, Route} from 'react-router-dom';
import About from "./pages/About"
import Home from "./pages/Home"
import Topnav from "./pages/Topnav";

const root = ReactDOM.createRoot(document.querySelector("#application")!);
root.render(
    <BrowserRouter>
        <Topnav></Topnav>
        <Routes> 
            <Route index element={<Home />} />
            <Route path="/about" element={<About />} />
        </Routes>
    </BrowserRouter>
);
