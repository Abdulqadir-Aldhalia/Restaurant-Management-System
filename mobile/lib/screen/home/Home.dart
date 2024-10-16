import 'package:flutter/material.dart';
import 'package:mobile/screen/AppDrawer.dart';
import '../components/VendorCard.dart';
import 'HomeViewModel.dart';

class HomeView extends StatefulWidget {
  const HomeView({super.key});

  @override
  _HomeViewState createState() => _HomeViewState();
}

class _HomeViewState extends State<HomeView> {
  late HomeViewModel viewModel;
  final ScrollController _itemsScrollController = ScrollController();
  final ScrollController _vendorsScrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    viewModel = HomeViewModel();

    // Initial load
    viewModel.fetchItems();
    viewModel.fetchVendors();

    _itemsScrollController.addListener(() {
      if (_itemsScrollController.position.maxScrollExtent ==
          _itemsScrollController.position.pixels &&
          !viewModel.isLoadingItems) {
        viewModel.fetchItems();
      }
    });

    _vendorsScrollController.addListener(() {
      if (_vendorsScrollController.position.maxScrollExtent ==
          _vendorsScrollController.position.pixels &&
          !viewModel.isLoadingVendors) {
        viewModel.fetchVendors();
      }
    });
  }

  @override
  void dispose() {
    _itemsScrollController.dispose();
    _vendorsScrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Home')),
      drawer: AppDrawer(),
      body: Column(
        children: [
          // Horizontally scrollable items
          SizedBox(
            height: 150,
            child: ListView.builder(
              controller: _itemsScrollController,
              scrollDirection: Axis.horizontal,
              itemCount: viewModel.items.length + (viewModel.isLoadingItems ? 1 : 0),
              itemBuilder: (context, index) {
                // Loading indicator for items
                if (index == viewModel.items.length && viewModel.isLoadingItems) {
                  return const Center(child: CircularProgressIndicator());
                }

                // Display item
                final item = viewModel.items[index];
                return Container(
                  width: 120,
                  margin: const EdgeInsets.all(8.0),
                  color: Colors.blueAccent,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Image.network(item.img, height: 60, width: 60, fit: BoxFit.cover),
                      const SizedBox(height: 8),
                      Text(
                        item.name,
                        style: const TextStyle(color: Colors.white),
                        textAlign: TextAlign.center,
                      ),
                      Text(
                        '\$${item.price.toStringAsFixed(2)}',
                        style: const TextStyle(color: Colors.white),
                      ),
                    ],
                  ),
                );
              },
            ),
          ),

          // Vertical separator
          const SizedBox(height: 16),

          // Vertically scrollable vendors
          Expanded(
            child: ListView.builder(
              controller: _vendorsScrollController,
              itemCount: viewModel.vendors.length + (viewModel.isLoadingVendors ? 1 : 0),
              itemBuilder: (context, index) {
                // Loading indicator for vendors
                if (index == viewModel.vendors.length && viewModel.isLoadingVendors) {
                  return const Center(child: CircularProgressIndicator());
                }

                // Display vendor
                final vendor = viewModel.vendors[index];
                return VendorCard(
                  title: vendor.name,
                  subtitle: vendor.description,
                  imagePath: vendor.img,
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}
