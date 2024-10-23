import React from "react"
import {NavLink} from "react-router-dom"

export default function Topnav() {
  return (
    <nav className="topnav">
        <ul>
            <li className="active"><NavLink to="/">Home</NavLink></li>
            <li><NavLink to="/about">About</NavLink></li>
        </ul>
    </nav>
  )
}
