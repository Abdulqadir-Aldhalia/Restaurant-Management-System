import 'package:flutter/material.dart';

import '../AppDrawer.dart';
import 'OrderViewModel.dart';

class Order extends StatelessWidget {
  const Order({super.key});

  @override
  Widget build(BuildContext context) {
    final OrderViewModel viewModel = OrderViewModel(); // Create the ViewModel for this page

    return Scaffold(
      appBar: AppBar(title: const Text('My Orders')),
      drawer: const AppDrawer(), // Use the shared AppDrawer
      body: const Center(
        child: Text('Order Page Content'),
      ),
    );
  }
}
