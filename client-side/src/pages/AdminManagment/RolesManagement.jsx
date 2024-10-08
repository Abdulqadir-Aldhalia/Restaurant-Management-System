import React, { useEffect, useState, useMemo } from "react";
import { Input, Button, Table, Select, notification, message } from "antd";
import axios from "axios";
import { baseUrl } from "../../const";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router";

const { Option } = Select;

const RolesManagement = () => {
  const [users, setUsers] = useState([]);
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [loading, setLoading] = useState(false);
  const [selectedRoles, setSelectedRoles] = useState({});
  const [totalUsers, setTotalUsers] = useState(0); // Total users state

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

  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response && error.response.status === 401) {
        message.error("Session expired. Please login again.");
        dispatch({ type: "LOGOUT" });
        navigate("/loginAdminPortal");
      }
      return Promise.reject(error);
    },
  );

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const response = await api.get(`/users`, {
        params: { query, page, per_page: perPage },
      });
      setUsers(response.data.data);
      // Set initial selected roles
      const initialRoles = {};
      response.data.data.forEach((user) => {
        initialRoles[user.id] = ""; // Initialize with no role selected
      });
      setSelectedRoles(initialRoles);

      // Store meta information for pagination
      const { meta } = response.data;
      setTotalUsers(meta.total_rows); // Store the total number of users
    } catch (error) {
      console.error("Failed to fetch users:", error);
      notification.error({
        message: "Error",
        description: "Failed to fetch users. Please try again later.",
      });
    }
    setLoading(false);
  };

  const handleRoleChange = async (userId, roleId, action) => {
    const url =
      action === "grant"
        ? `${baseUrl}/users/grant-role`
        : `${baseUrl}/users/revoke-role`;

    try {
      var formData = new FormData();
      formData.append("user_id", userId);
      formData.append("role_id", roleId);
      const response = await api.post(url, formData);

      if (response.status === 202) {
        notification.success({
          message: `Role successfully ${action === "grant" ? "assigned" : "revoked"}`,
        });
        fetchUsers(); // Refresh users list after action
      }
    } catch (error) {
      if (error.response && error.response.status === 409) {
        notification.error({
          message: "Conflict",
          description: "User already has this role.",
        });
      } else {
        notification.error({
          message: "Error",
          description: "Something went wrong. Please try again.",
        });
      }
    }
  };

  const handleSearch = () => {
    setPage(1); // Reset to the first page when searching
    fetchUsers();
  };

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Email",
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Actions",
      key: "actions",
      render: (text, record) => (
        <>
          <Select
            style={{ width: 120 }}
            value={selectedRoles[record.id] || ""}
            onChange={(value) =>
              setSelectedRoles({ ...selectedRoles, [record.id]: value })
            }
            placeholder="Select Role"
          >
            <Option value="1">Admin</Option>
            <Option value="2">Vendor</Option>
            <Option value="3">Customer</Option>
          </Select>
          <Button
            style={{ marginLeft: 8 }}
            type="primary"
            onClick={() =>
              handleRoleChange(record.id, selectedRoles[record.id], "grant")
            }
            disabled={!selectedRoles[record.id]}
          >
            Assign Role
          </Button>
          <Button
            style={{ marginLeft: 8 }}
            type="danger"
            onClick={() =>
              handleRoleChange(record.id, selectedRoles[record.id], "revoke")
            }
            disabled={!selectedRoles[record.id]}
          >
            Revoke Role
          </Button>
        </>
      ),
    },
  ];

  useEffect(() => {
    fetchUsers();
  }, [page, perPage]);

  return (
    <div>
      <Input.Search
        placeholder="Search by name or email"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onSearch={handleSearch}
        enterButton
        style={{ marginBottom: 16 }}
      />
      <Table
        dataSource={users}
        columns={columns}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: perPage,
          total: totalUsers, // Use the total from the state
          onChange: (page, pageSize) => {
            setPage(page);
            setPerPage(pageSize);
          },
        }}
      />
    </div>
  );
};

export default RolesManagement;
