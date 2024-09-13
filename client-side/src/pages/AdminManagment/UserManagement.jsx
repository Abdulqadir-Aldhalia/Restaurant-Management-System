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

  const loadInitialUsers = async () => {
    try {
      const response = await api.get("/users");
      setUsers(response.data);
      setLoading(false);
    } catch (error) {
      message.error("Failed to load users");
      setLoading(false);
    }
  };

  useEffect(() => {
    loadInitialUsers();
  }, [api]);

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
      setUsers([...users, response.data]);
      message.success("User added successfully");
      setIsModalVisible(false);
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
    } catch (error) {
      if (error.response && error.response.status === 409) {
        message.error("Email already exists, please use a different email");
      } else {
        message.error("Failed to add user");
      }
    } finally {
      setLoadingAction(false);
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
      const response = await api.put(`/users/${editingUser.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      // Get the updated user data from the response
      const updatedUser = response.data;

      // Update the users list locally with the new image or name/email/phone
      setUsers(
        users.map((user) =>
          user.id === editingUser.id
            ? {
                ...user,
                name: updatedUser.name,
                email: updatedUser.email,
                phone: updatedUser.phone,
                // Set the updated image URL returned from the server
                img: updatedUser.img,
              }
            : user,
        ),
      );

      message.success("User updated successfully");

      // Reset modal state
      setIsModalVisible(false);
      setEditingUser(null);
      form.resetFields();
      setImageUrl(""); // Reset image URL
      setImageFile(null); // Reset image file
    } catch (error) {
      if (error.response && error.response.status === 409) {
        message.error("Email already exists, please use a different email");
      } else {
        message.error("Failed to update user");
      }
    } finally {
      setLoadingAction(false);
    }
  };

  const handleEdit = (user) => {
    setEditingUser(user);
    setImageUrl(`${baseUrl}${user.img}`);
    form.setFieldsValue({
      name: user.name,
      email: user.email,
      phone: user.phone,
      pssword: user.password,
    });
    setIsModalVisible(true);
    setIsAddingUser(false);
  };

  const handleDeleteUser = async (userId) => {
    try {
      await api.delete(`/users/${userId}`);
      setUsers(users.filter((user) => user.id !== userId));
      message.success("User deleted successfully");
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

          <Button type="primary" htmlType="submit" loading={loadingAction}>
            {isAddingUser ? "Add User" : "Save Changes"}
          </Button>
        </Form>
      </Modal>
    </div>
  );
}

export default UserManagement;
