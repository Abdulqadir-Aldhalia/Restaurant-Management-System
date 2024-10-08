import React, { useEffect, useState, useMemo } from "react";
import {
  Card,
  Button,
  Modal,
  Input,
  Upload,
  message,
  Form,
  Spin,
  Pagination,
  InputNumber,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import axios from "axios";
import { baseUrl } from "../../const"; // Adjust path as necessary
import { useSelector } from "react-redux";
import "./VendorItemsManagement.css";
import { useOutletContext } from "react-router";

const VendorItemsManagement = () => {
  const { vendorId } = useOutletContext();
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [items, setItems] = useState([]);
  const [currentItem, setCurrentItem] = useState(null);
  const [form] = Form.useForm();
  const userToken = useSelector((state) => state.user.userToken);
  const [imageFile, setImageFile] = useState(null);
  const [imageUrl, setImageUrl] = useState("");
  const [isModalVisible, setIsModalVisible] = useState(false);

  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const api = useMemo(() => {
    return axios.create({
      baseURL: baseUrl,
      headers: {
        Authorization: `Bearer ${userToken}`,
      },
    });
  }, [userToken]);

  const loadItems = async (page) => {
    setLoading(true);
    try {
      const response = await api.get(
        `/items?filters=vendor_id:${vendorId}&page=${page}`,
      );
      setItems(response.data.data);
      setPagination((prev) => ({
        ...prev,
        total: response.data.meta.total_rows,
      }));
    } catch (error) {
      message.error("Failed to load items");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadItems(pagination.current);
  }, [pagination.current, api]);

  const handleImageUpload = ({ file }) => {
    setImageUrl(URL.createObjectURL(file));
    setImageFile(file);
  };

  // Function to add a new item
  const handleAddItem = async (values) => {
    console.log(values);
    setLoadingAction(true);
    const formData = new FormData();
    formData.append("name", values.name);
    formData.append("price", values.price);
    formData.append("vendor_id", vendorId);
    if (imageFile) {
      formData.append("img", imageFile);
    }

    try {
      await api.post("/items", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      message.success("Item added successfully");
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
      loadItems(pagination.current);
    } catch (error) {
      message.error("Failed to add item");
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false); // Close modal after action
    }
  };

  // Function to update an existing item
  const handleUpdateItem = async (values) => {
    setLoadingAction(true);
    const formData = new FormData();
    formData.append("name", values.name);
    formData.append("price", values.price);
    if (imageFile) {
      formData.append("img", imageFile);
    }

    try {
      await api.put(`/items/${currentItem.id}`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      message.success("Item updated successfully");
      form.resetFields();
      setImageUrl("");
      setImageFile(null);
      loadItems(pagination.current);
    } catch (error) {
      message.error("Failed to update item");
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false); // Close modal after action
    }
  };

  const handleEditItem = (item) => {
    setCurrentItem(item);
    form.setFieldsValue({
      name: item.name,
      price: item.price,
    });
    setImageUrl(`${baseUrl}/${item.img}`);
    setIsModalVisible(true); // Open modal for editing
  };

  const handleDeleteItem = async (itemId) => {
    setLoadingAction(true);
    try {
      await api.delete(`/items/${itemId}`);
      message.success("Item deleted successfully");
      loadItems(pagination.current); // Reload items
    } catch (error) {
      message.error("Failed to delete item");
    } finally {
      setLoadingAction(false);
    }
  };

  const handlePageChange = (page) => {
    setPagination((prev) => ({ ...prev, current: page }));
  };

  // Handle form submission, calling add or update function based on currentItem
  const handleFormSubmit = (values) => {
    if (currentItem) {
      handleUpdateItem(values); // Call update if we are editing an item
    } else {
      handleAddItem(values); // Call add if we are adding a new item
    }
  };

  return (
    <div className="vendor-items-container">
      <h1 className="header-title">Items Management</h1>
      <Button
        className="add-item-button"
        type="primary"
        icon={<PlusOutlined />}
        onClick={() => {
          setCurrentItem(null); // Ensure we're in "Add" mode
          form.resetFields();
          setImageUrl("");
          setImageFile(null);
          setIsModalVisible(true); // Show modal for adding new item
        }}
      >
        Add Item
      </Button>

      <Spin spinning={loading} tip="Loading...">
        <div className="items-list">
          {items.map((item) => (
            <Card key={item.id} hoverable className="item-card">
              <img
                src={`${baseUrl}/${item.img}`}
                alt={item.name}
                className="item-image"
              />
              <h3 className="item-name">{item.name}</h3>
              <p className="item-price">Price: ${item.price.toFixed(2)}</p>
              <div className="item-actions">
                <Button
                  type="primary"
                  onClick={() => handleEditItem(item)}
                  className="edit-button"
                >
                  Edit
                </Button>
                <Button
                  onClick={() => handleDeleteItem(item.id)}
                  danger
                  className="delete-button"
                >
                  Delete
                </Button>
              </div>
            </Card>
          ))}
        </div>
      </Spin>
      <Pagination
        className="pagination"
        current={pagination.current}
        pageSize={pagination.pageSize}
        total={pagination.total}
        onChange={handlePageChange}
      />

      <Modal
        title={currentItem ? "Edit Item" : "Add Item"}
        open={isModalVisible} // Use the state to control modal visibility
        onCancel={() => {
          setCurrentItem(null);
          form.resetFields();
          setImageUrl("");
          setImageFile(null);
          setIsModalVisible(false); // Hide the modal
        }}
        footer={null}
        centered
      >
        {/* Form and other content */}
        <Form form={form} onFinish={handleFormSubmit} layout="vertical">
          {/* Form fields */}
          <Form.Item
            name="name"
            label="Item Name"
            rules={[{ required: true, message: "Please enter item name" }]}
          >
            <Input placeholder="Enter item name" />
          </Form.Item>
          <Form.Item
            name="price"
            label="Price"
            rules={[
              { required: true, message: "Please enter item price" },
              {
                type: "number",
                min: 0.01,
                message: "Price must be greater than 0",
              },
            ]}
          >
            <InputNumber
              min={0.01}
              step={0.01}
              placeholder="Enter price"
              style={{ width: "100%" }}
            />
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
                alt="uploaded"
                className="profile-image-preview"
              />
            )}
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loadingAction}>
              {currentItem ? "Update Item" : "Add Item"}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default VendorItemsManagement;
