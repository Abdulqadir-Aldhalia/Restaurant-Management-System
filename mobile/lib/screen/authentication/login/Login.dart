import 'package:flutter/material.dart';
import 'package:mobile/screen/AppDrawer.dart';
import 'LoginViewModel.dart';

class LoginView extends StatefulWidget {
  @override
  _LoginViewState createState() => _LoginViewState();
}

class _LoginViewState extends State<LoginView> {
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final LoginViewModel viewModel = LoginViewModel(); // Create viewModel instance

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawer: const AppDrawer(),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextField(
                controller: emailController,
                decoration: InputDecoration(labelText: 'Email'),
              ),
              SizedBox(height: 10),
              TextField(
                controller: passwordController,
                decoration: InputDecoration(labelText: 'Password'),
                obscureText: true,
              ),
              SizedBox(height: 20),
              viewModel.isLoading
                  ? CircularProgressIndicator()
                  : ElevatedButton(
                onPressed: () async {
                  final email = emailController.text.trim();
                  final password = passwordController.text.trim();
                  await viewModel.login(email, password);
                  setState(() {}); // Manually trigger UI update
                },
                child: Text("Login"),
              ),
              if (viewModel.errorMessage != null)
                Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text(
                    viewModel.errorMessage!,
                    style: TextStyle(color: Colors.red),
                  ),
                ),
              if (viewModel.isLoggedIn) // Show success message only if login is successful
                const Padding(
                  padding: EdgeInsets.only(top: 8.0),
                  child: Text(
                    'Login successful!',
                    style: TextStyle(color: Colors.green),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
