import '../../model/VendorModel.dart';
import '../../repository/VendorRepository.dart';

class VendorsViewModel {
  final VendorRepository _vendorRepository = VendorRepository();

  Future<List<VendorModel>> fetchVendors() async {
    // Example pagination, adjust as needed
    const int page = 1;
    const int perPage = 10;

    final vendorResponse = await _vendorRepository.fetchVendors(page: page, perPage: perPage);
    return vendorResponse.data;
  }
}
