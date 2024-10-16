class UserModel {
  final String email;
  final String password;

  UserModel({required this.email, required this.password});
}

class LoginResponse {
  final String token;
  final String? message;

  LoginResponse({required this.token, this.message});

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      token: json['token'],
      message: json['message'],
    );
  }
}
