import 'package:flutter/material.dart';

class VendorCard extends StatelessWidget {
  final String title;
  final String subtitle;
  final String imagePath;

  const VendorCard({
    Key? key,
    required this.title,
    required this.subtitle,
    required this.imagePath,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: Column(
        children: [
          ListTile(
            leading: const Icon(Icons.store),
            title: Text(title),
            subtitle: Text(subtitle),
          ),
          Padding(
            padding: const EdgeInsets.all(16.0),
            // Wrap the image in an AspectRatio to maintain consistent height-to-width ratio
            child: AspectRatio(
              aspectRatio: 16 / 9, // Set the desired aspect ratio (16:9 in this case)
              child: Image.network(
                imagePath,
                fit: BoxFit.cover, // Ensures the image fills the space while maintaining its aspect ratio
              ),
            ),
          ),
        ],
      ),
    );
  }
}
