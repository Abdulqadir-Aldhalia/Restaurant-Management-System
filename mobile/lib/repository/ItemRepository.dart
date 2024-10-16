import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:mobile/repository/Statics.dart';
import '../model/ItemModel.dart';

class ItemRepository {
  final storage = FlutterSecureStorage(); // Access to saved token

  Future<ItemsResponse> fetchItems({required int page, required int perPage}) async {
    try {
      final uri = Uri.parse('${Statics.baseUrl}/items?page=$page&per_page=$perPage');

      final token = await storage.read(key: 'authToken');
      print("token = $token");

      final response = await http.get(
        uri,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Authorization': 'Bearer $token', // Add token if available
        },
      );

      if (response.statusCode == 200) {
        final responseBody = json.decode(response.body);
        return ItemsResponse.fromJson(responseBody);
      } else {
        throw Exception('Failed to load items');
      }
    } catch (e) {
      // Handle exceptions
      print('Error fetching items: $e');
      rethrow;
    }
  }
}
