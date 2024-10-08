import React from "react";
import "./UserHome.css"; // Assuming you create a separate CSS file for styling

const UserHome = () => {
  return (
    <div className="user-home-container">
      <h1>
        Hallo! The User Home page is still under development. We might complete
        it in the future.
      </h1>
      <h2>You can test the Admin management and Vendor management for now.</h2>

      <div className="links">
        <a href="/loginAdminPortal" className="portal-link">
          <h3>Admin Portal</h3>
        </a>
        <a href="/loginVendorPortal" className="portal-link">
          <h3>Vendor Portal</h3>
        </a>
      </div>
      <h2 className="note">
        Keep in mind, you need an admin account to log in to the Admin portal.
        The same applies for the Vendor portal, which requires a vendor account!
      </h2>
      <h2 className="note">
        If the user You have created is your first user then it will be the
        admin by default!
      </h2>
    </div>
  );
};

export default UserHome;
