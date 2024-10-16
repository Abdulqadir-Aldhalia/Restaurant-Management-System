import 'package:flutter/material.dart';

class AppDrawer extends StatelessWidget {
  const AppDrawer({super.key});

  @override
  Widget build(BuildContext context) {
    return Drawer(
      child: ListView(
        padding: EdgeInsets.zero,
        children: [
          const DrawerHeader(
            decoration: BoxDecoration(color: Colors.blue),
            child: Text("Sadeem"),
          ),
          ListTile(
            leading: const Icon(Icons.home),
            title: const Text('Home'),
            onTap: () {
              Navigator.pop(context); // Close the drawer
              Navigator.pushNamed(context, '/'); // Navigate to Home
            },
          ),
          ListTile(
            leading: const Icon(Icons.shop),
            title: const Text('Vendors'),
            onTap: () {
              Navigator.pop(context); // Close the drawer
              Navigator.pushNamed(context, '/vendors'); // Navigate to Vendors
            },
          ),
          ListTile(
            leading: const Icon(Icons.shopping_cart),
            title: const Text('My Cart'),
            onTap: () {
              Navigator.pop(context); // Close the drawer
              Navigator.pushNamed(context, '/cart'); // Navigate to Cart
            },
          ),
          ListTile(
            leading: const Icon(Icons.reorder),
            title: const Text('My Orders'),
            onTap: () {
              Navigator.pop(context); // Close the drawer
              Navigator.pushNamed(context, '/orders'); // Navigate to Orders
            },
          ),

          ListTile(
            leading: const Icon(Icons.login),
            title: const Text('Login'),
            onTap: () {
              Navigator.pop(context); // Close the drawer
              Navigator.pushNamed(context, '/login'); // Navigate to Orders
            },
          ),
        ],
      ),
    );
  }
}
