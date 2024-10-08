import React, { useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "antd";
import "./VendorManagementNavBar.css";
import { useDispatch } from "react-redux";
import { removeUserToken } from "../../redux/user/userSlice";
import { PoweroffOutlined, RightCircleFilled } from "@ant-design/icons";

const VendorNavBar = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const handleLogout = () => {
    localStorage.removeItem("userToken");
    dispatch(removeUserToken());
    navigate("/loginAdminPortal");
  };

  const goHome = () => {
    navigate("/vendorApp");
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
        <h2 className="logo">Vendor Management</h2>
        <div className="nav-links">
          <Link to="VendorOrderManagement" className="nav-link">
            Orders
          </Link>
          <Link to="VendorItemsManagement" className="nav-link">
            Items
          </Link>
          <Link to="VendorTableManagement" className="nav-link">
            Tables
          </Link>
          <Link to="VendorAdminManagement" className="nav-link">
            Admins
          </Link>
          <Link to="VendorProfileManagement" className="nav-link">
            Profile
          </Link>
        </div>

        {/* "My Vendors" Primary Button */}
        <button className="custom-button-primary" onClick={goHome}>
          <RightCircleFilled /> My Vendors
        </button>

        {/* "Logout" Danger Button */}
        {/* <button className="custom-button-danger" onClick={handleLogout}> */}
        {/*   <PoweroffOutlined /> Logout */}
        {/* </button> */}
      </div>
    </nav>
  );
};

export default VendorNavBar;
