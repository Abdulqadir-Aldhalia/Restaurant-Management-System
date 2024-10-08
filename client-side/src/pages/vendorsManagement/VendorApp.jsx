import React from "react";
import { Outlet } from "react-router-dom"; // Outlet is used to render child routes
import VendorNavBar from "./VendorNavBar";

const VendorApp = () => {
  return (
    <>
      <VendorNavBar />
      <div className="content">
        {/* Outlet will render the child routes defined in the parent route */}
        <Outlet />
      </div>
    </>
  );
};

export default VendorApp;
