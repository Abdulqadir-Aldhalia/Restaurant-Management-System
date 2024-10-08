import React, { useEffect, useState, useMemo } from "react";
import {
  Card,
  Button,
  Modal,
  Input,
  Upload,
  message,
  Form,
  Popconfirm,
  Spin,
  Pagination,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import axios from "axios";
import { baseUrl } from "../../const";
import { useSelector, useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import "./userManagement.css";

function UserManagement() {
  const [users, setUsers] = useState([]);
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [imageUrl, setImageUrl] = useState("");
  const [imageFile, setImageFile] = useState(null);
  const [form] = Form.useForm();
  const [isAddingUser, setIsAddingUser] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [searchQuery, setSearchQuery] = useState("");

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

  const loadInitialUsers = async (page = 1) => {
    setLoading(true);
    try {
      const response = await api.get(
        `/users?page=${page}&per_page=${pagination.pageSize}`,
      );
      setUsers(response.data.data);
      setFilteredUsers(response.data.data);
      setPagination((prev) => ({
        ...prev,
        current: page,
        total: response.data.meta.total_rows,
      }));
    } catch (error) {
      message.error("Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadInitialUsers(pagination.current);
  }, [api, pagination.current]);

  const handleSearch = async (query) => {
    if (!query) {
      setSearchQuery("");
      loadInitialUsers(pagination.current); // Reload all users if no query
      return;
    }
    setLoading(true);
    try {
      const response = await api.get(`/users?query=${query}`);
      setFilteredUsers(response.data.data || []);
      setPagination((prev) => ({
        ...prev,
        total: response.data.meta.total_rows,
        current: 1,
      })); // Reset pagination
    } catch (error) {
      console.error("Error searching for users:", error);
      alert("An error occurred while searching for users. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const debounce = (func, delay) => {
    let timeout;
    return (...args) => {
      clearTimeout(timeout);
      timeout = setTimeout(() => func.apply(this, args), delay);
    };
  };

  useEffect(() => {
    const debouncedSearch = debounce(handleSearch, 300);
    debouncedSearch(searchQuery);
  }, [searchQuery]);

  const handleAddUser = async (values) => {
    setLoadingAction(true);
    try {
      const formData = new FormData();
      formData.append("name", values.name);
      formData.append("email", values.email);
      formData.append("phone", values.phone);
      formData.append("password", values.password);
      if (imageFile) {
        formData.append("img", imageFile);
      }
      await api.post("/signup", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      message.success("User added successfully");
      form.resetFields();
      loadInitialUsers(pagination.current);
    } catch (error) {
      if (error.response && error.response.status === 409) {
        message.error("Email already exists, please use a different email");
      } else if (error.response && error.response.status == 400) {
        message.error(error.response.data.error);
      }
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false);
      setImageUrl("");
      setImageFile(null);
    }
  };

  const handleSaveUser = async (values) => {
    setLoadingAction(true);
    try {
      const formData = new FormData();
      formData.append("name", values.name);
      formData.append("email", values.email);
      formData.append("phone", values.phone);
      formData.append("password", values.password);
      if (imageFile) {
        formData.append("img", imageFile);
      }
      await api.put(`/users/${editingUser.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      message.success("User updated successfully");
      loadInitialUsers(pagination.current);
    } catch (error) {
      if (error.response && error.response.status === 409) {
        message.error("Email already exists, please use a different email");
      } else {
        message.error("Failed to update user");
      }
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false);
      setEditingUser(null);
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
    }
  };

  const handleEdit = (user) => {
    setEditingUser(user);
    setImageUrl(`${baseUrl}${user.img}`);
    form.setFieldsValue({
      name: user.name,
      email: user.email,
      phone: user.phone,
      password: user.password,
    });
    setIsModalVisible(true);
    setIsAddingUser(false);
  };

  const handleDeleteUser = async (userId) => {
    try {
      await api.delete(`/users/${userId}`);
      // setUsers((prev) => prev.filter((user) => user.id !== userId));
      message.success("User deleted successfully");
      loadInitialUsers(pagination.current);
    } catch (error) {
      message.error("Failed to delete user");
    }
  };

  const handleImageUpload = ({ file }) => {
    setImageUrl(URL.createObjectURL(file));
    setImageFile(file);
  };

  const showAddUserModal = () => {
    form.resetFields();
    setImageUrl("");
    setImageFile(null);
    setEditingUser(null);
    setIsModalVisible(true);
    setIsAddingUser(true);
  };

  const handlePageChange = (page) => {
    setPagination((prev) => ({ ...prev, current: page }));
  };

  return (
    <div>
      <Input
        placeholder="Search by name or email"
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        style={{ margin: "20px" }}
      />
      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={showAddUserModal}
        style={{ margin: "20px" }}
      >
        Add User
      </Button>
      {loading ? (
        <Spin tip="Loading..." />
      ) : (
        <div style={{ display: "flex", flexWrap: "wrap", gap: "16px" }}>
          {filteredUsers.map((user) => {
            const userImageUrl = `${baseUrl}${user.img}`;
            return (
              <Card
                key={user.id}
                hoverable
                style={{ width: 240, textAlign: "center" }}
                cover={
                  <img
                    alt="user"
                    src={userImageUrl}
                    className="profile-image"
                    onError={() =>
                      console.log("Image load failed:", userImageUrl)
                    }
                  />
                }
                actions={[
                  <Button onClick={() => handleEdit(user)}>Edit</Button>,
                  <Popconfirm
                    title="Are you sure you want to delete this user?"
                    onConfirm={() => handleDeleteUser(user.id)}
                    okText="Yes"
                    cancelText="No"
                  >
                    <Button danger>Delete</Button>
                  </Popconfirm>,
                ]}
              >
                <Card.Meta title={user.name} description={user.email} />
              </Card>
            );
          })}
        </div>
      )}
      <Pagination
        current={pagination.current}
        pageSize={pagination.pageSize}
        total={pagination.total}
        onChange={handlePageChange}
        style={{ marginTop: 20 }}
      />
      <Modal
        title={isAddingUser ? "Add User" : "Edit User"}
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form
          form={form}
          onFinish={isAddingUser ? handleAddUser : handleSaveUser}
        >
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item
            name="email"
            label="Email"
            rules={[{ required: true, type: "email" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="Phone" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item
            name="password"
            label="Password"
            rules={[{ required: isAddingUser }]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item label="Profile Image">
            <Upload
              beforeUpload={() => false} // Prevent automatic upload
              showUploadList={false}
              onChange={handleImageUpload}
            >
              <Button>Upload Image</Button>
            </Upload>
            {imageUrl && (
              <img
                src={imageUrl}
                alt="avatar"
                style={{
                  width: "100px",
                  marginTop: "10px",
                  marginLeft: "10px",
                }}
              />
            )}
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={loadingAction} htmlType="submit">
              {isAddingUser ? "Add" : "Save"}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

export default UserManagement;
