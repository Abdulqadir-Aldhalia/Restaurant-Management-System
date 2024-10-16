import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:mobile/repository/Statics.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../model/SigninModel.dart';

class LoginRepository {
  final storage = FlutterSecureStorage(); // For saving token securely

  Future<LoginResponse> loginUser(String email, String password) async {
    final uri = Uri.parse('${Statics.baseUrl}/signin?username=$email&password=$password');
    try {
      final response = await http.post(uri);
      if (response.statusCode == 200) {
        final responseData = jsonDecode(response.body);
        final token = responseData['token'];

        if (token == null) {
          throw Exception('No token received');
        }

        // Save the token for future requests
        await storage.write(key: 'authToken', value: token);

        return LoginResponse.fromJson(responseData);
      } else {
        throw Exception('Login failed');
      }
    } catch (e) {
      throw Exception('Error logging in: $e');
    }
  }


  Future<void> logout() async {
    await storage.delete(key: 'authToken');
  }
}
