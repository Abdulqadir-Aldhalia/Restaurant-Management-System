import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:mobile/screen/item/ItemViewModel.dart';

import '../AppDrawer.dart';

class Item extends StatelessWidget {
  const Item({super.key});

  @override
  Widget build(BuildContext context) {
    final ItemViewModel viewModel = ItemViewModel(); 

    return
      Scaffold(
        appBar: AppBar(title: const Text('Item')),
        drawer: const AppDrawer(), // Use the shared AppDrawer
        body: const Center(
          child: Text('Item Page Content'),
        ),
      );
  }
}
