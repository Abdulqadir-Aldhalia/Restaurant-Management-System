import 'package:flutter/material.dart';
import '../../model/VendorModel.dart';
import '../components/VendorCard.dart';
import 'VendorView.dart';
import 'VendorsViewModel.dart';

class VendorsView extends StatelessWidget {
  const VendorsView({super.key});

  @override
  Widget build(BuildContext context) {
    final VendorsViewModel viewModel = VendorsViewModel();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Vendors'),
      ),
      body: FutureBuilder<List<VendorModel>>(
        future: viewModel.fetchVendors(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(
              child: CircularProgressIndicator(),
            );
          } else if (snapshot.hasError) {
            return Center(
              child: Text('Error: ${snapshot.error}'),
            );
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(
              child: Text('No vendors available'),
            );
          }

          // If the data is successfully fetched
          final vendors = snapshot.data!;

          return ListView.builder(
            itemCount: vendors.length,
            itemBuilder: (context, index) {
              final vendor = vendors[index];
              return GestureDetector(
                onTap: () {
                  // Navigate to VendorView when a vendor is tapped
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => VendorView(vendorId: vendor.id), // Pass vendorId
                    ),
                  );
                },
                child: VendorCard(
                  title: vendor.name,
                  subtitle: vendor.description,
                  imagePath: vendor.img,
                ),
              );
            },
          );
        },
      ),
    );
  }
}
