import React, { useEffect, useState, useMemo } from "react";
import {
  Card,
  Button,
  Modal,
  Input,
  message,
  Form,
  Spin,
  Pagination,
  Switch,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import axios from "axios";
import { baseUrl } from "../../const"; // Adjust the path as necessary
import { useSelector } from "react-redux";
import { useOutletContext } from "react-router";
import "./VendorTableManagement.css"; // Add a custom CSS file for styling

const VendorTableManagement = () => {
  const { vendorId } = useOutletContext();
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [tables, setTables] = useState([]);
  const [currentTable, setCurrentTable] = useState(null);
  const [form] = Form.useForm();
  const userToken = useSelector((state) => state.user.userToken);
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

  // Load tables for the vendor
  const loadTables = async (page) => {
    setLoading(true);
    try {
      const response = await api.get(
        `/tables?filters=vendor_id:${vendorId}&page=${page}`,
      );
      setTables(response.data.data);
      setPagination((prev) => ({
        ...prev,
        total: response.data.meta.total_rows,
      }));
    } catch (error) {
      message.error("Failed to load tables");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTables(pagination.current);
  }, [pagination.current, api]);

  // Add a new table
  const handleAddTable = async (values) => {
    setLoadingAction(true);
    const formData = new FormData();
    formData.append("name", values.name);
    formData.append("vendor_id", vendorId);
    formData.append("is_available", values.is_available);
    formData.append("is_needs_service", values.is_needs_service);

    try {
      await api.post("/tables", formData);
      message.success("Table added successfully");
      form.resetFields(); // Reset the form fields correctly
      loadTables(pagination.current);
    } catch (error) {
      message.error("Failed to add table");
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false);
    }
  };

  const handleUpdateTable = async (values) => {
    setLoadingAction(true);
    const formData = new FormData();
    formData.append("name", values.name);
    formData.append("is_available", values.is_available);
    formData.append("is_needs_service", values.is_needs_service);

    try {
      await api.put(`/tables/${currentTable.id}`, formData);
      message.success("Table updated successfully");
      form.resetFields();
      loadTables(pagination.current);
    } catch (error) {
      message.error("Failed to update table");
    } finally {
      setLoadingAction(false);
      setIsModalVisible(false);
    }
  };

  const handleEditTable = (table) => {
    setCurrentTable(table);
    form.setFieldsValue({
      name: table.name,
      is_available: table.is_available,
      is_needs_service: table.is_needs_service,
    });
    setIsModalVisible(true); // Open modal for editing
  };

  const handleDeleteTable = async (tableId) => {
    setLoadingAction(true);
    try {
      await api.delete(`/tables/${tableId}`);
      message.success("Table deleted successfully");
      loadTables(pagination.current);
    } catch (error) {
      message.error("Failed to delete table");
    } finally {
      setLoadingAction(false);
    }
  };

  const handleEmptyTable = async (tableId) => {
    setLoadingAction(true);
    try {
      await api.delete(`/tables/${tableId}/empty`);
      message.success("Table emptied successfully");
      loadTables(pagination.current);
    } catch (error) {
      message.error("Failed to empty table");
    } finally {
      setLoadingAction(false);
    }
  };

  const handlePageChange = (page) => {
    setPagination((prev) => ({ ...prev, current: page }));
  };

  const handleFormSubmit = (values) => {
    if (currentTable) {
      handleUpdateTable(values); // Call update if we are editing a table
    } else {
      handleAddTable(values); // Call add if we are adding a new table
    }
  };

  return (
    <div className="vendor-tables-container">
      <h1 className="header-title">Tables Management</h1>
      <Button
        className="add-table-button"
        type="primary"
        icon={<PlusOutlined />}
        onClick={() => {
          setCurrentTable(null); // Ensure we're in "Add" mode
          form.resetFields();
          setIsModalVisible(true); // Show modal for adding new table
        }}
      >
        Add Table
      </Button>

      <Spin spinning={loading} tip="Loading...">
        <div className="tables-list">
          {tables.map((table) => (
            <Card
              key={table.id}
              hoverable
              className={`table-card ${
                table.is_available ? "table-available" : "table-unavailable"
              }`}
            >
              <h3 className="table-name">{table.name}</h3>
              <p className="table-availability">
                Available:{" "}
                <span className={table.is_available ? "yes" : "no"}>
                  {table.is_available ? "Yes" : "No"}
                </span>
              </p>
              <p className="table-service">
                Needs Service:{" "}
                <span className={table.is_needs_service ? "yes" : "no"}>
                  {table.is_needs_service ? "Yes" : "No"}
                </span>
              </p>
              <div className="table-actions">
                <Button
                  type="primary"
                  onClick={() => handleEditTable(table)}
                  className="edit-button"
                >
                  Edit
                </Button>
                <Button
                  onClick={() => handleDeleteTable(table.id)}
                  danger
                  className="delete-button"
                >
                  Delete
                </Button>
                {!table.is_available && (
                  <Button
                    onClick={() => handleEmptyTable(table.id)}
                    className="empty-button"
                  >
                    Empty
                  </Button>
                )}
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
        title={currentTable ? "Edit Table" : "Add Table"}
        open={isModalVisible} // Use the state to control modal visibility
        onCancel={() => {
          setCurrentTable(null);
          form.resetFields();
          setIsModalVisible(false); // Hide the modal
        }}
        footer={null}
        centered
      >
        <Form form={form} onFinish={handleFormSubmit} layout="vertical">
          <Form.Item
            name="name"
            label="Table Name"
            rules={[{ required: true, message: "Please enter table name" }]}
          >
            <Input placeholder="Enter table name" />
          </Form.Item>
          <Form.Item
            name="is_available"
            label="Available"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
          <Form.Item
            name="is_needs_service"
            label="Needs Service"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loadingAction}>
              {currentTable ? "Update" : "Add"}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default VendorTableManagement;
