# Restaurant Management System

Welcome to the Restaurant Management System repository! This open-source project provides a full-stack solution to manage restaurant operations efficiently and enhance the overall customer experience.

## 🌟 Project Overview

This system manages various aspects of a restaurant business:

- **Menu Management**: Easily add, update, and organize menu items.
- **Order Processing**: Streamline the order placement and tracking process.
- **Reservation System**: Optimize table usage and reduce wait times.
- **Customer Feedback**: Collect insights to improve service quality.

## 🧠 How the System Works

The system serves different types of users with unique roles:

- 🛠 **System Administrator Portal**: Handles global settings, access control, and user roles.
- 🧑‍🍳 **Restaurant Administrator Portal**: Manages restaurant-specific settings, menu, orders, and tables.
- 🍽 **Customer Interface** (In Progress): Allows customers to explore the restaurant, make reservations, and place orders.
- 📱 **Mobile App**:
  - ✅ View restaurants and menu items.
  - 🚧 Placing orders and sending data is still under development.

## ✅ Project Status

| Component                     | Status         |
|------------------------------|----------------|
| Backend (Go)                 | ✅ Completed    |
| Rust Components              | ✅ Completed    |
| Frontend - Admin Portals     | ✅ Completed    |
| Frontend - Customer UI       | 🚧 In Progress  |
| Mobile App (Flutter/Dart)    | 🚧 Partial (UI only) |

## 🛠 Technologies Used

- **Frontend**: JavaScript, HTML, CSS
- **Backend**: Go, Rust (performance-critical)
- **Mobile**: Dart & Flutter

## 📁 Repository Structure

```bash
restaurant-management-system/
├── client-side/        # Frontend code (Admin UI)
├── server-side/        # Backend services (Go)
├── rust/               # Performance modules (Rust)
└── mobile/             # Mobile app (Flutter)
```

## 🚀 Getting Started

Clone the repository from GitHub and set up each part:

```bash
git clone https://github.com/Abdulqadir-Aldhalia/Sadeem-Restaurant.git
cd Sadeem-Restaurant
```

### Install Dependencies

#### Frontend (Admin Interface)

```bash
cd client-side
npm install
npm start
```

#### Backend (Go)

```bash
cd ../server-side
go get ./...
go run main.go
```

#### Mobile (Flutter)

```bash
cd ../mobile
flutter pub get
flutter run
```

## 🤝 Contributing

We welcome your contributions!

1. **Fork** the project.
2. Create a new branch:
   ```bash
   git checkout -b feature/your-feature
   ```
3. Make your changes and commit:
   ```bash
   git commit -m "Add feature: ..."
   ```
4. Push to your branch:
   ```bash
   git push origin feature/your-feature
   ```
5. Open a **Pull Request** via GitHub.

Please make sure your changes follow our code standards and include relevant documentation or tests.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

Thank you for checking out the **Restaurant Management System** project!  
For inquiries or support, reach out to **[Abdulqadir Aldhalia](mailto:Abdulqadir.Aldhalia@hotmail.com)**.
