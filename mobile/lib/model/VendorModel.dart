import 'package:mobile/repository/Statics.dart';

import 'MetaModel.dart';

class VendorModel {
  final String id;
  final String name;
  final String description;
  final String img;
  final DateTime createdAt;
  final DateTime updatedAt;

  VendorModel({
    required this.id,
    required this.name,
    required this.description,
    required this.img,
    required this.createdAt,
    required this.updatedAt,
  });

  // Factory method to create a Vendor object from JSON
  factory VendorModel.fromJson(Map<String, dynamic> json) {
    return VendorModel(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      img: Statics.prependBaseUrl(json['img']),
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

class VendorResponse {
  final Meta meta;
  final List<VendorModel> data;

  VendorResponse({
    required this.meta,
    required this.data,
  });

  factory VendorResponse.fromJson(Map<String, dynamic> json) {
    print(json['data']);
    return VendorResponse(
      meta: Meta.fromJson(json['meta']),
      data: List<VendorModel>.from(json['data'].map((vendor) => VendorModel.fromJson(vendor))),
    );
  }
}


