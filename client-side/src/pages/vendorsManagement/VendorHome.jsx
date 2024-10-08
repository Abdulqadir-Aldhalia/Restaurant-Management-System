import { useEffect, useState } from "react";
import { Spinner } from "react-bootstrap";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import { message } from "antd";
import "bootstrap/dist/css/bootstrap.min.css";
import { baseUrl } from "../../const";
import "./VendorHome.css";

const VendorAppHome = () => {
  const [vendors, setVendors] = useState(null);
  const [loading, setLoading] = useState(true);
  const userToken = useSelector((state) => state.user.userToken);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const api = axios.create({
    baseURL: baseUrl,
    headers: {
      Authorization: `Bearer ${userToken}`,
    },
  });

  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response && error.response.status === 401) {
        message.error("Session expired. Please login again.");
        dispatch({ type: "LOGOUT" });
        navigate("/login");
      } else if (error.response.status === 403) {
        message.error("Unauthorized to perform this action.");
      }
      return Promise.reject(error);
    },
  );

  useEffect(() => {
    const fetchVendors = async () => {
      try {
        const response = await api.get("/me/vendors");
        const vendorsData = response.data?.data || response.data;
        setVendors(vendorsData);
      } catch (error) {
        message.error("Failed to load vendors");
      } finally {
        setLoading(false);
      }
    };

    fetchVendors();
  }, [userToken]);

  if (loading) {
    return (
      <div
        className="d-flex justify-content-center align-items-center"
        style={{ height: "100vh" }}
      >
        <Spinner animation="border" />
      </div>
    );
  }

  if (!vendors || vendors.length === 0) {
    return <div className="text-center">No vendors available</div>;
  }

  return (
    <>
      <h1>Manage Your Vendors</h1>
      <div className="containers">
        <div className="row">
          {vendors.map((vendor) => (
            <div key={vendor.id} className="card-container">
              <div className="card">
                <img
                  className="card-img"
                  src={`${baseUrl}/${vendor.img}`}
                  alt={vendor.name}
                />
                <div className="card-body">
                  <h5 className="card-title">{vendor.name}</h5>
                  <p className="card-text">
                    {vendor.description || "No description available."}
                  </p>
                  <button
                    className="btn"
                    onClick={() => navigate(`/vendorManagement/${vendor.id}`)}
                  >
                    Manage Vendor
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  );
};

export default VendorAppHome;
