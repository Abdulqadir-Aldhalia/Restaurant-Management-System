import React from "react";
import { Layout, Menu, Button } from "antd";
import {
  UserOutlined,
  ShopOutlined,
  PoweroffOutlined,
} from "@ant-design/icons";
import UserManagement from "./UserManagement";
import VendorManagement from "./VendorManagement";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { removeUserToken } from "../../redux/user/userSlice";

const { Header, Sider, Content } = Layout;

function AdminDashboard() {
  const [currentView, setCurrentView] = React.useState("users");
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const handleMenuClick = (e) => {
    setCurrentView(e.key);
  };

  const handleLogout = () => {
    localStorage.removeItem("userToken");
    dispatch(removeUserToken());
    navigate("/loginAdminPortal");
  };

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider
        collapsible
        theme="dark"
        style={{ background: "#001529", transition: "background 0.3s ease" }}
      >
        <div
          className="logo"
          style={{
            padding: "16px",
            color: "#fff",
            textAlign: "center",
            fontSize: "18px",
            fontWeight: "bold",
          }}
        >
          Admin Panel
        </div>
        <Menu
          theme="dark"
          selectedKeys={[currentView]}
          mode="inline"
          onClick={handleMenuClick}
          style={{ paddingTop: "20px" }}
        >
          <Menu.Item key="users" icon={<UserOutlined />}>
            Users
          </Menu.Item>
          <Menu.Item key="vendors" icon={<ShopOutlined />}>
            Vendors
          </Menu.Item>
        </Menu>
      </Sider>
      <Layout>
        <Header
          style={{
            background: "#fff",
            padding: "0 16px",
            boxShadow: "0 2px 8px rgba(0, 0, 0, 0.1)",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <h2 style={{ margin: 0, color: "#333" }}>Admin Dashboard</h2>
          <Button
            type="primary"
            icon={<PoweroffOutlined />}
            onClick={handleLogout}
            style={{
              marginLeft: "16px",
              backgroundColor: "#e63946", // Red-orange color for logout button
              borderColor: "#e63946", // Match border with button color
              color: "#fff", // Ensure text color is readable
            }}
          >
            Logout
          </Button>
        </Header>
        <Content
          style={{
            margin: "16px",
            padding: "24px",
            background: "#fff",
            borderRadius: "8px",
            boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
          }}
        >
          {currentView === "users" && <UserManagement />}
          {currentView === "vendors" && <VendorManagement />}
        </Content>
      </Layout>
    </Layout>
  );
}

export default AdminDashboard;
