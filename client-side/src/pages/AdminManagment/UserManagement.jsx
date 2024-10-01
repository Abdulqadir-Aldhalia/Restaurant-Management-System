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
        navigate("/login");
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
      const response = await api.post("/signup", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      message.success("User added successfully");
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
      loadInitialUsers(pagination.current); // Reload users after adding
    } catch (error) {
      if (error.response && error.response.status === 409) {
        message.error("Email already exists, please use a different email");
      } else {
        message.error("Failed to add user");
      }
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false);
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

      // Append the image file if a new one is selected
      if (imageFile) {
        formData.append("img", imageFile);
      }

      // Update the user data in the backend
      await api.put(`/users/${editingUser.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      message.success("User updated successfully");
      loadInitialUsers(pagination.current); // Reload users after editing
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
      setImageUrl(""); // Reset image URL
      setImageFile(null); // Reset image file
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
      message.success("User deleted successfully");
      loadInitialUsers(pagination.current); // Reload users after deletion
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
      {loading ? (
        <Spin tip="Loading..." />
      ) : (
        <div style={{ display: "flex", flexWrap: "wrap", gap: "16px" }}>
          {users.map((user) => {
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
                <p>Phone: {user.phone}</p>
              </Card>
            );
          })}
        </div>
      )}

      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={showAddUserModal}
        style={{ marginTop: "20px" }}
      >
        Add User
      </Button>

      <Pagination
        current={pagination.current}
        pageSize={pagination.pageSize}
        total={pagination.total}
        onChange={handlePageChange}
        style={{ marginTop: "20px", textAlign: "center" }}
      />

      <Modal
        title={editingUser ? "Edit User" : "Add User"}
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form
          form={form}
          onFinish={isAddingUser ? handleAddUser : handleSaveUser}
          layout="vertical"
        >
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: "Please enter name" }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="email"
            label="Email"
            rules={[{ required: true, message: "Please enter email" }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="phone"
            label="Phone"
            rules={[{ required: true, message: "Please enter phone" }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="password"
            label="Password"
            rules={[
              { required: isAddingUser, message: "Please enter password" },
            ]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item label="Upload Image">
            <Upload
              name="image"
              showUploadList={false}
              beforeUpload={() => false}
              onChange={handleImageUpload}
            >
              <Button icon={<PlusOutlined />}>Upload Image</Button>
            </Upload>
            {imageUrl ? (
              <img src={imageUrl} alt="user" className="profile-image" />
            ) : (
              <p>No image uploaded</p>
            )}
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loadingAction}>
              Submit
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

export default UserManagement;
