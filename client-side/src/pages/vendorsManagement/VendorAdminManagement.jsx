import React, { useState, useEffect, useMemo } from "react";
import axios from "axios";
import { useNavigate, useOutletContext } from "react-router";
import { baseUrl } from "../../const";
import { useDispatch, useSelector } from "react-redux";
import { Card, Button } from "antd"; // Assuming you're using Ant Design
import "./VendorAdminManagement.css"; // Import your CSS file

const VendorAdminManagement = () => {
  const [admins, setAdmins] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [feedbackMessage, setFeedbackMessage] = useState("");
  const { vendorId } = useOutletContext();
  const userToken = useSelector((state) => state.user.userToken);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const api = useMemo(() => {
    return axios.create({
      baseURL: baseUrl,
      headers: {
        Authorization: `Bearer ${userToken}`,
      },
    });
  }, [userToken]);

  useEffect(() => {
    fetchAdmins();
  }, [api]);

  const fetchAdmins = async () => {
    setLoading(true);
    try {
      const response = await api.get(`/vendors/${vendorId}/admins`);
      console.log("Fetched Admins:", response.data);
      setAdmins(response.data || []);
    } catch (error) {
      console.error("Error fetching admins:", error);
      setFeedbackMessage("Error fetching admins. Please try again later.");
    } finally {
      setLoading(false);
    }
  };

  // Debounce function
  const debounce = (func, delay) => {
    let timeout;
    return (...args) => {
      clearTimeout(timeout);
      timeout = setTimeout(() => func.apply(this, args), delay);
    };
  };

  const handleSearch = async (query) => {
    if (!query) {
      setUsers([]);
      return;
    }
    setLoading(true);
    try {
      const response = await api.get(`/users?query=${query}`);
      console.log("Search Results:", response.data);
      setUsers(response.data.data || []);
    } catch (error) {
      console.error("Error searching for users:", error);
      alert("An error occurred while searching for users. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  // Use effect for live search
  useEffect(() => {
    const debouncedSearch = debounce(handleSearch, 300); // 300ms debounce
    debouncedSearch(searchQuery);
  }, [searchQuery]);

  const addAdmin = async (userId) => {
    try {
      await api.post(
        `/vendors/assign-admin?user_id=${userId}&vendor_id=${vendorId}`,
      );
      fetchAdmins(); // Refresh the admin list
      setFeedbackMessage("Admin added successfully!");
      setTimeout(() => setFeedbackMessage(""), 3000); // Clear feedback after 3 seconds
    } catch (error) {
      console.error("Error adding admin:", error);
      setFeedbackMessage("Error adding admin. Please try again.");
    }
  };

  const revokeAdmin = async (userId) => {
    try {
      await api.post(
        `/vendors/revoke-admin?user_id=${userId}&vendor_id=${vendorId}`,
      );
      fetchAdmins(); // Refresh the admin list
      setFeedbackMessage("Admin revoked successfully!");
      setTimeout(() => setFeedbackMessage(""), 3000); // Clear feedback after 3 seconds
    } catch (error) {
      console.error("Error revoking admin:", error);
      setFeedbackMessage("Error revoking admin. Please try again.");
    }
  };

  return (
    <div className="admin-management-container">
      <h1>Admin Management</h1>
      {feedbackMessage && <p className="feedback-message">{feedbackMessage}</p>}
      <form onSubmit={(e) => e.preventDefault()} className="search-form">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search for users..."
          className="search-input"
        />
      </form>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <ul className="user-list">
          {users.map((user) => (
            <li key={user.id} className="user-item">
              <Card
                hoverable
                style={{ width: 240, textAlign: "center", margin: "10px" }}
                cover={
                  <img
                    alt="user"
                    src={
                      user.img ? `${baseUrl}${user.img}` : "default-avatar.png"
                    }
                    className="profile-image"
                    onError={() => console.log("Image load failed:", user.img)}
                  />
                }
              >
                <Card.Meta title={user.name} description={user.email} />
                <p>Phone: {user.phone}</p>
                <Button onClick={() => addAdmin(user.id)}>Add as Admin</Button>
              </Card>
            </li>
          ))}
        </ul>
      )}

      <h2>Current Admins</h2>
      <div className="admin-cards-container">
        {admins.length > 0 ? (
          admins.map((admin) => (
            <Card
              key={admin.id}
              hoverable
              style={{ width: 240, textAlign: "center", margin: "10px" }}
              cover={
                <img
                  alt="user"
                  src={
                    admin.img ? `${baseUrl}${admin.img}` : "default-avatar.png"
                  }
                  className="profile-image"
                  onError={() => console.log("Image load failed:", admin.img)}
                />
              }
              actions={[
                <Button onClick={() => revokeAdmin(admin.id)}>
                  Revoke Admin
                </Button>,
              ]}
            >
              <Card.Meta title={admin.name} description={admin.email} />
              <p>Phone: {admin.phone}</p>
            </Card>
          ))
        ) : (
          <p>No admins found.</p>
        )}
      </div>
    </div>
  );
};

export default VendorAdminManagement;
