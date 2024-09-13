import React from "react";
import { useNavigate } from "react-router";

const Header = () => {
  const navigate = useNavigate();
  return (
    <header className="header">
      <div className="logo">App</div>
      <nav>
        <ul>
          <li>
            <a href="#vendors" className="header-items">
              Vendors
            </a>
          </li>
          <li>
            <a href="#restaurants" className="header-items">
              Restaurants
            </a>
          </li>
          <li>
            <a onClick={() => navigate("/login")} className="login">
              Login
            </a>
          </li>
        </ul>
      </nav>
    </header>
  );
};

export default Header;
