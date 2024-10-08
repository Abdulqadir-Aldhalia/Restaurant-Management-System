import React, { useEffect, useState, useMemo } from "react";
import { Card, Button, Modal, Input, Upload, message, Form, Spin } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import axios from "axios";
import { baseUrl } from "../../const"; // Make sure to adjust the path to your constants file
import { useSelector } from "react-redux";
import "./VendorProfileManagement.css";
import { useOutletContext } from "react-router";

const VendorAdminProfile = () => {
  const { vendorId } = useOutletContext();
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [vendorData, setVendorData] = useState(null);
  const [imageFile, setImageFile] = useState(null);
  const [imageUrl, setImageUrl] = useState("");
  const [form] = Form.useForm();
  const userToken = useSelector((state) => state.user.userToken); // Adjust according to your state management

  const api = useMemo(() => {
    return axios.create({
      baseURL: baseUrl,
      headers: {
        Authorization: `Bearer ${userToken}`,
      },
    });
  }, [userToken]);

  // Fetch user data on component mount
  useEffect(() => {
    const fetchUserData = async () => {
      setLoading(true);
      try {
        const response = await api.get(`/vendors/${vendorId}`);
        setVendorData(response.data);
        form.setFieldsValue({
          name: response.data.name,
          description: response.data.description,
        });
        if (response.data.img) {
          setImageUrl(`${baseUrl}${response.data.img}`);
        }
      } catch (error) {
        message.error("Failed to load user data");
      } finally {
        setLoading(false);
      }
    };
    fetchUserData();
  }, [api, form]);

  const handleImageUpload = ({ file }) => {
    setImageUrl(URL.createObjectURL(file));
    setImageFile(file);
  };

  const handleUpdateProfile = async (values) => {
    setLoadingAction(true);
    try {
      const formData = new FormData();
      formData.append("name", values.name);
      formData.append("description", values.description);

      if (imageFile) {
        formData.append("img", imageFile);
      }

      await api.put(`/vendors/${vendorData.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      message.success("Profile updated successfully");

      // Refresh user data after update
      const response = await api.get(`/vendors/${vendorId}`);
      setVendorData(response.data);
      form.resetFields(); // Reset the form after updating
    } catch (error) {
      // Check if the error response exists and has a message
      if (
        error.response &&
        error.response.data &&
        error.response.data.message
      ) {
        message.error(error.response.data.message); // Show the specific error message from backend
      } else {
        message.error("Failed to update profile"); // Fallback message
      }
    } finally {
      setLoadingAction(false);
    }
  };

  return (
    <div>
      {loading ? (
        <Spin tip="Loading..." />
      ) : (
        <Card
          title="Vendor Admin Profile"
          style={{ width: 400, margin: "0 auto" }}
        >
          <Form form={form} onFinish={handleUpdateProfile} layout="vertical">
            <Form.Item
              name="name"
              label="Name"
              rules={[
                { required: true, message: "Please enter your Business Name" },
              ]}
            >
              <Input />
            </Form.Item>

            <Form.Item
              name="description"
              label="Description"
              rules={[
                { required: false, message: "Please enter your descripton" },
              ]}
            >
              <Input />
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
                <img src={imageUrl} alt="vendor" className="profile-image" />
              ) : (
                <p>No image uploaded</p>
              )}
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loadingAction}>
                Update Profile
              </Button>
            </Form.Item>
          </Form>
        </Card>
      )}
    </div>
  );
};

export default VendorAdminProfile;
