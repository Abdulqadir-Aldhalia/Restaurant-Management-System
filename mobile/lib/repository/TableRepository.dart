import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../repository/Statics.dart';
import '../model/TableModel.dart';

class TableRepository {
  final storage = const FlutterSecureStorage(); // Access to saved token

  Future<List<TableModel>> fetchTables({required String vendorId}) async {
    try {
      final uri = Uri.parse('${Statics.baseUrl}/tables?filter=vendor_id:$vendorId');

      final token = await storage.read(key: 'authToken');
      print("token = $token");

      final response = await http.get(
        uri,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final responseBody = json.decode(response.body);
        return (responseBody['data'] as List).map((table) => TableModel.fromJson(table)).toList();
      } else {
        throw Exception('Failed to load tables');
      }
    } catch (e) {
      print('Error fetching tables: $e');
      rethrow;
    }
  }
}
