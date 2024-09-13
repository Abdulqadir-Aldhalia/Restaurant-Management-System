import React, { useEffect, useState } from "react";
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
import { Switch } from "antd";

function VendorManagement() {
  const [vendors, setVendors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [editingVendor, setEditingVendor] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [imageUrl, setImageUrl] = useState("");
  const [subscriptionStatus, setSubscriptionStatus] = useState(null);
  const [imageFile, setImageFile] = useState(null);
  const [form] = Form.useForm();
  const [isAddingVendor, setIsAddingVendor] = useState(false);

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
    let isMounted = true;

    const fetchVendors = async () => {
      try {
        const response = await api.get("/vendors", {
          headers: { Authorization: `Bearer ${userToken}` },
        });
        if (isMounted) {
          setVendors(response.data);
          setLoading(false);
        }
      } catch (error) {
        if (isMounted) {
          message.error("Failed to load vendors");
          setLoading(false);
        }
      }
    };

    fetchVendors();

    // Cleanup function
    return () => {
      isMounted = false;
    };
  }, [userToken]);

  const handleAddVendor = async (values) => {
    setLoadingAction(true);
    try {
      const formData = new FormData();
      formData.append("name", values.name);
      formData.append("description", values.description);
      if (imageFile) {
        formData.append("img", imageFile);
      }
      const response = await api.post("/vendors", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      setVendors([...vendors, response.data]);
      message.success("Vendor added successfully");
      setIsModalVisible(false);
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
    } catch (error) {
      message.error("Failed to add vendor");
    } finally {
      setLoadingAction(false);
    }
  };

  const handleSaveVendor = async (values) => {
    setLoadingAction(true);
    try {
      const formData = new FormData();
      formData.append("name", values.name);
      formData.append("description", values.description);

      // Append the image file if a new one is selected
      if (imageFile) {
        formData.append("img", imageFile);
      }

      // Update the vendor data in the backend
      const response = await api.put(`/vendors/${editingVendor.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      // Get the updated vendor data from the response
      const updatedVendor = response.data;

      // Update the vendor list locally with the new image or name/description
      setVendors(
        vendors.map((vendor) =>
          vendor.id === editingVendor.id
            ? {
                ...vendor,
                name: updatedVendor.name,
                description: updatedVendor.description,
                // Set the updated image URL returned from the server
                img: updatedVendor.img,
              }
            : vendor,
        ),
      );

      message.success("Vendor updated successfully");

      // Reset modal state
      setIsModalVisible(false);
      setEditingVendor(null);
      form.resetFields();
      setImageUrl(""); // Reset image URL
      setImageFile(null); // Reset image file
    } catch (error) {
      message.error("Failed to update vendor");
    } finally {
      setLoadingAction(false);
    }
  };
  const handleEdit = (vendor) => {
    setEditingVendor(vendor);
    setImageUrl(`${baseUrl}${vendor.img}`);
    form.setFieldsValue({
      name: vendor.name,
      description: vendor.description,
    });
    setIsModalVisible(true);
    setIsAddingVendor(false);
  };

  const handleDeleteVendor = async (vendorId) => {
    try {
      await api.delete(`/vendors/${vendorId}`);
      setVendors(vendors.filter((vendor) => vendor.id !== vendorId));
      message.success("Vendor deleted successfully");
    } catch (error) {
      message.error("Failed to delete vendor");
    }
  };

  const handleImageUpload = ({ file }) => {
    setImageUrl(URL.createObjectURL(file));
    setImageFile(file);
  };

  const showAddVendorModal = () => {
    form.resetFields();
    setImageUrl("");
    setImageFile(null);
    setEditingVendor(null);
    setIsModalVisible(true);
    setIsAddingVendor(true);
  };

  const handleSubscriptionToggle = async (vendorId, isSubscribed) => {
    setLoadingAction(true);
    try {
      await api.put(`/vendors/${vendorId}/subscription`, {
        status: isSubscribed ? "enable" : "disable",
      });
      // Update the local state to reflect the change
      setVendors(
        vendors.map((vendor) =>
          vendor.id === vendorId ? { ...vendor, isSubscribed } : vendor,
        ),
      );
      message.success(
        `Subscription ${isSubscribed ? "enabled" : "disabled"} successfully`,
      );
    } catch (error) {
      message.error("Failed to update subscription status");
    } finally {
      setLoadingAction(false);
    }
  };

  return (
    <div>
      {loading ? (
        <Spin tip="Loading..." />
      ) : (
        <div style={{ display: "flex", flexWrap: "wrap", gap: "16px" }}>
          {vendors.map((vendor) => {
            const vendorImageUrl = `${baseUrl}${vendor.img}`;
            return (
              <Card
                key={vendor.id}
                hoverable
                style={{ width: 240 }}
                cover={
                  <img
                    alt="vendor"
                    src={vendorImageUrl}
                    onError={() =>
                      console.log("Image load failed:", vendorImageUrl)
                    }
                  />
                }
                actions={[
                  <Button onClick={() => handleEdit(vendor)}>Edit</Button>,
                  <Popconfirm
                    title="Are you sure you want to delete this vendor?"
                    onConfirm={() => handleDeleteVendor(vendor.id)}
                    okText="Yes"
                    cancelText="No"
                  >
                    <Button danger>Delete</Button>
                  </Popconfirm>,
                  <Switch
                    checked={vendor.isSubscribed} // Use the subscription status here
                    onChange={(checked) =>
                      handleSubscriptionToggle(vendor.id, checked)
                    }
                  />,
                ]}
              >
                <Card.Meta
                  title={vendor.name}
                  description={vendor.description}
                />
              </Card>
            );
          })}
        </div>
      )}

      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={showAddVendorModal}
        style={{ marginTop: "20px" }}
      >
        Add Vendor
      </Button>

      <Modal
        title={editingVendor ? "Edit Vendor" : "Add Vendor"}
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form
          form={form}
          onFinish={isAddingVendor ? handleAddVendor : handleSaveVendor}
          layout="vertical"
        >
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: "Please enter vendor name" }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="description"
            label="Description"
            rules={[
              { required: true, message: "Please enter vendor description" },
            ]}
          >
            <Input.TextArea />
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
            {imageUrl && (
              <img
                src={imageUrl}
                alt="vendor"
                style={{ width: 100, marginTop: 10 }}
              />
            )}
          </Form.Item>

          <Form.Item label="Subscription">
            <Switch
              checked={editingVendor ? editingVendor.isSubscribed : false}
              onChange={(checked) =>
                form.setFieldsValue({ subscription: checked })
              }
            />
          </Form.Item>

          <Button type="primary" htmlType="submit" loading={loadingAction}>
            {isAddingVendor ? "Add Vendor" : "Save Changes"}
          </Button>
        </Form>
      </Modal>
    </div>
  );
}

export default VendorManagement;
