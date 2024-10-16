class Statics {
  static const String baseUrl = "http://10.0.2.2:8000";

 static String prependBaseUrl(String imgPath) {
    if (!imgPath.startsWith('/')) {
      imgPath = '/$imgPath';
    }
    print("${Statics.baseUrl}$imgPath");
    return '${Statics.baseUrl}$imgPath';
  }

}

