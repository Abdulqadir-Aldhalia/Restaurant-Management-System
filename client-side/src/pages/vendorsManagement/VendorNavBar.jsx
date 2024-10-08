import React, { useEffect } from "react";
import { useDispatch } from "react-redux";

import { PoweroffOutlined } from "@ant-design/icons";
import { Button } from "antd";
import { Link, useNavigate } from "react-router-dom";
import { removeUserToken } from "../../redux/user/userSlice";
import "./VendorNavBar.css";

const VendorNavBar = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const handleLogout = () => {
    localStorage.removeItem("userToken");
    dispatch(removeUserToken());
    navigate("/loginVendorPortal");
  };

  useEffect(() => {
    const navLinks = document.querySelectorAll(".nav-link");

    const handleMouseOver = (event) => {
      event.target.style.transform = "scale(1.1)";
    };

    const handleMouseOut = (event) => {
      event.target.style.transform = "scale(1)";
    };

    const handleClick = (event) => {
      // event.preventDefault();
      const href = event.target.getAttribute("href");
      console.log(`Navigating to: ${href}`);
    };

    navLinks.forEach((link) => {
      link.addEventListener("mouseover", handleMouseOver);
      link.addEventListener("mouseout", handleMouseOut);
      link.addEventListener("click", handleClick);
    });

    return () => {
      navLinks.forEach((link) => {
        link.removeEventListener("mouseover", handleMouseOver);
        link.removeEventListener("mouseout", handleMouseOut);
        link.removeEventListener("click", handleClick);
      });
    };
  }, []);

  return (
    <nav className="navbar">
      <div className="nav-container">
        <h2 className="logo">Vendor Application</h2>
        <div className="nav-links">
          <Link to="VendorHome" className="nav-link">
            Home
          </Link>
          <Link to="VendorAdminProfile" className="nav-link">
            Profile
          </Link>
        </div>
        <button className="custom-button-danger" onClick={handleLogout}>
          <PoweroffOutlined /> Logout
        </button>
      </div>
    </nav>
  );
};

export default VendorNavBar;
