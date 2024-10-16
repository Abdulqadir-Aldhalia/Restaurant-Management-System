import 'package:flutter/material.dart';

import '../AppDrawer.dart';
import 'CartViewModel.dart';


class Cart extends StatelessWidget {
  const Cart({super.key});

  @override
  Widget build(BuildContext context) {
    final CartViewModel viewModel = CartViewModel(); // Create the ViewModel for this page

    return Scaffold(
      appBar: AppBar(title: const Text('My Cart')),
      drawer: const AppDrawer(), // Use the shared AppDrawer
      body: const Center(
        child: Text('Cart Page Content'),
      ),
    );
  }
}
