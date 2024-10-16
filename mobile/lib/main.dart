import 'package:flutter/material.dart';
import 'package:mobile/screen/authentication/login/Login.dart';
import 'package:mobile/screen/cart/Cart.dart';
import 'package:mobile/screen/home/Home.dart';
import 'package:mobile/screen/order/Order.dart';
import 'package:mobile/screen/vendor/VendorsView.dart';

void main() {
  runApp(const SadeemApp());
}

class SadeemApp extends StatelessWidget {
  const SadeemApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      routes: {
        '/vendors': (context) => const VendorsView(),
        '/cart': (context) => const Cart(),
        '/orders': (context) => const Order(),
        '/login': (context) => LoginView(),
      },
      home: const HomeView(),
    );
  }
}
