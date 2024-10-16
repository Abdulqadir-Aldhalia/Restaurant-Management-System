class TableModel {
  final String id;
  final String vendorId;
  final String name;
  final bool isAvailable;
  final String? customerId;
  final bool isNeedsService;

  TableModel({
    required this.id,
    required this.vendorId,
    required this.name,
    required this.isAvailable,
    this.customerId,
    required this.isNeedsService,
  });

  factory TableModel.fromJson(Map<String, dynamic> json) {
    return TableModel(
      id: json['id'],
      vendorId: json['vendor_id'],
      name: json['name'],
      isAvailable: json['is_available'],
      customerId: json['customer_id'],
      isNeedsService: json['Is_needs_service'],
    );
  }
}
