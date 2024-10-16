
import 'package:flutter/material.dart';
import '../../model/ItemModel.dart';
import '../../model/VendorModel.dart';
import '../../repository/ItemRepository.dart';
import '../../repository/VendorRepository.dart';

class HomeViewModel extends ChangeNotifier {
  final ItemRepository itemRepository = ItemRepository();
  final VendorRepository vendorRepository = VendorRepository();

  final List<ItemModel> items = [];
  final List<VendorModel> vendors = [];

  int currentPageItems = 1;
  int currentPageVendors = 1;
  bool isLoadingItems = false;
  bool isLoadingVendors = false;

  int totalPagesItems = 5; // Should be fetched from API
  int totalPagesVendors = 5; // Should be fetched from API

  Future<void> fetchItems() async {
    if (isLoadingItems || currentPageItems > totalPagesItems) return;

    isLoadingItems = true;
    print("items is loading");
    notifyListeners();

    try {
      final response = await itemRepository.fetchItems(page: currentPageItems, perPage: 10);
      items.addAll(response.data); // Add fetched items
      totalPagesItems = response.meta.totalPages; // Update total pages if available
      currentPageItems++;
    } catch (e) {
      // Handle errors (e.g., show an error message)

      print(e);
    } finally {
      isLoadingItems = false;
      notifyListeners();
    }
  }

  Future<void> fetchVendors() async {
    if (isLoadingVendors || currentPageVendors > totalPagesVendors) return;

    isLoadingVendors = true;
    notifyListeners();

    try {
      final response = await vendorRepository.fetchVendors(page: currentPageVendors, perPage: 10);
      vendors.addAll(response.data); // Add fetched vendors
      totalPagesVendors = response.meta.totalPages; // Update total pages if available
      currentPageVendors++;
    } catch (e) {
      print("form vendor");
      print(e);
    } finally {
      isLoadingVendors = false;
      notifyListeners();
    }
  }
}
