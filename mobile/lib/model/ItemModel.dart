import '../repository/Statics.dart';
import 'MetaModel.dart';


class ItemModel {
  final String id;
  final String vendorId;
  final String name;
  final double price;
  final String img;
  final DateTime createdAt;
  final DateTime updatedAt;

  ItemModel({
    required this.id,
    required this.vendorId,
    required this.name,
    required this.price,
    required this.img,
    required this.createdAt,
    required this.updatedAt,
  });

  factory ItemModel.fromJson(Map<String, dynamic> json) {
    return ItemModel(
      id: json['id'],
      vendorId: json['vendor_id'],
      name: json['name'],
      price: json['price'].toDouble(),
      img: Statics.prependBaseUrl(json['img']),
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'vendor_id': vendorId,
      'name': name,
      'price': price,
      'img': img,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }


}



class ItemsResponse {
  final Meta meta;
  final List<ItemModel> data;

  ItemsResponse({
    required this.meta,
    required this.data,
  });

  factory ItemsResponse.fromJson(Map<String, dynamic> json) {
    print(json['data']);
    return ItemsResponse(
      meta: Meta.fromJson(json['meta']),
      data: List<ItemModel>.from(json['data'].map((item) => ItemModel.fromJson(item))),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'meta': meta.toJson(),
      'data': List<dynamic>.from(data.map((item) => item.toJson())),
    };
  }
}
