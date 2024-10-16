import 'package:flutter/foundation.dart';
import 'package:mobile/repository/LoginRepository.dart';
import '../../../AuthService.dart';

class LoginViewModel extends ChangeNotifier {
  final LoginRepository loginRepository = LoginRepository();

  bool _isLoading = false;
  String? _errorMessage;
  bool _isLoggedIn = false; // Track login success

  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage;
  bool get isLoggedIn => _isLoggedIn;

  Future<void> login(String email, String password) async {
    _isLoading = true;
    _errorMessage = null;
    _isLoggedIn = false; // Reset login state
    notifyListeners();

    try {
      final result = await loginRepository.loginUser(email, password);
      await AuthService.saveToken(result.token);

      // If login successful, set isLoggedIn to true
      _isLoggedIn = true;
    } catch (error) {
      print('Error during login: $error');
      _errorMessage = error.toString();
      _isLoggedIn = false; // Login failed
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
