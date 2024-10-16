import 'package:flutter/material.dart';
import '../../model/ItemModel.dart';
import '../../model/TableModel.dart';
import 'VendorViewModel.dart';

class VendorView extends StatelessWidget {
  final String vendorId;

  const VendorView({required this.vendorId, super.key});

  @override
  Widget build(BuildContext context) {
    final VendorViewModel viewModel = VendorViewModel();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Vendor Details'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Vendor Tables Section
            FutureBuilder<List<TableModel>>(
              future: viewModel.fetchVendorTables(vendorId),
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                } else if (snapshot.hasError) {
                  return Center(child: Text('Error: ${snapshot.error}'));
                } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
                  final tables = snapshot.data!;
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SectionHeader(title: 'Tables'),
                      ListView.builder(
                        itemCount: tables.length,
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        itemBuilder: (context, index) {
                          final table = tables[index];
                          return TableCard(table: table);
                        },
                      ),
                    ],
                  );
                } else {
                  return const Center(child: Text('No tables available.'));
                }
              },
            ),
            const SizedBox(height: 30),
            // Vendor Items Section
            FutureBuilder<List<ItemModel>>(
              future: viewModel.fetchVendorItems(vendorId),
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                } else if (snapshot.hasError) {
                  return Center(child: Text('Error: ${snapshot.error}'));
                } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
                  final items = snapshot.data!;
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SectionHeader(title: 'Items'),
                      // Wrap GridView in a Container with a fixed height
                      Container(
                        height: 300, // Adjust the height as needed
                        child: GridView.builder(
                          itemCount: items.length,
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          gridDelegate:
                              const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 2,
                            crossAxisSpacing: 16,
                            mainAxisSpacing: 16,
                            childAspectRatio: 0.8,
                          ),
                          itemBuilder: (context, index) {
                            final item = items[index];
                            return ItemCard(item: item);
                          },
                        ),
                      ),
                    ],
                  );
                } else {
                  return const Center(child: Text('No items available.'));
                }
              },
            ),
          ],
        ),
      ),
    );
  }
}

// Section Header Widget
class SectionHeader extends StatelessWidget {
  final String title;

  const SectionHeader({required this.title, super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Text(
        title,
        style: const TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.bold,
          color: Colors.black87,
        ),
      ),
    );
  }
}

class ItemCard extends StatelessWidget {
  final ItemModel item;
  final ValueNotifier<int> quantityNotifier = ValueNotifier<int>(1);

  ItemCard({required this.item, super.key});

  @override
  Widget build(BuildContext context) {
    return Card(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15),
      ),
      elevation: 4,
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(10),
              child: Image.network(
                item.img,
                height: 100,
                width: double.infinity,
                fit: BoxFit.cover,
              ),
            ),
            Text(
              item.name,
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            Text(
              '\$${item.price.toStringAsFixed(2)}',
              style: const TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
            ),
            // Add to Cart Button
            ElevatedButton(
              onPressed: () {
                // Show quantity dialog when the button is pressed
                showQuantityDialog(context);
              },
              child: const Text('Add to Cart'),
            ),
          ],
        ),
      ),
    );
  }

  // Method to show the quantity dialog
  void showQuantityDialog(BuildContext context) {
    int selectedQuantity = 1;

    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Select Quantity'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                'Item: ${item.name}',
                style: const TextStyle(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 10),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  IconButton(
                    onPressed: () {
                      if (selectedQuantity > 1) {
                        selectedQuantity--;
                      }
                      quantityNotifier.value = selectedQuantity;
                    },
                    icon: const Icon(Icons.remove),
                  ),
                  ValueListenableBuilder<int>(
                    valueListenable: quantityNotifier,
                    builder: (context, value, child) {
                      return Text(
                        value.toString(),
                        style: const TextStyle(fontSize: 18),
                      );
                    },
                  ),
                  IconButton(
                    onPressed: () {
                      selectedQuantity++;
                      quantityNotifier.value = selectedQuantity;
                    },
                    icon: const Icon(Icons.add),
                  ),
                ],
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context); // Close the dialog
              },
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                // Add item to cart with selectedQuantity
                addToCart(selectedQuantity);
                Navigator.pop(context); // Close the dialog
              },
              child: const Text('Add to Cart'),
            ),
          ],
        );
      },
    );
  }

  // Method to handle adding to cart logic
  void addToCart(int quantity) {
    // You can handle adding the item to the cart with the specified quantity here
    print('Added ${item.name} with quantity $quantity to cart');
  }
}


// Table Card Widget
class TableCard extends StatelessWidget {
  final TableModel table;

  const TableCard({required this.table, super.key});

  @override
  Widget build(BuildContext context) {
    return Card(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15),
      ),
      elevation: 4,
      child: ListTile(
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        title: Text(
          table.name,
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Available: ${table.isAvailable ? 'Yes' : 'No'}',
              style: TextStyle(
                color: table.isAvailable ? Colors.green : Colors.red,
              ),
            ),
            if (table.isNeedsService)
              const Padding(
                padding: EdgeInsets.only(top: 4.0),
                child: Text(
                  'Needs Service',
                  style: TextStyle(color: Colors.red),
                ),
              ),
          ],
        ),
        trailing: table.isAvailable
            ? ElevatedButton(
                onPressed: () {
                  // Implement logic to reserve the table
                },
                child: const Text('Reserve'),
              )
            : const Icon(Icons.cancel, color: Colors.red),
      ),
    );
  }
}
