import 'dart:convert';

class Meta {
  final int page;
  final int perPage;
  final int totalRows;
  final int totalPages;
  final int from;
  final int to;

  Meta({
    required this.page,
    required this.perPage,
    required this.totalRows,
    required this.totalPages,
    required this.from,
    required this.to,
  });

  factory Meta.fromJson(Map<String, dynamic> json) {
    return Meta(
      page: json['page'],
      perPage: json['per_page'],
      totalRows: json['total_rows'],
      totalPages: json['total_pages'],
      from: json['from'],
      to: json['to'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'page': page,
      'per_page': perPage,
      'total_rows': totalRows,
      'total_pages': totalPages,
      'from': from,
      'to': to,
    };
  }
}
