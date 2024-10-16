import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;
import 'package:mobile/Exceptions.dart';
import 'package:mobile/repository/Statics.dart';

import '../model/VendorModel.dart';

class VendorRepository {
  final storage = FlutterSecureStorage(); // Access to saved token

  Future<VendorResponse> fetchVendors({required int page, required int perPage}) async {
    try {
      final uri = Uri.parse('${Statics.baseUrl}/vendors?page=$page&per_page=$perPage');

      final token = await storage.read(key: 'authToken');
      print("token from vendor $token");

      final response = await http.get(
        uri,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final jsonData = jsonDecode(response.body);
        return VendorResponse.fromJson(jsonData);
      } else if (response.statusCode == 404) {
        throw NotFoundException('');
      }
      print(response.body.toString());
      throw Exception('Failed to load vendors');
    } catch (e) {
      print("Error from vendor");
      throw Exception('Error fetching vendors: $e');
    }
  }
}
