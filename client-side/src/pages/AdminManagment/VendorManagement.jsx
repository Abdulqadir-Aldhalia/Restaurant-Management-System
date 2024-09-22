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
  Select,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import axios from "axios";
import { baseUrl } from "../../const";
import { useSelector, useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { Switch } from "antd";

function VendorManagement() {
  const [vendors, setVendors] = useState([]);
  const [users, setUsers] = useState([]); // List of users to assign as admins
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [editingVendor, setEditingVendor] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [imageUrl, setImageUrl] = useState("");
  const [imageFile, setImageFile] = useState(null);
  const [currentAdmin, setCurrentAdmin] = useState(null); // Current admin of the vendor
  const [form] = Form.useForm();
  const [isAddingVendor, setIsAddingVendor] = useState(false);
  const [selectedAdmin, setSelectedAdmin] = useState(null); // Selected admin for the vendor

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
        const response = await api.get("/vendors");
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

    const fetchUsers = async () => {
      try {
        const response = await api.get("/users");
        setUsers(response.data);
      } catch (error) {
        message.error("Failed to load users");
      }
    };

    fetchVendors();
    fetchUsers();

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

      if (imageFile) {
        formData.append("img", imageFile);
      }

      const response = await api.put(`/vendors/${editingVendor.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      const updatedVendor = response.data;

      setVendors(
        vendors.map((vendor) =>
          vendor.id === editingVendor.id
            ? {
                ...vendor,
                name: updatedVendor.name,
                description: updatedVendor.description,
                img: updatedVendor.img,
              }
            : vendor,
        ),
      );

      message.success("Vendor updated successfully");
      setIsModalVisible(false);
      setEditingVendor(null);
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
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
    setCurrentAdmin(vendor.currentAdmin);
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

  const handleAssignAdmin = async () => {
    if (!selectedAdmin) {
      message.error("Please select a user to assign as admin.");
      return;
    }

    setLoadingAction(true);
    try {
      const form = new FormData();
      form.append("user_id", selectedAdmin);
      form.append("vendor_id", editingVendor.id);
      await api.post("/vendors/assign-vendor-admin", form);

      // Update the local state with the new admin
      setVendors(
        vendors.map((vendor) =>
          vendor.id === editingVendor.id
            ? { ...vendor, currentAdmin: selectedAdmin }
            : vendor,
        ),
      );
      message.success("Admin assigned successfully");
      setSelectedAdmin(null);
    } catch (error) {
      message.error("Failed to assign admin");
    } finally {
      setLoadingAction(false);
    }
  };

  const handleRevokeAdmin = async (vendor) => {
    if (!vendor.currentAdmin) {
      message.error("No admin to revoke.");
      return;
    }

    setLoadingAction(true);
    const form = new FormData();
    form.append("user_id", vendor.currentAdmin);
    form.append("vendor_id", vendor.id); // Extract vendor.id correctly

    try {
      await api.post("/vendors/revoke-vendor-admin", form);

      // Update the local state to remove the admin
      setVendors(
        vendors.map((v) =>
          v.id === vendor.id ? { ...v, currentAdmin: null } : v,
        ),
      );
      message.success("Admin revoked successfully");
    } catch (error) {
      message.error("Failed to revoke admin");
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
                    checked={vendor.isSubscribed}
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
                <p>
                  Current Admin:{" "}
                  {vendor.currentAdmin ? (
                    <span>{vendor.currentAdmin}</span>
                  ) : (
                    <span>No admin assigned</span>
                  )}
                </p>
                <Popconfirm
                  title="Are you sure you want to revoke the current admin?"
                  onConfirm={() => handleRevokeAdmin(vendor)}
                  okText="Yes"
                  cancelText="No"
                >
                  <Button disabled={!vendor.currentAdmin} danger>
                    Revoke Admin
                  </Button>
                </Popconfirm>
              </Card>
            );
          })}

          <Button
            type="dashed"
            style={{ width: 240, height: 240 }}
            onClick={showAddVendorModal}
          >
            <PlusOutlined />
            <br />
            Add Vendor
          </Button>
        </div>
      )}

      <Modal
        title={isAddingVendor ? "Add Vendor" : "Edit Vendor"}
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        onOk={() => {
          isAddingVendor
            ? form.submit()
            : handleSaveVendor(form.getFieldsValue());
        }}
        okText={isAddingVendor ? "Add" : "Save"}
        confirmLoading={loadingAction}
      >
        <Form form={form} layout="vertical" onFinish={handleAddVendor}>
          <Form.Item
            label="Name"
            name="name"
            rules={[{ required: true, message: "Please input vendor name!" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            label="Description"
            name="description"
            rules={[
              { required: true, message: "Please input vendor description!" },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item label="Image">
            <Upload
              beforeUpload={() => false}
              onChange={handleImageUpload}
              showUploadList={false}
            >
              <Button icon={<PlusOutlined />}>Upload Image</Button>
            </Upload>
            {imageUrl && (
              <img
                src={imageUrl}
                alt="Vendor"
                style={{ width: "100%", marginTop: 10 }}
              />
            )}
          </Form.Item>

          {!isAddingVendor && (
            <>
              <Form.Item label="Assign Admin">
                <Select
                  placeholder="Select an admin"
                  onChange={(value) => setSelectedAdmin(value)}
                  style={{ width: "100%" }}
                >
                  {users.map((user) => (
                    <Select.Option key={user.id} value={user.id}>
                      {user.name}
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>
              <Button
                onClick={handleAssignAdmin}
                disabled={!selectedAdmin || loadingAction}
                type="primary"
              >
                Assign Admin
              </Button>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
}

export default VendorManagement;
