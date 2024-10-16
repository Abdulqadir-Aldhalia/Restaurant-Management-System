import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:mobile/repository/LoginRepository.dart';

class AuthService {
  final LoginRepository loginRepository = LoginRepository();
  static final storage = FlutterSecureStorage();

  static Future<void> saveToken(String token) async {
    await storage.write(key: 'authToken', value: token);
  }

  static Future<String?> getToken() async {
    return await storage.read(key: 'authToken');
  }

  static Future<void> removeToken() async {
    await storage.delete(key: 'authToken');
  }


}
