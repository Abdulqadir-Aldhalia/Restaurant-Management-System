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
  Pagination,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import axios from "axios";
import { baseUrl } from "../../const";
import { useSelector, useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";

function VendorManagement() {
  const [vendors, setVendors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [editingVendor, setEditingVendor] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [imageUrl, setImageUrl] = useState("");
  const [imageFile, setImageFile] = useState(null);
  const [form] = Form.useForm();
  const [isAddingVendor, setIsAddingVendor] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalVendors, setTotalVendors] = useState(0);
  const [vendorsPerPage] = useState(10);
  const [searchQuery, setSearchQuery] = useState(""); // State for search query

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
        navigate("/adminLoginPortal");
      } else if (error.response.status === 403) {
        message.error("Unauthorized to perform this action.");
      }
      return Promise.reject(error);
    },
  );

  useEffect(() => {
    const fetchVendors = async () => {
      setLoading(true); // Ensure loading is true when fetching
      try {
        const response = await api.get(
          `/vendors?page=${currentPage}&limit=${vendorsPerPage}&query=${searchQuery}`, // Include search query
        );
        setVendors(response.data.data); // Adjusted to match the new API response
        setTotalVendors(response.data.total); // Assuming total vendors count comes in the response
      } catch (error) {
        message.error("Failed to load vendors");
      } finally {
        setLoading(false); // Set loading to false after fetching
      }
    };

    fetchVendors();
  }, [userToken, currentPage, searchQuery]); // Fetch vendors when currentPage or searchQuery changes

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
      setVendors((prev) => [response.data, ...prev]);
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

      setVendors((prev) =>
        prev.map((vendor) =>
          vendor.id === editingVendor.id
            ? { ...vendor, ...updatedVendor }
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
  };

  const handleDeleteVendor = async (vendorId) => {
    try {
      await api.delete(`/vendors/${vendorId}`);
      setVendors((prev) => prev.filter((vendor) => vendor.id !== vendorId));
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

  return (
    <div>
      <Input
        placeholder="Search by name"
        onChange={(e) => setSearchQuery(e.target.value)} // Update search query
        style={{ marginBottom: "20px", width: "300px" }}
      />
      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={showAddVendorModal}
        style={{ margin: "20px" }}
      >
        Add Vendor
      </Button>
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
      <Pagination
        current={currentPage}
        pageSize={vendorsPerPage}
        total={totalVendors}
        onChange={(page) => setCurrentPage(page)}
        style={{ marginTop: "16px", textAlign: "center" }}
      />
      <Modal
        title={isAddingVendor ? "Add Vendor" : "Edit Vendor"}
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form
          form={form}
          onFinish={isAddingVendor ? handleAddVendor : handleSaveVendor}
        >
          <Form.Item
            label="Name"
            name="name"
            rules={[
              { required: true, message: "Please input the vendor name!" },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            label="Description"
            name="description"
            rules={[
              {
                required: true,
                message: "Please input the vendor description!",
              },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item label="Image">
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
                alt="Vendor"
                style={{ width: 100, marginTop: 10, marginLeft: 10 }}
              />
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

export default VendorManagement;
