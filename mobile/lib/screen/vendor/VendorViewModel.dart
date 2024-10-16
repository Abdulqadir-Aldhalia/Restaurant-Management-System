import '../../model/ItemModel.dart';
import '../../model/TableModel.dart';
import '../../repository/ItemRepository.dart';
import '../../repository/TableRepository.dart';
import '../../repository/VendorRepository.dart';

class VendorViewModel {
  final VendorRepository _vendorRepository = VendorRepository();
  final ItemRepository _itemRepository = ItemRepository();
  final TableRepository _tableRepository = TableRepository();

  Future<List<ItemModel>> fetchVendorItems(String vendorId) async {
    // TODO: put send pagination as an argument
    const int page = 1;
    const int perPage = 10;
    final itemsResponse = await _itemRepository.fetchItems(page: page, perPage: perPage);
    return itemsResponse.data.where((item) => item.vendorId == vendorId).toList();
  }

  Future<List<TableModel>> fetchVendorTables(String vendorId) async {
    return await _tableRepository.fetchTables(vendorId: vendorId);
  }
}
