import React, { useEffect, useState, useMemo } from "react";
import axios from "axios";
import { baseUrl } from "../../const";
import { useNavigate, useOutletContext } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import "./VendorOrderManagement.css";
import { message } from "antd";

const VendorOrderManagement = () => {
  const navigate = useNavigate();
  const userToken = useSelector((state) => state.user.userToken);
  const [loading, setLoading] = useState(true);
  const [updatingOrder, setUpdatingOrder] = useState(null);

  const { vendorId } = useOutletContext();

  const dispatch = useDispatch();
  const [orders, setOrders] = useState([]);

  const api = useMemo(() => {
    return axios.create({
      baseURL: baseUrl,
      headers: {
        Authorization: `Bearer ${userToken}`,
      },
    });
  }, [userToken]);

  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response && error.response.status === 401) {
        message.error("Session expired. Please login again.");
        dispatch({ type: "LOGOUT" });
        navigate("/login");
      }
      return Promise.reject(error);
    },
  );

  const loadInitialOrders = async (vendorId) => {
    setLoading(true);
    try {
      const response = await api.get(`/orders/vendors/${vendorId}`);
      setOrders(response.data);
      message.success("Orders loaded successfully");
    } catch (error) {
      message.error("Failed to load orders");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadInitialOrders(vendorId);
  }, [api]);

  const updateOrderStatus = (orderId, orderNumber, newStatus) => {
    setUpdatingOrder(orderId);
    const updatedOrders = orders.map((order) =>
      order.order_id === orderId ? { ...order, status: newStatus } : order,
    );

    api
      .put(`/orders/${orderId}?vendor_id=${vendorId}&status=${newStatus}`)
      .then((response) => {
        setOrders(updatedOrders);
        message.success(`Order #${orderNumber} status updated to ${newStatus}`);
      })
      .catch((error) => {
        console.error("Error updating order status:", error);
        message.error("Failed to update order status. Please try again.");
      })
      .finally(() => setUpdatingOrder(null)); // reset updating state after completion
  };

  const getStatusColor = (status) => {
    switch (status) {
      case "PENDING":
      case "1":
        return "red"; // Pending
      case "PREPEARING":
      case "2":
        return "orange"; // Preparing
      case "READY":
      case "3":
        return "green"; // Ready
      default:
        return "black"; // Default color
    }
  };

  const getStatus = (status) => {
    switch (status) {
      case "PENDING":
      case "1":
        return "PENDING"; // Pending
      case "PREPEARING":
      case "2":
        return "PREPEARING"; // Preparing
      case "READY":
      case "3":
        return "READY"; // Ready
      default:
        return "UNKNOWN"; // Default value
    }
  };

  const SkeletonLoader = () => (
    <div className="skeleton-loader">
      <div className="skeleton-title"></div>
      <div className="skeleton-line"></div>
      <div className="skeleton-line"></div>
    </div>
  );

  return (
    <>
      <h1>Order Management</h1>
      <div className="containers">
        {loading
          ? Array(3)
              .fill()
              .map((_, i) => <SkeletonLoader key={i} />)
          : orders.map((order) => (
              <div
                className="order"
                key={`${order.order_id}-${order.order_number}`}
              >
                <div className="order-header">
                  <span className="order-id">Order #{order.order_number}</span>

                  {/* Status display */}

                  {/* Status dropdown with placeholder */}
                  <select
                    className="status-select"
                    id={`status-${order.order_id}`}
                    defaultValue=""
                  >
                    <option value="" disabled>
                      Change Status
                    </option>
                    <option value="1">Pending</option>
                    <option value="2">Preparing</option>
                    <option value="3">Ready</option>
                  </select>
                </div>

                <div className="order-details">
                  <p>
                    <strong>Table</strong> {order.table_name}
                  </p>
                  <p>
                    <strong>Total Cost</strong> {order.total_order_cost}
                  </p>
                  <p>
                    <strong>Status</strong>{" "}
                    <span
                      className="status-text"
                      style={{ color: getStatusColor(order.status) }}
                    >
                      {getStatus(order.status)}
                    </span>
                  </p>
                </div>

                <button
                  className="update-btn"
                  onClick={() =>
                    updateOrderStatus(
                      order.order_id,
                      order.order_number,
                      document.getElementById(`status-${order.order_id}`).value,
                    )
                  }
                >
                  Update Status
                </button>
              </div>
            ))}
      </div>
    </>
  );
};

export default VendorOrderManagement;
