import React from "react";
import NavBar from "./VendorManagementNavBar";
import { Outlet, useParams } from "react-router-dom"; // Outlet is used to render child routes

const VendorManagement = () => {
  const { vendorId } = useParams();
  return (
    <>
      <NavBar />
      <div className="content">
        {/* <h2>Managing Vendor ID: {vendorId}</h2> */}
        {/* Outlet will render the child routes defined in the parent route */}
        <Outlet context={{ vendorId }} />{" "}
        {/* Pass vendorId to child components */}
      </div>
    </>
  );
};

export default VendorManagement;
